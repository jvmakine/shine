#ifndef __MEM__
#define __MEM__

#define MEM_CLOSURE 1
#define MEM_STRUCT 2
#define MEM_PVECTOR 3

typedef struct RefCount {
    uint8_t type;
    uint32_t count;
} RefCount;

void *heap_malloc(int size);
void *heap_calloc(int count, int size);
void free_rc(RefCount *ptr);

#endif