#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "memory.h"

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

void free_structure(void *cls) {
    if (cls == NULL) {
        return;
    }
    char *cptr = cls;
    uint32_t refcount = *(uint32_t*)cptr;
    if (refcount == 1) {
        uint16_t clscount = *(uint16_t*)(cptr + 4);
        uint16_t strucount = *(uint16_t*)(cptr + 6);
        int8_t** ptr = (int8_t**)(cptr + 8);
        while (clscount > 0) {
            ptr++; // pass the function pointer
            free_structure(*ptr);
            ptr++;
            clscount--;
        }
        while (strucount > 0) {
            free_structure(*ptr);
            ptr++;
            strucount--;
        }
        free(cls);
    } else if (refcount > 1) {
        *((uint32_t*)cls) = refcount - 1;
    }
}

void increase_refcount(void *ref) {
    if (ref != NULL) {
        uint32_t refcount = *((uint32_t*)ref);
        if (refcount >= 1) {
            *((uint32_t*)ref) = refcount + 1;
        }
    }
}