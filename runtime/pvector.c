#include <stdlib.h>
#include <string.h>
#include "pvector.h"

void OOM() {
    // TODO: error message
    exit(1);
}

PVHead* pvector_new() {
    PVHead* head = malloc(sizeof(PVHead));
    if (head == 0) { OOM(); }
    head->refcount = 1;
    head->size = 0;
    head->node = 0;
    return head;
}

PVNode* pnode_new() {
    PVNode* node = malloc(sizeof(PVNode));
    if (node == 0) { OOM(); }
    node->refcount = 1;
    return node;
}

PVLeaf_uint16 *pleaf_uint16_new() {
    PVLeaf_uint16* leaf = malloc(sizeof(PVLeaf_uint16));
    if (leaf == 0) { OOM(); }
    leaf->refcount = 1;
    return leaf;
}

PVLeaf_ptr *pleaf_ptr_new() {
    PVLeaf_ptr* leaf = malloc(sizeof(PVLeaf_ptr));
    if (leaf == 0) { OOM(); }
    leaf->refcount = 1;
    return leaf;
}

uint32_t pvector_length(PVHead *vector) {
    return vector->size;
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    PVHead *head = pvector_new();
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
        if (old_size) {
            ((PVNode*)node)->children[0] = vector->node;
            ((PVNode*)vector->node)->refcount++;
        }
    } else {
        node = malloc(sizeof(PVNode));
        if (node == 0) { OOM(); }
        memcpy(node, vector->node, sizeof(PVNode));
        ((PVNode*)node)->refcount = 1;
    }
    head->node = (PVNode*)node;
    while (depth) {
        uint32_t key = (old_size >> depth*BITS) & MASK;
        uint32_t okey = ((old_size - 1) >> depth*BITS) & MASK;
        void *nn = 0;
        depth--;
        if (depth) {
            if (key != okey) {
                nn = pnode_new();
            } else {
                nn = (PVNode*)malloc(sizeof(PVNode));
                if (nn == 0) { OOM(); }
                memcpy(nn, ((PVNode*)node)->children[key], sizeof(PVNode));
            }
        } else {
            if (key != okey) {
                nn = pleaf_uint16_new();
            } else {
                nn = (PVLeaf_uint16*)malloc(sizeof(PVLeaf_uint16));
                if (nn == 0) { OOM(); }
                memcpy(nn, ((PVNode*)node)->children[key], sizeof(PVLeaf_uint16));
            }
        }
        ((PVNode*)node)->children[key] = nn;
        node = nn;
    }
    ((PVLeaf_uint16*)node)->values[old_size & MASK] = value;
    head->size = new_size;
    return head;
}

uint16_t  pvector_get_uint16(PVHead *vector, uint32_t index) {
    uint8_t depth = 0;
    uint32_t i = vector->size;
    while (i) {
        depth++;
        i = i >> BITS;
    }
    void *node = vector->node;
    while (depth) {
        uint32_t key = (index >> depth*BITS) & MASK;
        depth--;
        node = ((PVNode*)node)->children[key];
    }
    return ((PVLeaf_uint16*)node)->values[index & MASK];
}