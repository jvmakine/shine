#include <stdio.h>
#include <stdint.h>

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

void print_string(int8_t *b) {
    char *cptr = (char*)b;
    printf("%s\n", cptr);
}