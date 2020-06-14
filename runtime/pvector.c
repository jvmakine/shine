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
    return node;
}

void pleaf_free(PVLeaf_header *leaf) {
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

PVLeaf_header *pleaf_header_new(uint8_t element_size) {
    PVLeaf_header* leaf = heap_malloc(sizeof(PVLeaf_header) + (element_size << BITS) );
    leaf->refcount = 1;
    memset((leaf + 1), 0, element_size << BITS );
    return leaf;
}

uint32_t pvector_length(PVHead *vector) {
    return vector->size;
}

PVNode *copy_pnode(PVNode* node) {
    PVNode* res = heap_calloc(1,sizeof(PVNode));
    memcpy(res, node, sizeof(PVNode));
    res->refcount = 1;
    return res;
}

PVLeaf_header *copy_pleaf(PVLeaf_header *leaf, uint8_t element_size) {
    PVLeaf_header *res = (PVLeaf_header*)heap_malloc(sizeof(PVLeaf_header) + (element_size << BITS));
    memcpy(res, leaf, sizeof(PVLeaf_header) + (element_size << BITS));
    res->refcount = 1;
    return res;
}

PVHead* pvector_append_leaf(PVHead *vector, uint8_t element_size, void **retval) {
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
            node = pleaf_header_new(element_size);
        }
    } else {
        node = copy_pnode(vector->node);
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
                uint32_t rc = ((PVLeaf_header*)children[i])->refcount;
                if (rc > 0) {
                    ((PVLeaf_header*)(children[i]))->refcount = rc + 1;
                }
            }
            if (children[key] == 0) {
                node = pleaf_header_new(element_size);
            } else {
                node = copy_pleaf(children[key], element_size);
            }
        }
        children[key] = node;
    }
    *retval = ((PVLeaf_header*)node) + 1;
    return head;
}

void *pvector_get_leaf(PVHead *vector, uint32_t index, uint8_t element_size) {
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
    return ((void*)(((PVLeaf_header*)node)+1));
}

uint16_t pvector_get_uint16(PVHead *vector, uint32_t index) {
    uint16_t *ptr = pvector_get_leaf(vector, index, sizeof(uint16_t));
    return ptr[index & MASK];
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    uint16_t *ptr;
    PVHead *head = pvector_append_leaf(vector, sizeof(uint16_t), (void*)&ptr);
    ptr[vector->size & MASK] = value;
    return head;
}

PVHead* pvector_combine_uint16(PVHead *a, PVHead *b) {
    uint32_t lenb = pvector_length(b);
    PVHead *n = a;
    for(uint32_t i = 0; i < lenb; ++i) {
        uint16_t c = pvector_get_uint16(b, i);
        PVHead *u = pvector_append_uint16(n, c);
        // TODO: Optimise
        pvector_free(n);
        n = u;
    }
    return n;
}

uint8_t pleaf_equals(PVLeaf_header *a, PVLeaf_header *b, uint8_t element_size) {
    // TODO: Implement deep equals for pointers
    if (a == b) {
        return 1;
    }
    uint32_t len = BRANCH * element_size;
    uint8_t* ptra = (uint8_t*)(a + 1);
    uint8_t* ptrb = (uint8_t*)(b + 1);
    for(uint16_t i = 0; i < len; ++i) {
        if(ptra[i] != ptrb[i]) {
            return 0;
        }
    }
    return 1;
}

uint8_t pnode_equals(PVNode *a, PVNode *b, uint8_t depth, uint8_t element_size) {
    if (a == b) {
        return 1;
    }
    for(uint8_t i = 0; i < BRANCH; i++) {
        if (a->children[i] == 0) {
            return 1;
        }
        if (depth > 0) {
            if(!pnode_equals((PVNode*)a->children[i], (PVNode*)b->children[i], depth - 1, element_size)) {
                return 0;
            }
        } else {
            if(!pleaf_equals((PVLeaf_header*)a->children[i], (PVLeaf_header*)b->children[i], element_size)) {
                return 0;
            }
        }
    }
    return 1;
}

uint8_t pvector_equals(PVHead *a, PVHead *b, uint8_t element_size) {
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
        return pnode_equals(a->node, b->node, depth, element_size);
    } else {
        return pleaf_equals(a->node, b->node, element_size);
    }
}