#include <stdlib.h>
#include <stdint.h>

void free_closure(void *cls) {
    int32_t refcount = *((int32_t*)cls);
    if (refcount <= 1) {
        int16_t clscount = *((int16_t*)(((int32_t*)cls) + 1));
        int8_t* ptr = (int8_t*)((int16_t*)(((int32_t*)cls) + 1) + 1);
        while (clscount > 0) {
            ptr++;
            free_closure(ptr);
            ptr++;
            clscount--;
        }
        free(cls);
    } else {
        *((int32_t*)cls) = refcount - 1;
    }
}