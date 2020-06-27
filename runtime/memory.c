#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "memory.h"
#include "pvector.h"
#include "structure.h"

void *heap_malloc(int size) {
    void *result = malloc(size);
    if (result == 0) {
        fprintf(stderr, "malloc failed\n");
        exit(1);
    }
    return result;
}

void *heap_calloc(int count, int size) {
    void *result = calloc(count, size);
    if (result == 0) {
        fprintf(stderr, "calloc failed\n");
        exit(1);
    }
    return result;
}

void increase_refcount(RefCount *ref) {
    if (ref != NULL) {
        uint32_t refcount = ref->count;
        if (refcount != CONSTANT_REF) {
            ref->count = refcount + 1;
        }
    }
}

void free_rc(RefCount* ref) {
    if (ref == 0) {
        return;
    }
    if (ref->count == CONSTANT_REF) {
        return;
    }
    uint8_t t = ref->type;
    if (t == MEM_PVECTOR) {
        pv_free((PVHead*)ref);
    } else if (t == MEM_STRUCT) {
        free_structure((Structure*)ref);
    } else if (t == MEM_CLOSURE) {
        free_closure((Closure*)ref);
    } else {
        fprintf(stderr, "invalid pointer %p\n", ref);
        exit(1);
    }
}