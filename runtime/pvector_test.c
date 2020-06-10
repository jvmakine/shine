#include <stdlib.h>
#include <stdio.h>
#include "pvector.h"

void test_append_increases_length() ;
void test_append_adds_elements();
void test_append_branches();

int main() {
    test_append_increases_length();
    test_append_adds_elements();
    test_append_branches();
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