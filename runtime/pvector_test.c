#include <stdlib.h>
#include <stdio.h>
#include "pvector.h"

void test_append_increases_length() ;
void test_append_adds_elements();
void test_append_branches();
void test_pvector_append_performance();
void test_pvector_combine();
void test_pvector_equality();

int main() {
    test_append_increases_length();
    test_append_adds_elements();
    test_append_branches();
    test_pvector_combine();
    test_pvector_equality();
    test_pvector_append_performance();
}

void test_append_increases_length() {
    printf("test_append_increases_length: ");
    PVHead *head = pvector_new();
    if (pvector_length(head) != 0) {
        printf("expected new vector size to be 0. Got %d\n", pvector_length(head));
        exit(1);
    }
    PVHead *appended = pvector_append_uint16(pvector_append_uint16(head, (uint16_t)5), (uint16_t)5);
    if (pvector_length(head) != 0) {
        printf("expected vector size not to change after append. Got %d\n", pvector_length(head));
        exit(1);
    }
    if (pvector_length(appended) != 2) {
        printf("expected appended vector size to be 2. Got %d\n", pvector_length(appended));
        exit(1);
    }
    printf("OK\n");
}

void test_append_adds_elements() {
     printf("test_append_adds_elements: ");
     PVHead *head = pvector_new();
     PVHead *appended = pvector_append_uint16(pvector_append_uint16(head, (uint16_t)5), (uint16_t)8);
     if (pvector_get_uint16(appended, 0) != 5) {
        printf("expected appended(0) == 5. Got %d\n", pvector_get_uint16(appended, 0));
        exit(1);
    }
    if (pvector_get_uint16(appended, 1) != 8) {
        printf("expected appended(0) == 8. Got %d\n", pvector_get_uint16(appended, 2));
        exit(1);
    }
     printf("OK\n");
}

void test_append_branches() {
     printf("test_append_branches: ");
     PVHead *head = pvector_append_uint16(pvector_append_uint16(pvector_new(), (uint16_t)5), (uint16_t)8);
     PVHead *b1 = pvector_append_uint16(head, 12);
     PVHead *b2 = pvector_append_uint16(head, 15);
     if (pvector_get_uint16(b1, 2) != 12) {
        printf("expected b1(2) == 12. Got %d\n", pvector_get_uint16(b1, 2));
        exit(1);
    }
    if (pvector_get_uint16(b2, 2) != 15) {
        printf("expected b1(0) == 15. Got %d\n", pvector_get_uint16(b2, 2) );
        exit(1);
    }
     printf("OK\n");
}

void test_pvector_combine() {
    printf("test_pvector_combine: ");
    PVHead *a = pvector_append_uint16(pvector_append_uint16(pvector_new(), 1), 2);
    PVHead *b = pvector_append_uint16(pvector_append_uint16(pvector_new(), 3), 4);
    PVHead *res = pvector_combine_uint16(a, b);
    if (pvector_length(res) != 4) {
        printf("expected new vector size to be 4. Got %d\n", pvector_length(res));
        exit(1);
    }
    if (pvector_get_uint16(res, 1) != 2) {
        printf("expected res(1) == 2. Got %d\n", pvector_get_uint16(res, 1));
        exit(1);
    }
    if (pvector_get_uint16(res, 2) != 3) {
        printf("expected res(2) == 3. Got %d\n", pvector_get_uint16(res, 2));
        exit(1);
    }
    if (pvector_get_uint16(res, 0) != 1) {
        printf("expected res(0) == 1. Got %d\n", pvector_get_uint16(res, 0));
        exit(1);
    }
    printf("OK\n");
}

void test_pvector_equality() {
    printf("test_pvector_equality: ");
    PVHead *a = pvector_append_uint16(pvector_append_uint16(pvector_new(), 1), 2);
    PVHead *b = pvector_append_uint16(pvector_append_uint16(pvector_new(), 1), 2);
    if (!pvector_equals(a, b, sizeof(uint16_t))) {
        printf("independent identical vectors were not equal\n");
        exit(1);
    }
    PVHead *aa = pvector_append_uint16(a, 4);
    PVHead *bb = pvector_append_uint16(b, 5);
    if (pvector_equals(aa, bb, sizeof(uint16_t))) {
        printf("different vectors were equal\n");
        exit(1);
    }
    printf("OK\n");
}

void test_pvector_append_performance() {
    printf("test_pvector_append_performance: ");
    PVHead *head = pvector_new();
    for (int i = 0; i < 1000000; ++i) {
        PVHead *updated = pvector_append_uint16(head, i);
        pvector_free(head);
        head = updated;
    }
    for (int i = 0; i < 1000000; ++i) {
        uint16_t val = pvector_get_uint16(head, i);
        uint16_t exp = i;
        if (val != exp) {
            printf("expected vector(%d) == %d. Got %d\n", i, exp, val);
            exit(1);
        }
    }
    printf("OK\n");
}