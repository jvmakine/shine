#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

void free_structure(void *cls) {
    if (cls == NULL) {
        return;
    }
    char *cptr = cls;
    uint32_t refcount = *(uint32_t*)cptr;
    if (refcount <= 1) {
        uint16_t clscount = *(uint16_t*)(cptr + 4);
        uint16_t strucount = *(uint16_t*)(cptr + 6);
        int8_t** ptr = (int8_t**)(cptr + 8);
        while (clscount > 0) {
            free_structure(*ptr);
            ptr++;
            ptr++; // pass the function pointer
            clscount--;
        }
        while (strucount > 0) {
            free_structure(*ptr);
            ptr++;
            strucount--;
        }
        free(cls);
    } else {
        *((uint32_t*)cls) = refcount - 1;
    }
}

void increase_refcount(void *ref) {
    if (ref != NULL) {
        uint32_t refcount = *((uint32_t*)ref);
        *((uint32_t*)ref) = refcount + 1;
    }
}