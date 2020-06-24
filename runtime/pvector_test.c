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
void test_pvector_balancing_performance();

int main() {
    test_append_increases_length();
    test_append_adds_elements();
    test_append_branches();
    test_pvector_combine();
    test_pvector_equality();
    test_pvector_append_performance();
    test_pvector_combine_performance();
    test_pvector_balancing_performance();
}

PVHead* make_pvector(uint32_t length) {
    PVHead *head = pv_new();
    for (int i = 0; i < length; ++i) {
        PVHead *updated = pv_uint16_append(head, i);
        pv_free(head);
        head = updated;
    }
    return head;
}

void test_append_increases_length() {
    printf("test_append_increases_length: ");
    PVHead *head = pv_new();
    if (pv_length(head) != 0) {
        printf("expected new vector size to be 0. Got %d\n", pv_length(head));
        exit(1);
    }
    PVHead *appended = pv_uint16_append(pv_uint16_append(head, (uint16_t)5), (uint16_t)5);
    if (pv_length(head) != 0) {
        printf("expected vector size not to change after append. Got %d\n", pv_length(head));
        exit(1);
    }
    if (pv_length(appended) != 2) {
        printf("expected appended vector size to be 2. Got %d\n", pv_length(appended));
        exit(1);
    }
    printf("OK\n");
}

void test_append_adds_elements() {
     printf("test_append_adds_elements: ");
     PVHead *head = pv_new();
     PVHead *appended = pv_uint16_append(pv_uint16_append(head, (uint16_t)5), (uint16_t)8);
     if (pv_uint16_get(appended, 0) != 5) {
        printf("expected appended(0) == 5. Got %d\n", pv_uint16_get(appended, 0));
        exit(1);
    }
    if (pv_uint16_get(appended, 1) != 8) {
        printf("expected appended(0) == 8. Got %d\n", pv_uint16_get(appended, 2));
        exit(1);
    }
     printf("OK\n");
}

void test_append_branches() {
     printf("test_append_branches: ");
     PVHead *head = pv_uint16_append(pv_uint16_append(pv_new(), (uint16_t)5), (uint16_t)8);
     PVHead *b1 = pv_uint16_append(head, 12);
     PVHead *b2 = pv_uint16_append(head, 15);
     if (pv_uint16_get(b1, 2) != 12) {
        printf("expected b1(2) == 12. Got %d\n", pv_uint16_get(b1, 2));
        exit(1);
    }
    if (pv_uint16_get(b2, 2) != 15) {
        printf("expected b1(0) == 15. Got %d\n", pv_uint16_get(b2, 2) );
        exit(1);
    }
     printf("OK\n");
}

void test_pvector_combine() {
    printf("test_pvector_combine: ");

    // Test joining two nodes where the resulting size is less than BRANCH
    PVHead *a = pv_uint16_append(pv_uint16_append(pv_new(), 1), 2);
    PVHead *b = pv_uint16_append(pv_uint16_append(pv_new(), 3), 4);
    PVHead *res = pv_concatenate(a, b);
    if (pv_length(res) != 4) {
        printf("expected new vector size to be 4. Got %d\n", pv_length(res));
        exit(1);
    }
    if (pv_uint16_get(res, 1) != 2) {
        printf("expected res(1) == 2. Got %d\n", pv_uint16_get(res, 1));
        exit(1);
    }
    if (pv_uint16_get(res, 2) != 3) {
        printf("expected res(2) == 3. Got %d\n", pv_uint16_get(res, 2));
        exit(1);
    }
    if (pv_uint16_get(res, 0) != 1) {
        printf("expected res(0) == 1. Got %d\n", pv_uint16_get(res, 0));
        exit(1);
    }
    pv_free(a);
    pv_free(b);
    pv_free(res);

    // Test joining two leaf nodes where resulting size is greater than BRANCH
    a = make_pvector(20);
    b = make_pvector(20);
    res = pv_concatenate(a, b);
    if (pv_uint16_get(res, 32) != 12) {
        printf("expected res(32) == 12. Got %d\n", pv_uint16_get(res, 32));
        exit(1);
    }
    pv_free(a);
    pv_free(b);
    pv_free(res);

    // Test joining to nodes of differing sizes
    a = make_pvector(BRANCH / 2);
    b = make_pvector(BRANCH * 2);
    res = pv_concatenate(a, b);
    for (uint32_t i = 0; i < BRANCH / 2; ++i) {
        if (pv_uint16_get(res, i) != i) {
            printf("differing sizes 1: ");
            printf("expected res(%d) == %d. Got %d\n", i, i, pv_uint16_get(res, i));
            exit(1);
        }
    }
    for (uint32_t i = 0; i < BRANCH * 2; ++i) {
        uint32_t ri = i + BRANCH / 2;
        if (pv_uint16_get(res, ri) != i) {
            printf("differing sizes 2: ");
            printf("expected res(%d) == %d. Got %d\n", ri, i, pv_uint16_get(res, ri));
            exit(1);
        }
    }
    pv_free(a);
    pv_free(b);
    pv_free(res);

    // Test joining multiple small vectors from the left
    res = make_pvector(10);
    for (uint32_t i = 1; i < 200; ++i) {
        a = make_pvector(10);
        b = pv_concatenate(res, a);
        pv_free(a);
        pv_free(res);
        res = b;
    }
    if (pv_length(res) != 2000) {
        printf("expected new vector size to be 1000. Got %d\n", pv_length(res));
        exit(1);
    }
    for (uint32_t i = 0; i < 10*200; ++i) {
        if (pv_uint16_get(res, i) != i % 10) {
            printf("\nexpected res(%d) == %d. Got %d\n", i, i % 10, pv_uint16_get(res, i));
            exit(1);
        }
    }
    pv_free(res);

    // Test joining multiple small vectors from the right
    res = make_pvector(10);
    for (uint32_t i = 1; i < 200; ++i) {
        a = make_pvector(10);
        b = pv_concatenate(a, res);
        pv_free(a);
        pv_free(res);
        res = b;
    }
    if (pv_length(res) != 2000) {
        printf("expected new vector size to be 1000. Got %d\n", pv_length(res));
        exit(1);
    }

    for (uint32_t i = 0; i < 10*200; ++i) {
        if (pv_uint16_get(res, i) != i % 10) {
            printf("expected res(%d) == %d. Got %d\n", i, i % 10, pv_uint16_get(res, i));
            exit(1);
        }
    }
    pv_free(res);

    PVHead *res2;
    
    // Test joins that need rebalancing
    a = make_pvector(30);
    b = make_pvector(30);

    for (uint32_t i = 0; i < 6; ++i) {
        res = pv_concatenate(a, b);
        res2 = pv_concatenate(a, b);
        pv_free(a);
        pv_free(b);

        a = res;
        b = res2;
    }
    res = pv_concatenate(a, b);
    pv_free(a);
    pv_free(b);

    if (pv_length(res) != 30 << 7) {
        printf("expected new vector size to be 30 << 7. Got %d\n", pv_length(res));
        exit(1);
    }

    for (uint32_t i = 0; i < pv_length(res); ++i) {
        if (pv_uint16_get(res, i) != i % 30) {
            printf("expected res(%d) == %d. Got %d\n", i, i % 30, pv_uint16_get(res, i));
            exit(1);
        }
    }

    printf("OK\n");
}

void test_pvector_equality() {
    printf("test_pvector_equality: ");
    PVHead *a = pv_uint16_append(pv_uint16_append(pv_new(), 1), 2);
    PVHead *b = pv_uint16_append(pv_uint16_append(pv_new(), 1), 2);
    if (!pv_uint16_equals(a, b)) {
        printf("independent identical vectors were not equal\n");
        exit(1);
    }
    PVHead *aa = pv_uint16_append(a, 4);
    PVHead *bb = pv_uint16_append(b, 5);
    if (pv_uint16_equals(aa, bb)) {
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
    printf("construct %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);
    
    gettimeofday(&tval_before, NULL);
    for (int i = 0; i < 1000000; ++i) {
        uint16_t val = pv_uint16_get(head, i);
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
    PVHead *v1 = pv_new();
    PVHead *v2 = pv_new();
    for (int i = 0; i < 500000; ++i) {
        PVHead *updated = pv_uint16_append(v1, i);
        pv_free(v1);
        v1 = updated;
        updated = pv_uint16_append(v2,500000 + i);
        pv_free(v2);
        v2 = updated;
    }
    struct timeval tval_before, tval_after, tval_result;
    gettimeofday(&tval_before, NULL);
    PVHead *combined = pv_concatenate(v1, v2);
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);

    pv_free(v1);
    pv_free(v2);
    if (combined->size != 1000000) {
        printf("expected combined->size == %d. Got %d\n", 1000000, combined->size);
        exit(1);
    }
    printf("construct %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);

    gettimeofday(&tval_before, NULL);
    for (int i = 0; i < 1000000; ++i) {
        uint16_t val = pv_uint16_get(combined, i);
        uint16_t exp = i;
        if (val != exp) {
            printf("expected combined(%d) == %d. Got %d\n", i, exp, val);
            exit(1);
        }
    }
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("index %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);
    pv_free(combined);
    printf("OK\n");
}

void test_pvector_balancing_performance() {
    printf("test_pvector_balancing_performance: ");
    PVHead *a = make_pvector(7);
    PVHead *b = make_pvector(7);
    PVHead *res;
    PVHead *res2;

    struct timeval tval_before, tval_after, tval_result;
    gettimeofday(&tval_before, NULL);

    for (uint32_t i = 0; i < 16; ++i) {
        PVHead *res = pv_concatenate(a, b);
        res2 = pv_concatenate(a, b);
        pv_free(a);
        pv_free(b);

        a = res;
        b = res2;
    }
    res = pv_concatenate(a, b);
    pv_free(a);
    pv_free(b);

    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("construct %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);

    gettimeofday(&tval_before, NULL);
    for (uint32_t i = 0; i < pv_length(res); ++i) {
        if (pv_uint16_get(res, i) != i % 7) {
            printf("expected res(%d) == %d. Got %d\n", i, i % 15, pv_uint16_get(res, i));
            exit(1);
        }
    }
    gettimeofday(&tval_after, NULL);
    timersub(&tval_after, &tval_before, &tval_result);
    printf("index %ld.%06lds ", (long int)tval_result.tv_sec, (long int)tval_result.tv_usec);

    printf("OK\n");
}