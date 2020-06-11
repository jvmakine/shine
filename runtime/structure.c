#include <stdint.h>
#include <stdlib.h>
#include "memory.h"
#include "structure.h"

void free_structure(Structure *s) {
    if (s == NULL) {
        return;
    }
    uint32_t refcount = s->ref.count;
    if (refcount == 1) {
        uint16_t strucount = s->strucount;
        char* cptr = (char*)(s + 1);
        // TODO: remove need for padding adjustment
        cptr = cptr + 4;
        void **ptr = (void**)cptr;
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

void free_closure(Closure *s) {
    if (s == NULL) {
        return;
    }
    uint32_t refcount = s->ref.count;
    if (refcount == 1) {
        uint16_t strucount = s->strucount;
        char* cptr = (char*)(s + 1);
        void **ptr = (void**)cptr;
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
