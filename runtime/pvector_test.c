#include <stdlib.h>
#include <stdio.h>
#include <sys/time.h>
#include "pvector.h"

void test_append_increases_length() ;
void test_append_adds_elements();
void test_append_branches();
void test_pvector_combine();
void test_pvector_equality();
void test_pvector_depth();
void test_pvector_append_performance();
void test_pvector_combine_performance();

int main() {
    test_pvector_depth();
    test_append_increases_length();
    test_append_adds_elements();
    test_append_branches();
    test_pvector_combine();
    test_pvector_equality();
    test_pvector_append_performance();
    test_pvector_combine_performance();
}

void test_pvector_depth() {
    printf("test_pvector_depth: ");
    PVHead *head = pvector_new();
    head->size = 0;
    uint8_t d = pvector_depth(head);
    if (d != 0) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 0, d);
        exit(1);
    }

    head->size = BRANCH + 3;
    d = pvector_depth(head);
    if (d != 1) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 1, d);
        exit(1);
    }

    head->size = BRANCH * BRANCH + 3;
    d = pvector_depth(head);
    if (d != 2) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 2, d);
        exit(1);
    }

    head->size = BRANCH / 2;
    d = pvector_depth(head);
    if (d != 0) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 0, d);
        exit(1);
    }

    head->size = BRANCH;
    d = pvector_depth(head);
    if (d != 0) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 0, d);
        exit(1);
    }

    head->size = BRANCH * BRANCH;
    d = pvector_depth(head);
    if (d != 1) {
        printf("expected new vector size %d to have depth %d. Got %d\n", head->size, 1, d);
        exit(1);
    }
    printf("OK\n");
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
    if (!pvector_equals(a, b, sizeof(PVLeaf_uint16))) {
        printf("independent identical vectors were not equal\n");
        exit(1);
    }
    PVHead *aa = pvector_append_uint16(a, 4);
    PVHead *bb = pvector_append_uint16(b, 5);
    if (pvector_equals(aa, bb, sizeof(PVLeaf_uint16))) {
        printf("different vectors were equal\n");
        exit(1);
    }
    printf("OK\n");
}

void test_pvector_append_performance() {
    printf("test_pvector_append_performance: ");
    PVHead *head = pvector_new();
    struct timeval tval_before, tval_after, tval_result;
    gettimeofday(&tval_before, NULL);
    for (int i = 0; i < 1000000; ++i) {
        PVHead *updated = pvector_append_uint16(head, i);
        pvector_free(head);
        head = updated;
    }
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("append %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);
    
    gettimeofday(&tval_before, NULL);
    for (int i = 0; i < 1000000; ++i) {
        uint16_t val = pvector_get_uint16(head, i);
        uint16_t exp = i;
        if (val != exp) {
            printf("expected vector(%d) == %d. Got %d\n", i, exp, val);
            exit(1);
        }
    }
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("index %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);
    printf("OK\n");
}

void test_pvector_combine_performance() {
    printf("test_pvector_combine_performance: ");
    PVHead *v1 = pvector_new();
    PVHead *v2 = pvector_new();
    for (int i = 0; i < 500000; ++i) {
        PVHead *updated = pvector_append_uint16(v1, i);
        pvector_free(v1);
        v1 = updated;
        updated = pvector_append_uint16(v2,500000 + i);
        pvector_free(v2);
        v2 = updated;
    }

    struct timeval tval_before, tval_after, tval_result;
    gettimeofday(&tval_before, NULL);
    PVHead *combined = pvector_combine_uint16(v1, v2);
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);

    pvector_free(v1);
    pvector_free(v2);
    if (combined->size != 1000000) {
        printf("expected combined->size == %d. Got %d\n", 1000000, combined->size);
        exit(1);
    }
    printf("combine %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);

    gettimeofday(&tval_before, NULL);
    for (int i = 0; i < 1000000; ++i) {
        uint16_t val = pvector_get_uint16(combined, i);
        uint16_t exp = i;
        if (val != exp) {
            printf("expected combined(%d) == %d. Got %d\n", i, exp, val);
            exit(1);
        }
    }
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("index %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);
    pvector_free(combined);
    printf("OK\n");
}