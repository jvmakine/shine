#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "pvector.h"

#define RRB_ERROR 1

uint8_t pv_depth(PVHead *vector) {
   PVH* node = vector->node;
    if (node == 0) {
        return 0;
    }
    return node->depth;
}

uint32_t pv_length(PVHead *vector) {
    return vector->size;
}

void pn_incr_ref(PVH *p) {
    if (p->refcount != CONSTANT_REF) {
        p->refcount++;
    }
}

PVHead* pv_new() {
    PVHead* head = heap_calloc(1, sizeof(PVHead));
    head->ref.count = 1;
    head->ref.type = MEM_PVECTOR;
    return head;
}

PVHead* pv_construct(PVH *node) {
    PVHead* head = pv_new();
    head->node = node;
    head->size = node->size;
    pn_incr_ref(node);
    return head;
}

PVNode* pn_new(uint8_t depth) {
    PVNode* node = heap_calloc(1, sizeof(PVNode));
    node->header.depth = depth;
    return node;
}

PVH* pl_new(uint32_t leaf_size) {
    return heap_calloc(1, leaf_size);
}

void pn_free(PVH *n) {
    if (n == 0) {
        return;
    }
    uint32_t rc = n->refcount;
    if (rc == CONSTANT_REF) {
        return;
    }
    if (rc > 1) {
        n->refcount = rc - 1;
        return;
    }
    if (n->depth > 0) {
        PVNode *node = (PVNode*)n;
        for(int i = 0; i < BRANCH; ++i) {
            if (node->children[i] != 0) {
                pn_free(node->children[i]);
            }
        }
        if (node->indextable != 0) {
            free(node->indextable);
        }
    }
    free(n);
}

void pv_free(PVHead *vector) {
    if (vector == 0 || vector->ref.count == CONSTANT_REF) {
        return;
    }
    if (vector->ref.count > 1) {
        vector->ref.count = vector->ref.count - 1;
        return;
    }
    if (vector->node) {
        pn_free(vector->node);
    }
    free(vector);
}

void pn_increment_children_ref(PVNode *node) {
    PVH **children = node->children;
    for (uint8_t i = 0; i < BRANCH; ++i) {
        // Can not break here, as sometimes we disable children to be replaced before the call
        if (children[i] != 0) {
            pn_incr_ref(children[i]);
        }
    }
}

PVNode *pn_copy(PVNode* node) {
    PVNode* res = heap_malloc(sizeof(PVNode));
    memcpy(res, node, sizeof(PVNode));
    if (node->indextable != 0) {
        res->indextable = heap_malloc(BRANCH*sizeof(uint32_t));
        memcpy(res->indextable, node->indextable, BRANCH*sizeof(uint32_t));
    }
    pn_increment_children_ref(res);
    return res;
}

void *pl_copy(PVH *leaf, uint32_t leaf_size) {
    PVH *res = heap_malloc(leaf_size);
    memcpy(res, leaf, leaf_size);
    return res;
}

void pn_set_child(PVNode *n, PVH *c, uint8_t index) {
    uint32_t os = 0;
    uint32_t ns = 0;
    PVH *old_child = n->children[index];
    if (old_child) {
        os = old_child->size;
        pn_free(old_child);
    }
    n->children[index] = c;
    if (c) {
        ns = c->size;
        pn_incr_ref(c);
    }
    n->header.size += ns;
    n->header.size -= os;
}   

void* pv_get_leaf(PVHead *vector, uint32_t *inder_ptr) {
    uint32_t index = *inder_ptr;
    if (index >= vector->size) {
        fprintf(stderr, "pvector index out of bounds: got %d, size %d\n", index, vector->size);
        exit(1);
    }
    uint8_t depth = pv_depth(vector);
    void *node = vector->node;
    uint32_t *it = ((PVNode*)node)->indextable;

    // If the index table is in use, we need to use it to adjust the index
    while (depth && it != 0) {
        uint8_t r = BRANCH - 1;
        uint8_t l = 0;
        // TODO: Optimise finding the right edge
        while(r > 0 && it[r] == 0) { r--; }
        // Binary search for the right child
        while(r > l) {
            uint8_t mid = l + ((r - l) >> 1);
            uint32_t mv = it[mid];
            if (mv == index) {
                r = mid + 1;
                break;
            } else if (mv < index) {
                l = mid + 1;
            } else {
                r = mid;
            }
        }
        if (r > 0) {
            index -= it[r - 1];
        }
        node = ((PVNode*)node)->children[r];
        it = ((PVNode*)node)->indextable;
        depth--;
    }

    // Fully balanced subtree does not need the index lookup
    while (depth) {
        uint32_t key = (index >> depth*BITS) & MASK;
        depth--;
        node = ((PVNode*)node)->children[key];
    }
    *inder_ptr = index;
    return (void*)node;
}

uint16_t pv_uint16_get(PVHead *vector, uint32_t index) {
    uint32_t index_v = index;
    PVLeaf_uint16 *leaf = pv_get_leaf(vector, &index_v);
    return leaf->data[index_v & MASK];
}

uint8_t pn_right_child_index(PVNode *n) {
    if (n->indextable) {
        int8_t i = BRANCH;
        while(i >= 0 && n->children[--i] == 0);
        return i;
    } else {
        uint8_t depth = n->header.depth;
        uint32_t s = 1;
        for(uint8_t d = depth; d > 0; d--) {
            s = s << BITS;
        }
        uint8_t i = 0;
        uint32_t size = n->header.size;
        while(size > s) {
            size -= s;
            i++;
        }
        return i;
    }
}

void *pn_right_child(PVNode *n) {
    return(n->children[pn_right_child_index(n)]);
}

void pn_update_index_table(PVNode *n) {
    uint32_t indices[BRANCH];
    uint8_t depth = n->header.depth;
    uint8_t needed = 0;
    uint32_t sum = 0;
    uint32_t size = 0;
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (n->children[i]) {
            uint32_t full = 1;
            for (uint8_t d = n->header.depth; d > 0; d--) {
                full = full << BITS;
            }
            size = ((PVH*)(n->children[i]))->size;
            if (size > 0 && size < full && i < BRANCH - 1 && n->children[i + 1] != 0) {
                needed = 1;
            }
            sum += size;
            indices[i] = sum;
        } else {
            indices[i] = 0;
        }
    }
    if (n->indextable) {
        free(n->indextable);
        n->indextable = 0;
    }
    if (needed) {
        n->indextable = heap_malloc(BRANCH * sizeof(uint32_t));
        memcpy(n->indextable, indices, BRANCH * sizeof(uint32_t));
    }
}

uint8_t pn_branch_count(PVH* n) {
    if (n->depth > 0) {
        uint8_t i = 0;
        for(; i < BRANCH && ((PVNode*)n)->children[i]; ++i);
        return i;
    } else {
        return n->size;
    }
}

uint32_t pn_branch_sum(PVNode* node) {
    uint32_t p = 0;
    uint8_t depth = node->header.depth;
    if (depth == 0) {
        return ((PVH*)node)->size;
    }
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (!node->children[i]) {
            break;
        }
        p += pn_branch_count(node->children[i]);
    }
    return p;
}

uint8_t pn_needs_rebalancing(PVNode* left, PVNode* right) {
    if (left == 0 || right == 0) return 0;
    uint32_t p = pn_branch_sum(left) + pn_branch_sum(right);
    uint32_t a = pn_branch_count((PVH*)left) + pn_branch_count((PVH*)right);
    int e = a - ((p - 1) >> BITS) - 1;
    if (e > RRB_ERROR) {
        return 1;
    }
    return 0;
}

PVLeaf_uint16* pl_16_concatenate(PVLeaf_uint16 *a, PVLeaf_uint16 *b, PVLeaf_uint16 **overflow) {
    if (a->header.size + b->header.size <= BRANCH) {
        if (overflow) {
            *overflow = 0;
        }
        PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pl_new(sizeof(PVLeaf_uint16));
        memcpy(leaf->data, a->data, a->header.size * sizeof(uint16_t));
        memcpy(leaf->data + a->header.size, b->data, b->header.size * sizeof(uint16_t));
        leaf->header.size = a->header.size + b->header.size;
        return leaf;
    } else {
        if (!overflow) {
            fprintf(stderr, "overflow required\n");
            exit(1);
        }
        uint32_t overflow_size = (a->header.size + b->header.size) - BRANCH;
        *overflow = (PVLeaf_uint16*)pl_new(sizeof(PVLeaf_uint16));
        PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pl_new(sizeof(PVLeaf_uint16));
        memcpy(leaf->data, a->data, a->header.size * sizeof(uint16_t));
        memcpy(leaf->data + a->header.size, b->data, (BRANCH - a->header.size) * sizeof(uint16_t));
        leaf->header.size = BRANCH;

        (*overflow)->header.size = overflow_size;
        memcpy((*overflow)->data, b->data + (BRANCH - a->header.size), overflow_size * sizeof(uint16_t));
        return leaf;
    }
}

PVH* pn_join_nodes(PVH* left, PVH* right, PVH **overflow) {
    if (((PVLeaf_uint16*)left)->header.depth == 0 && ((PVLeaf_uint16*)right)->header.depth == 0) {
        return (PVH*)pl_16_concatenate((PVLeaf_uint16*)left, (PVLeaf_uint16*)right, (PVLeaf_uint16**)overflow);
    } else {
        PVNode *a = (PVNode*)left;
        PVNode *b = (PVNode*)right;
        uint32_t asize = pn_branch_count(left);
        uint32_t bsize = pn_branch_count(right);
        uint8_t depth = ((PVNode*)a)->header.depth;
        if (((PVNode*)b)->header.depth != depth) {
            fprintf(stderr, "join error, depth mismatch\n");
            exit(1);
        }

        if (asize + bsize <= BRANCH) {
            if (overflow) {
                *overflow = 0;
            }
            PVNode *node = pn_new(depth);
            memcpy(node->children, a->children, asize * sizeof(void*));
            memcpy(node->children + asize, b->children, bsize * sizeof(void*));

            pn_update_index_table(node);
            pn_increment_children_ref(node);
            node->header.size = a->header.size + b->header.size;
            return (PVH*)node;
        } else {
            if (!overflow) {
                fprintf(stderr, "overflow required\n");
                exit(1);
            }
            uint32_t overflow_branches = (asize + bsize) - BRANCH;
            PVNode *ofn = pn_new(depth);
            PVNode *node = pn_new(depth);
            for (uint8_t i = 0; i < asize; ++i) {
                pn_set_child(node, a->children[i], i);
            }
            for (uint8_t i = 0; i < (BRANCH - asize); ++i) {
                pn_set_child(node, b->children[i], i + asize);
            }
             for (uint8_t i = 0; i < overflow_branches; ++i) {
                 pn_set_child(ofn, b->children[i + (BRANCH - asize)], i);
            }
            pn_update_index_table(node);
            pn_update_index_table(ofn);
            *overflow = (PVH*)ofn;
            return (PVH*)node;
        }
    }
}

PVNode *pn_replace_child(PVNode *node, uint8_t index, PVH* new_child) {
    PVNode *n = pn_copy(node);
    pn_set_child(n, new_child, index);
    pn_update_index_table(n);
    return n;
}

PVNode *pn_remove_child(PVNode *node, uint8_t index) {
    PVNode *copy = pn_copy(node);
    PVH *child = node->children[index];
    if (!child) {
        return copy;
    }
    uint32_t csize = child->size;
    pn_free(child);
    copy->children[index] = 0;
    for (uint8_t i = index; i < BRANCH - 1; ++i) {
        copy->children[i] = copy->children[i + 1];
    }
    copy->children[BRANCH - 1] = 0;
    copy->header.size -= csize;
    pn_update_index_table(copy);
    return copy;
}

uint8_t pn_fits_into_one_node(PVH *l, PVH *r) {
    if (l->depth == 0) {
        if (l->size + r->size <= BRANCH) {
            return 1;
        } else {
            return 0;
        }
    } else {
        if (pn_branch_count(l) + pn_branch_count(r) <= BRANCH) {
            return 1;
        } else {
            return 0;
        }
    }
}

PVNode *pn_make_parent(PVH *child) {
    PVNode *n = pn_new(child->depth + 1);
    pn_set_child(n, child, 0);
    return n;
}

void rassign(PVH **to, PVH *from) {
    if (*to) {
        pn_free(*to);
    }
    if (from) {
        pn_incr_ref(from);
    }
    *to = from;
}

PVHead* pv_concatenate(PVHead *a, PVHead *b) {
    PVH* na = a->node;
    PVH* nb = b->node;

    if (!na) {
        b->ref.count++;
        return b;
    }
    if (!nb) {
        a->ref.count++;
        return a;
    }

    // Construct the paths to the rightmost leaf of left value and leftmost leaf of right value
    PVH* patha[na->depth + 1];
    PVH* pathb[nb->depth + 1];

    uint8_t ia = 0;
    uint8_t ib = 0;

    while (na->depth) {
        patha[ia++] = na;
        na = pn_right_child((PVNode*)na);
    }
    patha[ia] = na;

    while (nb->depth) {
        pathb[ib++] = nb;
        nb = ((PVNode*)nb)->children[0];
    }
    pathb[ib] = nb;

    PVH* l = 0;
    PVH* r = 0;
    rassign(&l, patha[ia]);
    rassign(&r, pathb[ib]);
    while (ia > 0 || ib > 0) {
        uint8_t balanced = 0;
        if (l && r && l->depth > 0 && pn_needs_rebalancing((PVNode*)l, (PVNode*)r)) {
            PVNode *lr, *rr;
            pn_balance_level((PVNode*)l, (PVNode*)r, &lr, &rr);
            rassign(&l, (PVH*)lr);
            rassign(&r, (PVH*)rr);
            balanced = 1;
        }
        if (l && r) {
            if (ia == 0) {
                if (pn_fits_into_one_node(l, r)) {
                    PVNode *n = (PVNode*)pathb[ib - 1];
                    PVH *join = pn_join_nodes(l, r, 0);
                    rassign(&l, 0);
                    rassign(&r, (PVH*)pn_replace_child(n, 0, join));
                } else {
                    rassign(&l, (PVH*)pn_make_parent(l));
                }
            } 
            if (ib == 0) {
                if (pn_fits_into_one_node(l, r)) {
                    PVNode *n = (PVNode*)patha[ia - 1];
                    uint8_t index = pn_right_child_index(n);
                    PVH *join = pn_join_nodes(l, r, 0);
                    rassign(&l, (PVH*)pn_replace_child(n, index, join));
                    rassign(&r, 0);
                } else {
                    rassign(&r, (PVH*)pn_make_parent(r));
                }
            } 
            if (ib > 0 && (r == pathb[ib] || balanced)) {
                rassign(&r, pathb[ib - 1]);
            }
            if (ia > 0 && (l == patha[ia] || balanced)) {
                rassign(&l, patha[ia - 1]);
            }
        } else if (l) {
            if (ia > 0) {
                PVNode *n = (PVNode*)patha[ia - 1];
                uint8_t index = pn_right_child_index(n);
                PVNode *nn = pn_replace_child(n, index, l);
                rassign(&l, (PVH*)nn);
            } else {
                rassign(&l, (PVH*)pn_make_parent(l));
            }
            // Child was lost in balancing
            if (ib > 0) {
                PVNode *n = (PVNode*)pathb[ib - 1];
                PVNode *nn = pn_remove_child(n, 0);
                rassign(&r, (PVH*)nn);
            }
        } else if (ib > 0) {
            if (ib > 0) {
                PVNode *n = (PVNode*)pathb[ib - 1];
                PVNode *nn = pn_replace_child(n, 0, r);
                rassign(&r, (PVH*)nn);
            } else {
                rassign(&r, (PVH*)pn_make_parent(r));
            }
            // Child was lost in balancing
            if (ia > 0) {
                PVNode *n = (PVNode*)patha[ia - 1];
                PVNode *nn = pn_remove_child(n, pn_right_child_index(n));
                rassign(&l, (PVH*)nn);
            }
        }
        if (ia > 0) ia--;
        if (ib > 0) ib--;
    }
    PVH *result;
    PVH *overflow;
    if (l && r) {
        result = pn_join_nodes(l, r, &overflow);
        pn_free(l);
        pn_free(r);
        if (overflow) {
            PVNode *node = pn_new(result->depth + 1);
            pn_set_child(node, result, 0);
            pn_set_child(node, overflow, 1);
            pn_update_index_table(node);
            result = (PVH*)node;
        }
    } else if (l) {
        result = l;
    } else {
        result = r;
    }
    return pv_construct(result);
}

void pn_balance_level(PVNode* left, PVNode* right, PVNode **leftOut, PVNode **rightOut) {
    PVNode *new_left = pn_new(left->header.depth);
    PVNode *new_right = pn_new(right->header.depth);
    PVH *l = left->children[0];
    PVH *r = 0;
    uint8_t writeTo = 0;
    for (uint8_t i = 1; i < (BRANCH << 1); ++i) {
        if (l) {
            if (i < BRANCH) {
                r = left->children[i];
            } else {
                r = right->children[i - BRANCH];
            }
            if (r) {
                PVH *overflow;
                PVH *join = pn_join_nodes(l, r, &overflow);
                if (writeTo > 0 && l) {
                    pn_free(l);
                }
                l = overflow;
                if (writeTo < BRANCH) {
                    pn_set_child(new_left, join, writeTo);
                } else {
                    pn_set_child(new_right, join, writeTo - BRANCH);
                }
                writeTo++;
            }
        } else {
            if (i < BRANCH) {
                r = left->children[i];
            } else {
                r = right->children[i - BRANCH];
            }
            if (writeTo < BRANCH) {
                pn_set_child(new_left, r, writeTo);
            } else {
                pn_set_child(new_right, r, writeTo - BRANCH);
            }
            if (r) writeTo++;
        }
    }
    if (l) {
        fprintf(stderr, "balance failed to compress a vector\n");
        exit(1);
    }
    if (new_right->header.size == 0) {
       pn_free((PVH*)new_right);
       new_right = 0;
    } else {
        pn_update_index_table(new_right);
    }
    pn_update_index_table(new_left);
    *leftOut = new_left;
    *rightOut = new_right;
}

uint8_t pv_uint16_equals(PVHead *a, PVHead *b) {
    if (a->size != b->size) return 0;
    if (a == b || a->size == 0) return 1;
    for (uint32_t i = 0; i < a->size; ++i) {
        if (pv_uint16_get(a, i) != pv_uint16_get(b, i)) {
            return 0;
        }
    }
    return 1;
}

PVHead* pv_uint16_append(PVHead *vector, uint16_t value) {
    PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pl_new(sizeof(PVLeaf_uint16));
    leaf->header.size = 1;
    leaf->data[0] = value;
    PVHead *head = pv_construct((PVH*)leaf);
    PVHead *result = pv_concatenate(vector, head);
    pv_free(head);
    return result;
}