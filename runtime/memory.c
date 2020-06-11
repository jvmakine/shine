#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "memory.h"
#include "pvector.h"

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

typedef struct Structure {
    RefCount ref;
    uint16_t clscount;
    uint16_t strucount;
} Structure;

void free_structure(Structure *s) {
    if (s == NULL) {
        return;
    }
    uint32_t refcount = s->ref.count;
    if (refcount == 1) {
        uint16_t clscount = s->clscount;
        uint16_t strucount = s->strucount;
        char* cptr = (char*)(s + 1);
        cptr = cptr + 4; // TODO: remove need for padding adjustment
        void **ptr = (void**)cptr;
        while (clscount > 0) {
            ptr++; // pass the function pointer
            free_rc(*ptr);
            ptr++;
            clscount--;
        }
        while (strucount > 0) {
            free_rc(*ptr);
            ptr++;
            strucount--;
        }
        free(s);
    } else if (refcount > 1) {
        s->ref.count = refcount - 1;
    }
}

void increase_refcount(RefCount *ref) {
    if (ref != NULL) {
        uint32_t refcount = ref->count;
        if (refcount >= 1) {
            ref->count = refcount + 1;
        }
    }
}

void free_rc(RefCount* ref) {
    if (ref == 0) {
        return;
    }
    uint8_t t = ref->type;
    if (t == MEM_PVECTOR) {
        pvector_free((PVHead*)ref);
    } else if (t == MEM_STRUCT) {
        free_structure((Structure*)ref);
    } else {
        fprintf(stderr, "invalid pointer\n");
        exit(1);
    }
}