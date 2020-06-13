#include <stdio.h>
#include <stdint.h>
#include "strings.h"
#include "pvector.h"

void print_int(long i) {
    printf("%ld\n", i);
}

void print_real(double d) {
    printf("%f\n", d);
}

void print_bool(int8_t b) {
    if (b) {
       printf("true\n");
    } else {
       printf("false\n");
    }
}

void print_string(PVHead *str) {
    for (int i = 0; i < str->size; ++i) {
        uint16_t c = pvector_get_uint16(str, i);
        // TODO: UTF-8 conversion
        printf("%c", c);
    }
    printf("\n");
}