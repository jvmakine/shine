#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "pvector.h"

uint8_t pvector_depth(PVHead *vector) {
    uint32_t size = vector->size;
    if (size == 0) {
        return 0;
    }
    uint8_t depth = 0;
    uint32_t i = (size - 1) >> BITS;
    while (i) {
        depth++;
        i = i >> BITS;
    }
    return depth;
}

PVHead* pvector_new() {
    PVHead* head = heap_calloc(1, sizeof(PVHead));
    head->ref.count = 1;
    head->ref.type = MEM_PVECTOR;
    head->size = 0;
    head->node = 0;
    return head;
}

PVNode* pnode_new() {
    PVNode* node = heap_calloc(1, sizeof(PVNode));
    node->refcount = 1;
    node->indextable = 0;
    return node;
}

void *pleaf_new(uint32_t leaf_size) {
    PVLeaf_uint16* leaf = heap_calloc(1, leaf_size);
    leaf->refcount = 1;
    return (void*)leaf;
}

void pleaf_free(PVLeaf_uint16 *leaf) {
    // TODO: release references properly
    uint32_t rc = leaf->refcount;
    if (rc == 0) {
        return;
    }
    if (rc > 1) {
        leaf->refcount = rc - 1;
        return;
    }
    free(leaf);
}

void pnode_free(PVNode *node, int depth) {
    if (node == 0) {
        return;
    }
    uint32_t rc = node->refcount;
    if (rc == 0) {
        return;
    }
    if (rc > 1) {
        node->refcount = rc - 1;
        return;
    }
    for(int i = 0; i < BRANCH; ++i) {
        if (node->children[i] != 0) {
            if (depth > 1) {
                pnode_free(node->children[i], depth - 1);
            } else {
                pleaf_free(node->children[i]);
            }
        }
    }
    if (node->indextable != 0) {
        free(node->indextable);
    }
    free(node);
}

void pvector_free(PVHead *vector) {
    if (vector == 0 || vector->ref.count == 0) {
        return;
    }
    uint8_t depth = pvector_depth(vector);
    if (vector->ref.count > 1) {
        vector->ref.count = vector->ref.count - 1;
        return;
    }
    if (vector->size > 0) {
        if (depth > 0) {
            pnode_free(vector->node, depth);
        } else {
            pleaf_free(vector->node);
        }
    }
    free(vector);
}

uint32_t pvector_length(PVHead *vector) {
    return vector->size;
}

PVNode *copy_pnode(PVNode* node) {
    PVNode* res = heap_malloc(sizeof(PVNode));
    memcpy(res, node, sizeof(PVNode));
    res->refcount = 1;
    if (node->indextable != 0) {
        node->indextable = 0;
        printf("%p\n", node->indextable);
        res->indextable = heap_malloc(BRANCH*sizeof(uint32_t));
        memcpy(res->indextable, node->indextable, BRANCH*sizeof(uint32_t));
    }
    return res;
}

void *copy_pleaf(void *leaf, uint32_t leaf_size) {
    void *res = heap_malloc(leaf_size);
    memcpy(res, leaf, leaf_size);
    ((PVLeaf_uint16*)res)->refcount = 1;
    return res;
}

PVHead* pvector_append_leaf(PVHead *vector, uint32_t leaf_size, void **retval) {
    uint32_t old_size = vector->size;
    uint32_t new_size = old_size + 1;
    char new_node = 0;
    uint8_t depth = 0;
    uint32_t o = (old_size - 1) >> BITS;
    uint32_t n = (new_size - 1) >> BITS;
    if (old_size > 0) {
        while(n) {
            if (!o) { new_node = 1; }
            o = o >> BITS;
            n = n >> BITS;
            depth++;
        }
    } else {
        depth = 0;
        new_node = 1;
    }
    void *node = 0;
    if (new_node || vector->node == 0) {
        if (depth > 0) {
            node = pnode_new();
            PVNode *vn = vector->node;
            if (vn) {
                ((PVNode*)node)->children[0] = vn;
                uint32_t rc = vn->refcount;
                if (rc > 0) {
                    vn->refcount = rc + 1;
                }
            }
        } else {
            node = pleaf_new(leaf_size);
        }
    } else {
        if (depth > 0) {
            node = copy_pnode(vector->node);
         } else {
             node = copy_pleaf(vector->node, leaf_size);
         }
    }

    PVHead *head = pvector_new();
    head->node = (PVNode*)node;
    head->size = new_size;

    while (depth) {
        uint8_t shift = depth*BITS;
        depth--;
        uint32_t key = (old_size >> shift) & MASK;
        void** children = ((PVNode*)node)->children;
        if (depth) {
            for(uint8_t i = 0; i < key; i++) {
                uint32_t rc = ((PVNode*)children[i])->refcount;
                if (rc > 0) {
                    ((PVNode*)children[i])->refcount = rc + 1;
                }
            }
            if (children[key] == 0) {
                node = pnode_new();
            } else {
                node = copy_pnode(children[key]);
            }
        } else {
            for(uint8_t i = 0; i < key; i++) {
                uint32_t rc = ((PVLeaf_uint16*)children[i])->refcount;
                if (rc > 0) {
                    ((PVLeaf_uint16*)(children[i]))->refcount = rc + 1;
                }
            }
            if (children[key] == 0) {
                node = pleaf_new(leaf_size);
            } else {
                node = copy_pleaf(children[key], leaf_size);
            }
        }
        children[key] = node;
    }
    ((PVLeaf_uint16*)node)->size++;
    *retval = node;
    return head;
}

void* pvector_get_leaf(PVHead *vector, uint32_t index) {
    if (index >= vector->size) {
        fprintf(stderr, "pvector index out of bounds: got %d, size %d\n", index, vector->size);
        exit(1);
    }
    uint8_t depth = pvector_depth(vector);
    void *node = vector->node;
    while (depth) {
        uint32_t key = (index >> depth*BITS) & MASK;
        depth--;
        node = ((PVNode*)node)->children[key];
    }
    return (void*)node;
}

uint16_t pvector_get_uint16(PVHead *vector, uint32_t index) {
    PVLeaf_uint16 *leaf = pvector_get_leaf(vector, index);
    return leaf->data[index & MASK];
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    PVLeaf_uint16 *leaf;
    PVHead *head = pvector_append_leaf(vector, sizeof(PVLeaf_uint16), (void*)&leaf);
    leaf->data[vector->size & MASK] = value;
    return head;
}

PVHead* pvector_combine_uint16(PVHead *a, PVHead *b) {
    uint32_t lenb = pvector_length(b);
    PVHead *n = a;
    for(uint32_t i = 0; i < lenb; ++i) {
        uint16_t c = pvector_get_uint16(b, i);
        PVHead *u = pvector_append_uint16(n, c);
        // TODO: Optimise
        if (i > 0) {
            pvector_free(n);
        }
        n = u;
    }
    return n;
}

uint8_t pleaf_equals(void *a, void *b, uint32_t leaf_size) {
    if (a == b) {
        return 1;
    }
    return !memcmp(a, b, leaf_size);
}

uint8_t pnode_equals(PVNode *a, PVNode *b, uint8_t depth, uint8_t leaf_size) {
    if (a == b) {
        return 1;
    }
    for(uint8_t i = 0; i < BRANCH; i++) {
        if (a->children[i] == 0) {
            return 1;
        }
        if (depth > 0) {
            if(!pnode_equals((PVNode*)a->children[i], (PVNode*)b->children[i], depth - 1, leaf_size)) {
                return 0;
            }
        } else {
            if(!pleaf_equals(a->children[i], b->children[i], leaf_size)) {
                return 0;
            }
        }
    }
    return 1;
}

uint8_t pvector_equals(PVHead *a, PVHead *b, uint32_t leaf_size) {
    if (a->size != b->size) {
        return 0;
    }
    if (a == b) {
        return 1;
    }
    if (a->size == 0 && b->size == 0) {
        return 1;
    }
    uint32_t depth = pvector_depth(a);
    if (depth > 0) {
        return pnode_equals(a->node, b->node, depth, leaf_size);
    } else {
        return pleaf_equals(a->node, b->node, leaf_size);
    }
}