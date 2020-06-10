#include <stdlib.h>
#include <string.h>
#include "memory.h"
#include "pvector.h"

PVHead* pvector_new() {
    PVHead* head = heap_alloc(sizeof(PVHead));
    head->refcount = 1;
    head->size = 0;
    head->node = 0;
    return head;
}

PVNode* pnode_new() {
    PVNode* node = heap_alloc(sizeof(PVNode));
    node->refcount = 1;
    return node;
}

PVLeaf_header *pleaf_header_new(uint8_t element_size) {
    PVLeaf_header* leaf = heap_alloc(sizeof(PVLeaf_header) + (element_size << BITS) );
    leaf->refcount = 1;
    return leaf;
}

uint32_t pvector_length(PVHead *vector) {
    return vector->size;
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
        node = heap_alloc(sizeof(PVNode));
        memcpy(node, vector->node, sizeof(PVNode));
        ((PVNode*)node)->refcount = 1;
    }

    PVHead *head = pvector_new();
    head->node = (PVNode*)node;
    head->size = new_size;

    while (depth) {
        uint8_t shift = depth*BITS;
        uint32_t key = (old_size >> shift) & MASK;
        uint32_t okey = ((old_size - 1) >> shift) & MASK;
        void *nn = 0;
        if (--depth) {
            for(uint8_t i = 0; i < key; i++) {
                ((PVNode*)(((PVNode*)node)->children[i]))->refcount++;
            }
            if (key != okey) {
                nn = pnode_new();
            } else {
                nn = (PVNode*)heap_alloc(sizeof(PVNode));
                memcpy(nn, ((PVNode*)node)->children[key], sizeof(PVNode));
            }
        } else {
            for(uint8_t i = 0; i < key; i++) {
                ((PVLeaf_header*)(((PVNode*)node)->children[i]))->refcount++;
            }
            if (key != okey) {
                nn = pleaf_header_new(element_size);
            } else {
                nn = (PVLeaf_header*)heap_alloc(sizeof(PVLeaf_header) + (element_size << BITS));
                memcpy(nn, ((PVNode*)node)->children[key], sizeof(PVLeaf_header) + (element_size << BITS));
            }
        }
        ((PVNode*)node)->children[key] = nn;
        node = nn;
    }
    *retval = ((PVLeaf_header*)node) + 1;
    return head;
}

void *pvector_get_leaf(PVHead *vector, uint32_t index, uint8_t element_size) {
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
    return ((void*)(((PVLeaf_header*)node)+1));
}

uint16_t pvector_get_uint16(PVHead *vector, uint32_t index) {
    uint16_t *ptr = pvector_get_leaf(vector, index, sizeof(uint16_t));
    return ptr[index & MASK];
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    void *ptr;
    PVHead *head = pvector_append_leaf(vector, sizeof(uint16_t), &ptr);
    ((uint16_t*)ptr)[vector->size & MASK] = value;
    return head;
}