#include <stdio.h>
#include <stdint.h>
#include "strings.h"

uint8_t strings_equal(struct String *s1, struct String *s2) {
    if (s1 == s2) {
        return 1;
    }
    // TODO: take appended strings into account
    if (s1->base == s2->base) {
        return 1;
    }
    return 0;
}
