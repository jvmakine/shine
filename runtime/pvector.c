#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "memory.h"
#include "pvector.h"

PVHead* pvector_new() {
    PVHead* head = heap_calloc(1, sizeof(PVHead));
    head->refcount = 1;
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
    if (leaf->refcount == 0) {
        return;
    }
    if (leaf->refcount > 1) {
        leaf->refcount--;
        return;
    }
    free(leaf);
}

void pnode_free(PVNode *node, int depth) {
    if (depth <= 0 || node == 0) {
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

uint8_t pvector_depth(PVHead *vector) {
    uint8_t depth = 0;
    uint32_t i = vector->size;
    while (i) {
        depth++;
        i = i >> BITS;
    }
    return depth;
}

void pvector_free(PVHead *vector) {
    if (vector == 0 || vector->refcount == 0) {
        return;
    }
    uint8_t depth = pvector_depth(vector);
    if (vector->refcount > 1) {
        vector->refcount = vector->refcount - 1;
        return;
    }
    pnode_free(vector->node, depth);
    free(vector);
}

PVLeaf_header *pleaf_header_new(uint8_t element_size) {
    PVLeaf_header* leaf = heap_malloc(sizeof(PVLeaf_header) + (element_size << BITS) );
    leaf->refcount = 1;
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
    uint32_t o = old_size;
    uint32_t n = new_size;
    while(n) {
        if (!o) { new_node = 1; }
        o = o >> BITS;
        n = n >> BITS;
        depth++;
    }
    void *node = 0;
    if (new_node) {
        node = pnode_new();
        if (vector->node) {
            ((PVNode*)node)->children[0] = vector->node;
            ((PVNode*)vector->node)->refcount++;
        }
    } else {
        node = copy_pnode(vector->node);
    }

    PVHead *head = pvector_new();
    head->node = (PVNode*)node;
    head->size = new_size;

    while (depth) {
        uint8_t shift = depth*BITS;
        uint32_t key = (old_size >> shift) & MASK;
        void *nn = 0;
        void** children = ((PVNode*)node)->children;
        if (--depth) {
            for(uint8_t i = 0; i < key; i++) {
                ((PVNode*)children[i])->refcount++;
            }
            if (children[key] == 0) {
                nn = pnode_new();
            } else {
                nn = copy_pnode(children[key]);
            }
        } else {
            for(uint8_t i = 0; i < key; i++) {
                ((PVLeaf_header*)(children[i]))->refcount++;
            }
            if (children[key] == 0) {
                nn = pleaf_header_new(element_size);
            } else {
                nn = copy_pleaf(children[key], element_size);
            }
        }
        children[key] = nn;
        node = nn;
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
        PVHead *u = pvector_append_uint16(n, pvector_get_uint16(b, i));
        pvector_free(n);
        n = u;
    }
    return n;
}