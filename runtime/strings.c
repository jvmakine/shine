#include <stdio.h>
#include <stdint.h>

uint8_t strings_equal(void *s1, void *s2) {
    if (s1 == s2) {
        return 1;
    }
    char *cptr1 = (char*)s1;
    char *cptr2 = (char*)s2;
    // TODO: do not assume 16 byte pointers
    char **pptr1 = (char**)(cptr1 + 16);
    char **pptr2 = (char**)(cptr2 + 16);
    // TODO: take appended strings into account
    if (*pptr1 == *pptr2) {
        return 1;
    }
    return 0;
}