#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

void free_closure(void *cls) {
    if (cls == NULL) {
        return;
    }
    int32_t refcount = *((int32_t*)cls);
    if (refcount <= 1) {
        int16_t* cptr = (int16_t*)(cls + 1);
        int16_t clscount = *(cptr);
        int8_t** ptr = (int8_t**)(cptr + 1);
        while (clscount > 0) {
            ptr++; // pass the function pointer
            free_closure(*ptr);
            ptr++;
            clscount--;
        }
        free(cls);
    } else {
        *((int32_t*)cls) = refcount - 1;
    }
}

void increase_refcount(void *ref) {
    if (ref != NULL) {
        int32_t refcount = *((int32_t*)ref);
        *((int32_t*)ref) = refcount + 1;
    }
}