#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "pvector.h"

#define RRB_ERROR 1

uint8_t pvector_depth(PVHead *vector) {
   PVH* node = vector->node;
    if (node == 0) {
        return 0;
    }
    return node->depth;
}

uint32_t pvector_length(PVHead *vector) {
    return vector->size;
}

PVHead* pvector_new() {
    PVHead* head = heap_calloc(1, sizeof(PVHead));
    head->ref.count = 1;
    head->ref.type = MEM_PVECTOR;
    head->size = 0;
    head->node = 0;
    return head;
}

PVHead* pvector_from_node(PVH *node) {
    PVHead* head = pvector_new();
    head->node = node;
    head->size = node->size;
    node->refcount++;
    return head;
}

PVNode* pnode_new(uint8_t depth) {
    PVNode* node = heap_calloc(1, sizeof(PVNode));
    node->header.refcount = 1;
    node->header.depth = depth;
    node->header.size = 0;
    node->indextable = 0;
    return node;
}

PVH *pleaf_new(uint32_t leaf_size) {
    PVH* leaf = heap_calloc(1, leaf_size);
    leaf->refcount = 1;
    return leaf;
}

void pnode_free(PVH *n) {
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
                pnode_free(node->children[i]);
            }
        }
        if (node->indextable != 0) {
            free(node->indextable);
        }
    }
    free(n);
}

void pvector_free(PVHead *vector) {
    if (vector == 0 || vector->ref.count == CONSTANT_REF) {
        return;
    }
    if (vector->ref.count > 1) {
        vector->ref.count = vector->ref.count - 1;
        return;
    }
    if (vector->size > 0) {
        pnode_free(vector->node);
    }
    free(vector);
}

PVNode *copy_pnode(PVNode* node) {
    PVNode* res = heap_malloc(sizeof(PVNode));
    memcpy(res, node, sizeof(PVNode));
    res->header.refcount = 1;
    if (node->indextable != 0) {
        res->indextable = heap_malloc(BRANCH*sizeof(uint32_t));
        memcpy(res->indextable, node->indextable, BRANCH*sizeof(uint32_t));
    }
    return res;
}

void *copy_pleaf(PVH *leaf, uint32_t leaf_size) {
    PVH *res = heap_malloc(leaf_size);
    memcpy(res, leaf, leaf_size);
    res->refcount = 1;
    return res;
}

void pnode_set_child(PVNode *n, PVH *c, uint8_t index) {
    uint32_t os = 0;
    uint32_t ns = 0;
    if (n->children[index]) {
        os = n->children[index]->size;
    }
    n->children[index] = c;
    if (c) {
        ns = c->size;
        c->refcount++;
    }
    n->header.size += ns;
    n->header.size -= os;
}   

void increment_children_refcount(PVNode *node) {
    PVH **children = node->children;
    for (uint8_t i = 0; i < BRANCH; ++i) {
        // Can not break here, as sometimes we disable children to be replaced before the call
        if (children[i] != 0) {
            uint32_t rc = ((PVH*)children[i])->refcount;
            if (rc > 0) {
                ((PVH*)(children[i]))->refcount = rc + 1;
            }
        }
    }
}

void* pvector_get_leaf(PVHead *vector, uint32_t *inder_ptr) {
    uint32_t index = *inder_ptr;
    if (index >= vector->size) {
        fprintf(stderr, "pvector index out of bounds: got %d, size %d\n", index, vector->size);
        exit(1);
    }
    uint8_t depth = pvector_depth(vector);
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

uint16_t pvector_get_uint16(PVHead *vector, uint32_t index) {
    uint32_t index_v = index;
    PVLeaf_uint16 *leaf = pvector_get_leaf(vector, &index_v);
    return leaf->data[index_v & MASK];
}

uint8_t pvnode_right_child_index(PVNode *n) {
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

void *pvnode_right_child(PVNode *n) {
    return(n->children[pvnode_right_child_index(n)]);
}

void update_index_table(PVNode *n) {
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

uint8_t pvnode_branching(PVH* n) {
    if (n->depth > 0) {
        uint8_t i = 0;
        for(; i < BRANCH && ((PVNode*)n)->children[i]; ++i);
        return i;
    } else {
        return n->size;
    }
}

uint32_t branching_sum(PVNode* node) {
    uint32_t p = 0;
    uint8_t depth = node->header.depth;
    if (depth == 0) {
        return ((PVH*)node)->size;
    }
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (!node->children[i]) {
            break;
        }
        p += pvnode_branching(node->children[i]);
    }
    return p;
}

uint8_t needs_rebalancing(PVNode* left, PVNode* right) {
    if (left == 0 || right == 0) return 0;
    uint32_t p = branching_sum(left) + branching_sum(right);
    uint32_t a = pvnode_branching((PVH*)left) + pvnode_branching((PVH*)right);
    int e = a - ((p - 1) >> BITS) - 1;
    if (e > RRB_ERROR) {
        return 1;
    }
    return 0;
}

PVLeaf_uint16* combine_leaf_uint16(PVLeaf_uint16 *a, PVLeaf_uint16 *b, PVLeaf_uint16 **overflow) {
    if (a->header.size + b->header.size <= BRANCH) {
        if (overflow) {
            *overflow = 0;
        }
        PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pleaf_new(sizeof(PVLeaf_uint16));
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
        *overflow = (PVLeaf_uint16*)pleaf_new(sizeof(PVLeaf_uint16));
        PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pleaf_new(sizeof(PVLeaf_uint16));
        memcpy(leaf->data, a->data, a->header.size * sizeof(uint16_t));
        memcpy(leaf->data + a->header.size, b->data, (BRANCH - a->header.size) * sizeof(uint16_t));
        leaf->header.size = BRANCH;

        (*overflow)->header.size = overflow_size;
        memcpy((*overflow)->data, b->data + (BRANCH - a->header.size), overflow_size * sizeof(uint16_t));
        return leaf;
    }
}

PVH* join_nodes(PVH* left, PVH* right, PVH **overflow) {
    if (((PVLeaf_uint16*)left)->header.depth == 0 && ((PVLeaf_uint16*)right)->header.depth == 0) {
        return (PVH*)combine_leaf_uint16((PVLeaf_uint16*)left, (PVLeaf_uint16*)right, (PVLeaf_uint16**)overflow);
    } else {
        PVNode *a = (PVNode*)left;
        PVNode *b = (PVNode*)right;
        uint32_t asize = pvnode_branching(left);
        uint32_t bsize = pvnode_branching(right);
        uint8_t depth = ((PVNode*)a)->header.depth;
        if (((PVNode*)b)->header.depth != depth) {
            fprintf(stderr, "join error, depth mismatch\n");
            exit(1);
        }

        if (asize + bsize <= BRANCH) {
            if (overflow) {
                *overflow = 0;
            }
            PVNode *node = pnode_new(depth);
            memcpy(node->children, a->children, asize * sizeof(void*));
            memcpy(node->children + asize, b->children, bsize * sizeof(void*));

            update_index_table(node);
            increment_children_refcount(node);
            node->header.size = a->header.size + b->header.size;
            return (PVH*)node;
        } else {
            if (!overflow) {
                fprintf(stderr, "overflow required\n");
                exit(1);
            }
            uint32_t overflow_branches = (asize + bsize) - BRANCH;
            PVNode *ofn = pnode_new(depth);
            PVNode *node = pnode_new(depth);
            for (uint8_t i = 0; i < asize; ++i) {
                pnode_set_child(node, a->children[i], i);
            }
            for (uint8_t i = 0; i < (BRANCH - asize); ++i) {
                pnode_set_child(node, b->children[i], i + asize);
            }
             for (uint8_t i = 0; i < overflow_branches; ++i) {
                 pnode_set_child(ofn, b->children[i + (BRANCH - asize)], i);
            }
            update_index_table(node);
            update_index_table(ofn);
            *overflow = (PVH*)ofn;
            return (PVH*)node;
        }
    }
}

PVNode *pnode_replace_child(PVNode *node, uint8_t index, PVH* new_child) {
    if (node->header.depth <= 0) {
        fprintf(stderr, "no children in leaf nodes\n");
        exit(1);
    }
    PVNode *n = copy_pnode(node);
    uint32_t os = n->children[index]->size;
    uint32_t ns = new_child->size;
    n->children[index] = 0;
    increment_children_refcount(n);
    n->children[index] = new_child;
    n->header.size -= os;
    n->header.size += ns;
    update_index_table(n);
    return n;
}

PVNode *pnode_remove_child(PVNode *node, uint8_t index) {
    PVNode *copy = copy_pnode(node);
    if (!node->children[index]) {
        return copy;
    }
    uint32_t csize = copy->children[index]->size;
    copy->children[index] = 0;
    for (uint8_t i = index; i < BRANCH - 1; ++i) {
        copy->children[i] = copy->children[i + 1];
    }
    copy->children[BRANCH - 1] = 0;
    copy->header.size -= csize;
    increment_children_refcount(copy);
    update_index_table(copy);
    return copy;
}

uint8_t can_join(PVH *l, PVH *r) {
    if (l->depth == 0) {
        if (l->size + r->size <= BRANCH) {
            return 1;
        } else {
            return 0;
        }
    } else {
        if (pvnode_branching(l) + pvnode_branching(r) <= BRANCH) {
            return 1;
        } else {
            return 0;
        }
    }
}

PVNode *make_parent_node(PVH *child) {
    PVNode *n = pnode_new(child->depth + 1);
    n->header.size = child->size;
    n->children[0] = child;
    child->refcount++;
    return n;
}

PVHead* pvector_combine_uint16(PVHead *a, PVHead *b) {
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
        na = pvnode_right_child((PVNode*)na);
    }
    patha[ia] = na;

    while (nb->depth) {
        pathb[ib++] = nb;
        nb = ((PVNode*)nb)->children[0];
    }
    pathb[ib] = nb;

    PVH* l = patha[ia];
    PVH* r = pathb[ib];
    while (ia > 0 || ib > 0) {
        uint8_t balanced = 0;
        if (l && r && l->depth > 0 && needs_rebalancing((PVNode*)l, (PVNode*)r)) {
            PVNode *lr, *rr;
            balance_level((PVNode*)l, (PVNode*)r, &lr, &rr);
            l = (PVH*)lr;
            r = (PVH*)rr;
            balanced = 1;
        }
        if (l && r) {
            if (ia == 0) {
                if (can_join(l, r)) {
                    PVNode *n = (PVNode*)pathb[ib - 1];
                    r = (PVH*)pnode_replace_child(n, 0, join_nodes(l, r, 0));
                    l = 0;
                } else {
                    l = (PVH*)make_parent_node(l);
                }
            } 
            if (ib == 0) {
                if (can_join(l, r)) {
                    PVNode *n = (PVNode*)patha[ia - 1];
                    uint8_t index = pvnode_right_child_index(n);
                    l = (PVH*)pnode_replace_child(n, index, (PVH*)join_nodes(l, r, 0));
                    r = 0;
                } else {
                    r = (PVH*)make_parent_node(r);
                }
            } 
            if (ib > 0 && (r == pathb[ib] || balanced)) {
                r = pathb[ib - 1];
            }
            if (ia > 0 && (l == patha[ia] || balanced)) {
                l = patha[ia - 1];
            }
        } else if (l) {
            if (ia > 0) {
                PVNode *n = (PVNode*)patha[ia - 1];
                uint8_t index = pvnode_right_child_index(n);
                l = (PVH*)pnode_replace_child(n, index, l);
            } else {
                l = (PVH*)make_parent_node(l);
            }
            // Child was lost in balancing
            if (ib > 0) {
                r = pathb[ib - 1];
                r = (PVH*)pnode_remove_child((PVNode*)r, 0);
            }
        } else if (ib > 0) {
            if (ib > 0) {
                PVNode *n = (PVNode*)pathb[ib - 1];
                r = (PVH*)pnode_replace_child(n, 0, r);
            } else {
                r = (PVH*)make_parent_node(r);
            }
            // Child was lost in balancing
            if (ia > 0) {
                l = patha[ia - 1];
                l = (PVH*)pnode_remove_child((PVNode*)l, pvnode_right_child_index((PVNode*)l));
            }
        }
        if (ia > 0) ia--;
        if (ib > 0) ib--;
    }
    PVH *result;
    PVH *overflow;
    if (l && r) {
        result = join_nodes(l, r, &overflow);
        if (overflow) {
            PVNode *node = pnode_new(result->depth + 1);
            node->header.size = result->size + overflow->size;
            node->children[0] = result;
            node->children[1] = overflow;
            update_index_table(node);
            result = (PVH*)node;
        }
    } else if (l) {
        result = l;
    } else {
        result = r;
    }
    PVHead *head = pvector_from_node(result);
    return head;
}

void balance_level(PVNode* left, PVNode* right, PVNode **leftOut, PVNode **rightOut) {
    PVNode *new_left = copy_pnode(left);
    PVNode *new_right = copy_pnode(right);
    PVH *l = left->children[0];
    PVH *r = 0;
    uint32_t lsize = 0;
    uint32_t rsize = 0;
    uint8_t writeTo = 0;
    for (uint8_t i = 1; i < (BRANCH << 1); ++i) {
        if (l) {
            if (i < BRANCH) {
                r = new_left->children[i];
            } else {
                r = new_right->children[i - BRANCH];
            }
            if (r) {
                PVH *overflow;
                PVH *join = join_nodes(l, r, &overflow);
                if (writeTo > 0 && l) {
                    pnode_free(l);
                }
                l = overflow;
                if (writeTo < BRANCH) {
                    new_left->children[writeTo] = join;
                    lsize += join->size;
                } else {
                    new_right->children[writeTo - BRANCH] = join;
                    rsize += join->size;
                }
                writeTo++;
            }
        } else {
            if (i < BRANCH) {
                r = new_left->children[i];
            } else {
                r = new_right->children[i - BRANCH];
            }
            if (writeTo < BRANCH) {
                new_left->children[writeTo] = r;
                if (r) lsize += r->size;
            } else {
                new_right->children[writeTo - BRANCH] = r;
                if (r) rsize += r->size;
            }
            if (r) {
                r->refcount++;
                writeTo++;
            }
        }
    }
    if (l) {
        fprintf(stderr, "balance failed to compress a vector\n");
        exit(1);
    }
    new_right->children[BRANCH - 1] = 0;
    new_left->header.size = lsize;
    if (rsize == 0) {
       pnode_free((PVH*)new_right);
       new_right = 0;
    } else {
        new_right->header.size = rsize;
        update_index_table(new_right);
    }
    update_index_table(new_left);
    *leftOut = new_left;
    *rightOut = new_right;
}

uint8_t pvector_equals_uint16(PVHead *a, PVHead *b) {
    if (a->size != b->size) return 0;
    if (a == b || a->size == 0) return 1;
    for (uint32_t i = 0; i < a->size; ++i) {
        if (pvector_get_uint16(a, i) != pvector_get_uint16(b, i)) {
            return 0;
        }
    }
    return 1;
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    PVLeaf_uint16 *leaf = (PVLeaf_uint16*)pleaf_new(sizeof(PVLeaf_uint16));
    leaf->header.size = 1;
    leaf->data[0] = value;
    PVHead *head = pvector_from_node((PVH*)leaf);
    PVHead *result = pvector_combine_uint16(vector, head);
    pvector_free(head);
    return result;
}