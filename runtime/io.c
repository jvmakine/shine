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
    uint32_t size = str->size;
    for (int i = 0; i < size; ++i) {
        uint16_t c = pv_uint16_get(str, i);
        uint32_t code = c;
        if (code > 0xd7ff && code < 0xe000) {
            // see https://unicode.org/faq/utf_bom.html#utf16-3 
            uint16_t hi = c;
            i++;
            uint16_t lo = pv_uint16_get(str, i);
            uint32_t x = (hi & ((1 << 6) -1)) << 10 | (lo & ((1 << 10) -1));
            uint32_t w = (hi >> 6) & ((1 << 5) - 1);
            uint32_t u = w + 1;
            c = u << 16 | x;
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
    }
    printf("\n");
}