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
