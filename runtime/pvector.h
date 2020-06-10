#include <stdint.h>

#define BITS 5
#define BRANCH (1<<BITS)
#define MASK (BRANCH-1)

typedef struct PVHead {
    uint32_t refcount;
    uint32_t size;
    void* node; 
} PVHead;

typedef struct PVNode {
    uint32_t refcount;
    void* children[BRANCH]; 
} PVNode;

typedef struct PVLeaf_uint16 {
    uint32_t refcount;
    uint16_t values[BRANCH];
} PVLeaf_uint16;

typedef struct PVLeaf_ptr {
    uint32_t refcount;
    void* values[BRANCH];
} PVLeaf_ptr;

PVHead* pvector_new();
uint32_t pvector_length(PVHead *vector);

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value);
uint16_t  pvector_get_uint16(PVHead *vector, uint32_t index);
