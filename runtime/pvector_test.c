#include <stdlib.h>
#include <stdio.h>
#include <sys/time.h>
#include "pvector.h"

void test_append_increases_length() ;
void test_append_adds_elements();
void test_append_branches();
void test_pvector_combine();
void test_pvector_equality();
void test_pvector_append_performance();
void test_pvector_combine_performance();
void test_pvector_rebalancing();

int main() {
    test_append_increases_length();
    test_append_adds_elements();
    test_append_branches();
    test_pvector_combine();
    test_pvector_equality();
    test_pvector_rebalancing();
    test_pvector_append_performance();
    test_pvector_combine_performance();
}

PVHead* make_pvector(uint32_t length) {
    PVHead *head = pvector_new();
    for (int i = 0; i < length; ++i) {
        PVHead *updated = pvector_append_uint16(head, i);
        pvector_free(head);
        head = updated;
    }
    return head;
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

    // Test joining two nodes where the resulting size is less than BRANCH
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
    pvector_free(a);
    pvector_free(b);
    pvector_free(res);

    // Test joining two leaf nodes where resulting size is greater than BRANCH
    a = make_pvector(20);
    b = make_pvector(20);
    res = pvector_combine_uint16(a, b);
    if (pvector_get_uint16(res, 32) != 12) {
        printf("expected res(32) == 12. Got %d\n", pvector_get_uint16(res, 32));
        exit(1);
    }
    pvector_free(a);
    pvector_free(b);
    pvector_free(res);

    // Test joining to nodes of differing sizes
    a = make_pvector(BRANCH / 2);
    b = make_pvector(BRANCH * 2);
    res = pvector_combine_uint16(a, b);
    for (uint32_t i = 0; i < BRANCH / 2; ++i) {
        if (pvector_get_uint16(res, i) != i) {
            printf("differing sizes 1: ");
            printf("expected res(%d) == %d. Got %d\n", i, i, pvector_get_uint16(res, i));
            exit(1);
        }
    }
    for (uint32_t i = 0; i < BRANCH * 2; ++i) {
        uint32_t ri = i + BRANCH / 2;
        if (pvector_get_uint16(res, ri) != i) {
            printf("differing sizes 2: ");
            printf("expected res(%d) == %d. Got %d\n", ri, i, pvector_get_uint16(res, ri));
            exit(1);
        }
    }
    pvector_free(a);
    pvector_free(b);
    pvector_free(res);

    // Test joining multiple small vectors from the left
    res = make_pvector(10);
    for (uint32_t i = 1; i < 200; ++i) {
        a = make_pvector(10);
        b = pvector_combine_uint16(res, a);
        pvector_free(a);
        pvector_free(res);
        res = b;
    }
    if (pvector_length(res) != 2000) {
        printf("expected new vector size to be 1000. Got %d\n", pvector_length(res));
        exit(1);
    }
    for (uint32_t i = 0; i < 10*200; ++i) {
        if (pvector_get_uint16(res, i) != i % 10) {
            printf("expected res(%d) == %d. Got %d\n", i, i % 10, pvector_get_uint16(res, i));
            exit(1);
        }
    }
    pvector_free(res);

    // Test joining multiple small vectors from the right
    res = make_pvector(10);
    for (uint32_t i = 1; i < 200; ++i) {
        a = make_pvector(10);
        b = pvector_combine_uint16(a, res);
        pvector_free(a);
        pvector_free(res);
        res = b;
    }
    if (pvector_length(res) != 2000) {
        printf("expected new vector size to be 1000. Got %d\n", pvector_length(res));
        exit(1);
    }
    pvector_free(res);

    printf("OK\n");
}

void test_pvector_equality() {
    printf("test_pvector_equality: ");
    PVHead *a = pvector_append_uint16(pvector_append_uint16(pvector_new(), 1), 2);
    PVHead *b = pvector_append_uint16(pvector_append_uint16(pvector_new(), 1), 2);
    if (!pvector_equals_uint16(a, b)) {
        printf("independent identical vectors were not equal\n");
        exit(1);
    }
    PVHead *aa = pvector_append_uint16(a, 4);
    PVHead *bb = pvector_append_uint16(b, 5);
    if (pvector_equals_uint16(aa, bb)) {
        printf("different vectors were equal\n");
        exit(1);
    }
    printf("OK\n");
}

void test_pvector_append_performance() {
    printf("test_pvector_append_performance: ");
    struct timeval tval_before, tval_after, tval_result;
    gettimeofday(&tval_before, NULL);
    PVHead *head = make_pvector(1000000);
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

void test_pvector_rebalancing() {
    printf("test_pvector_rebalancing: ");
    PVHead *a = make_pvector(BRANCH * BRANCH + 1);
    PVHead *b = make_pvector(BRANCH * BRANCH + 1);
    if (needs_rebalancing((PVNode*)a->node, (PVNode*)b->node)) {
        printf("needs_rebalancing returned true\n");
        exit(1);
    }
    pvector_free(a);
    pvector_free(b);
    printf("OK\n");
}