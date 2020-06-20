#include <stdint.h>
#include "memory.h"

#define BITS 5
#define BRANCH (1<<BITS)
#define MASK (BRANCH-1)

typedef struct PVHead {
    RefCount ref;
    uint32_t size;
    void* node; 
} PVHead;

// Common header between nodes and leaves
typedef struct PVH {
    uint8_t depth;
    uint32_t refcount;
    uint32_t size;
} PVH;

typedef struct PVNode {
    PVH header;
    uint32_t *indextable;
    void* children[BRANCH];
} PVNode;

typedef struct PVLeaf_uint16 {
    PVH header;
    uint16_t data[BRANCH];
} PVLeaf_uint16;

PVHead* pvector_new();
uint32_t pvector_length(PVHead *vector);
void pvector_free(PVHead *vector);
uint8_t pvector_equals(PVHead *a, PVHead *b, uint32_t leaf_size);
uint8_t pvector_depth(PVHead *vector);

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value);
uint16_t  pvector_get_uint16(PVHead *vector, uint32_t index);
PVHead* pvector_combine_uint16(PVHead *a, PVHead *b);

uint8_t needs_rebalancing(PVNode* left, PVNode* right);