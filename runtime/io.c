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
        uint32_t code = c;
        if (code > 0xd7ff && code < 0xe000) {
            // TODO: large code points
        }
        if (code <= 0x7f) {
            printf("%c", c);
        } else if (code <= 0x7ff) {
            uint8_t b1 = 0b11000000 + ((code >> 6) & 0b00011111);
            uint8_t b2 = 0b10000000 + (code & 0b00111111);
            printf("%c%c", b1, b2);
        } else if (code <= 0xffff) {
            uint8_t b1 = 0b11100000 + ((code >> 12) & 0b00001111);
            uint8_t b2 = 0b10000000 + ((code >> 6) & 0b00111111);
            uint8_t b3 = 0b10000000 + (code & 0b00111111);
            printf("%c%c%c", b1, b2, b3);
        } else {
            uint8_t b1 = 0b11110000 + ((code >> 18) & 0b00000111);
            uint8_t b2 = 0b10000000 + ((code >> 12) & 0b00111111);
            uint8_t b3 = 0b10000000 + ((code >> 6) & 0b00111111);
            uint8_t b4 = 0b10000000 + (code & 0b00111111);
            printf("%c%c%c%c", b1, b2, b3, b4);
        }
        // TODO: UTF-8 conversion
    }
    printf("\n");
}