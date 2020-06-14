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

typedef struct PVNode {
    uint32_t refcount;
    void* children[BRANCH]; 
} PVNode;

typedef struct PVLeaf_header {
    uint32_t refcount;
} PVLeaf_header;

PVHead* pvector_new();
uint32_t pvector_length(PVHead *vector);
void pvector_free(PVHead *vector);
uint8_t pvector_equals(PVHead *a, PVHead *b, uint8_t element_size);
uint8_t pvector_depth(PVHead *vector);

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value);
uint16_t  pvector_get_uint16(PVHead *vector, uint32_t index);
PVHead* pvector_combine_uint16(PVHead *a, PVHead *b);
