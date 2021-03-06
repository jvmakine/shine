#include <stdint.h>
#include "memory.h"

#define BITS 5
#define BRANCH (1<<BITS)
#define MASK (BRANCH-1)

// Common header between nodes and leaves
typedef struct PVH {
    uint8_t depth;
    uint32_t refcount;
    uint32_t size;
} PVH;

typedef struct PVHead {
    RefCount ref;
    uint32_t size;
    PVH* node; 
} PVHead;

typedef struct PVNode {
    PVH header;
    uint32_t *indextable;
    PVH* children[BRANCH];
} PVNode;

typedef struct PVLeaf_uint16 {
    PVH header;
    uint16_t data[BRANCH];
} PVLeaf_uint16;

// Public interface

PVHead* pv_new();

void pv_free(PVHead *vector);

uint32_t pv_length(PVHead *vector);

PVHead* pv_concatenate(PVHead *a, PVHead *b);

uint8_t pv_uint16_equals(PVHead *a, PVHead *b);

PVHead* pv_uint16_append(PVHead *vector, uint16_t value);

uint16_t pv_uint16_get(PVHead *vector, uint32_t index);

// Internal functions

uint8_t pv_depth(PVHead *vector);

uint8_t pn_needs_rebalancing(PVNode* left, PVNode* right);

void pn_balance_level(PVNode* left, PVNode* right, PVNode **leftOut, PVNode **rightOut);