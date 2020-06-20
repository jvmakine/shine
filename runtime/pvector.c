#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "pvector.h"

#define RRB_ERROR 2

uint8_t pvector_depth(PVHead *vector) {
   PVNode* node = vector->node;
    if (node == 0) {
        return 0;
    }
    return node->header.depth;
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

PVNode* pnode_new(uint8_t depth) {
    PVNode* node = heap_calloc(1, sizeof(PVNode));
    node->header.refcount = 1;
    node->indextable = 0;
    node->header.depth = depth;
    return node;
}

void *pleaf_new(uint32_t leaf_size) {
    PVLeaf_uint16* leaf = heap_calloc(1, leaf_size);
    leaf->header.refcount = 1;
    return (void*)leaf;
}

void pleaf_free(PVLeaf_uint16 *leaf) {
    // TODO: release references properly
    uint32_t rc = leaf->header.refcount;
    if (rc == 0) {
        return;
    }
    if (rc > 1) {
        leaf->header.refcount = rc - 1;
        return;
    }
    free(leaf);
}

void pnode_free(PVNode *node) {
    if (node == 0) {
        return;
    }
    uint32_t rc = node->header.refcount;
    if (rc == 0) {
        return;
    }
    if (rc > 1) {
        node->header.refcount = rc - 1;
        return;
    }
    for(int i = 0; i < BRANCH; ++i) {
        if (node->children[i] != 0) {
            if (node->header.depth > 1) {
                pnode_free(node->children[i]);
            } else {
                pleaf_free(node->children[i]);
            }
        }
    }
    if (node->indextable != 0) {
        free(node->indextable);
    }
    free(node);
}

void pvector_free(PVHead *vector) {
    if (vector == 0 || vector->ref.count == 0) {
        return;
    }
    if (vector->ref.count > 1) {
        vector->ref.count = vector->ref.count - 1;
        return;
    }
    if (vector->size > 0) {
        if (pvector_depth(vector) > 0) {
            pnode_free(vector->node);
        } else {
            pleaf_free(vector->node);
        }
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

void *copy_pleaf(void *leaf, uint32_t leaf_size) {
    void *res = heap_malloc(leaf_size);
    memcpy(res, leaf, leaf_size);
    ((PVLeaf_uint16*)res)->header.refcount = 1;
    return res;
}

void increment_children_refcount(PVNode *node) {
    void **children = node->children;
    for (uint8_t i = 0; i < BRANCH; ++i) {
        if (children[i] == 0) break;
        uint32_t rc = ((PVH*)children[i])->refcount;
        if (rc > 0) {
            ((PVH*)(children[i]))->refcount = rc + 1;
        }
    }
}

PVHead* pvector_append_leaf(PVHead *vector, uint32_t leaf_size, void **retval) {
    uint32_t old_size = vector->size;
    uint32_t new_size = old_size + 1;
    char new_node = 0;
    uint8_t depth = 0;
    uint32_t o = (old_size - 1) >> BITS;
    uint32_t n = (new_size - 1) >> BITS;
    if (old_size > 0) {
        while(n) {
            if (!o) { new_node = 1; }
            o = o >> BITS;
            n = n >> BITS;
            depth++;
        }
    } else {
        depth = 0;
        new_node = 1;
    }
    void *node = 0;
    if (new_node || vector->node == 0) {
        if (depth > 0) {
            node = pnode_new(depth);
            PVNode *vn = vector->node;
            if (vn) {
                ((PVNode*)node)->children[0] = vn;
                increment_children_refcount((PVNode*)node);
            }
        } else {
            node = pleaf_new(leaf_size);
        }
    } else {
        if (depth > 0) {
            node = copy_pnode(vector->node);
         } else {
             node = copy_pleaf(vector->node, leaf_size);
         }
    }

    PVHead *head = pvector_new();
    head->node = (PVNode*)node;
    head->size = new_size;

    while (depth) {
        uint8_t shift = depth*BITS;
        depth--;
        uint32_t key = (old_size >> shift) & MASK;
        void** children = ((PVNode*)node)->children;
        increment_children_refcount((PVNode*)node);
        if (depth) {
            if (children[key] == 0) {
                node = pnode_new(depth);
            } else {
                node = copy_pnode(children[key]);
            }
        } else {
            if (children[key] == 0) {
                node = pleaf_new(leaf_size);
            } else {
                node = copy_pleaf(children[key], leaf_size);
            }
        }
        children[key] = node;
    }
    ((PVLeaf_uint16*)node)->size++;
    *retval = node;
    return head;
}

void* pvector_get_leaf(PVHead *vector, uint32_t index) {
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
    return (void*)node;
}

uint16_t pvector_get_uint16(PVHead *vector, uint32_t index) {
    PVLeaf_uint16 *leaf = pvector_get_leaf(vector, index);
    return leaf->data[index & MASK];
}

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    PVLeaf_uint16 *leaf;
    PVHead *head = pvector_append_leaf(vector, sizeof(PVLeaf_uint16), (void*)&leaf);
    leaf->data[vector->size & MASK] = value;
    return head;
}

void *pvnode_right_child(PVNode *n) {
    int8_t i = BRANCH - 1;
    while(i >= 0 && n->children[i--] == 0);
    return(n->children[i + 1]);
}

uint32_t pvnode_size(PVNode *n) {
    uint32_t sum = 0;
    if(n->indextable != 0) {
        uint32_t *it = n->indextable;
        for(uint8_t i = 0; i < BRANCH; ++i) {
            sum += it[i];
        }
    } else {
        void **children = n->children;
        for(uint8_t i = 0; i < BRANCH; ++i) {
            // TODO: Optimise
             //if (i == BRANCH - 1 || children[i + 1] == 0) {
                 if (children[i] == 0) break;
                 if (n->header.depth == 1) {
                     sum += ((PVLeaf_uint16*)children[i])->size;
                 } else {
                     sum += pvnode_size((PVNode*)children[i]);
                 }
                 //break;
             /*} else {
                 uint32_t s = 1 << BITS;
                 uint8_t d = n->depth;
                 while(d > 0) {
                     s = s << BITS;
                     d--;
                 }
                 sum += s;
             }*/
        }
    }
    return sum;
}

void update_index_table(PVNode *n) {
    uint32_t indices[BRANCH];
    uint8_t depth = n->header.depth;
    uint8_t needed = 0;
    uint32_t sum = 0;
    uint32_t size = 0;
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (n->children[i]) {
            uint32_t full = 1 << BITS;
            for (uint8_t d = n->header.depth; d > 0; d--) {
                full = full << BITS;
            }
            if (size > 0 && size < full && i < BRANCH - 1 && n->children[i + 1] != 0) {
                needed = 1;
            }
            if (depth <= 1) {
                size = ((PVLeaf_uint16*)n->children[i])->size;
            } else {
                size = pvnode_size(n->children[i]);
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

uint8_t pvnode_branches(PVNode* n) {
    uint8_t i = 0;
    for(; i < BRANCH && n->children[i]; ++i);
    return i;
}

void balance_level(PVNode** left, PVNode** right) {
    fprintf(stderr, "BALANCE NOT IMPLEMENTED\n");
    exit(1);
}

uint32_t branching_sum(PVNode* node) {
    uint32_t p = 0;
    uint8_t depth = node->header.depth;
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (!node->children[i]) {
            break;
        }
        if (depth > 1) {
            p += pvnode_branches((PVNode*)node->children[i]);
        } else {
            p += ((PVLeaf_uint16*)node->children[i])->size;
        }
    }
    return p;
}

uint8_t needs_rebalancing(PVNode* left, PVNode* right) {
    uint32_t p = branching_sum(left) + branching_sum(right);
    uint32_t a = pvnode_branches(left) + pvnode_branches(right);
    uint32_t e = a - ((p - 1) >> BITS) - 1;
    if (e > RRB_ERROR) {
        return 1;
    }
    return 0;
}

void combine_level(PVNode** left, PVNode** right) {
    if (needs_rebalancing(*left, *right)) {
        balance_level(left, right);
    }
}

PVLeaf_uint16* combine_leaf_uint16(PVLeaf_uint16 *a, PVLeaf_uint16 *b, PVLeaf_uint16 **overflow) {
    if (a->size + b->size <= BRANCH) {
        *overflow = 0;
        PVLeaf_uint16 *leaf = pleaf_new(sizeof(PVLeaf_uint16));
        memcpy(leaf->data, a->data, a->size * sizeof(uint16_t));
        memcpy(leaf->data + a->size, b->data, b->size * sizeof(uint16_t));
        leaf->size = a->size + b->size;
        return leaf;
    } else {
        uint32_t overflow_size = (a->size + b->size) - BRANCH;
        *overflow = pleaf_new(sizeof(PVLeaf_uint16));
        PVLeaf_uint16 *leaf = pleaf_new(sizeof(PVLeaf_uint16));
        memcpy(leaf->data, a->data, a->size * sizeof(uint16_t));
        memcpy(leaf->data + a->size, b->data, (BRANCH - a->size) * sizeof(uint16_t));
        leaf->size = BRANCH;

        (*overflow)->size = overflow_size;
        memcpy((*overflow)->data, b->data + (BRANCH - a->size), overflow_size * sizeof(uint16_t));
        return leaf;
    }
}

void* join_nodes(void* left, void* right, void **overflow) {
    if (((PVLeaf_uint16*)left)->header.depth == 0 && ((PVLeaf_uint16*)right)->header.depth == 0) {
        return combine_leaf_uint16((PVLeaf_uint16*)left, (PVLeaf_uint16*)right, (PVLeaf_uint16**)overflow);
    } else {
        PVNode *a = (PVNode*)left;
        PVNode *b = (PVNode*)right;
        uint32_t asize = pvnode_branches(a);
        uint32_t bsize = pvnode_branches(b);
        uint8_t depth = ((PVNode*)a)->header.depth;
        if (((PVNode*)b)->header.depth != depth) {
            fprintf(stderr, "join error, depth mismatch\n");
            exit(1);
        }

        if (asize + bsize <= BRANCH) {
            *overflow = 0;
            PVNode *node = pnode_new(depth);
            memcpy(node->children, a->children, asize * sizeof(void*));
            memcpy(node->children + asize, b->children, bsize * sizeof(void*));

            update_index_table(node);
            increment_children_refcount(node);
            return node;
        } else {
            uint32_t overflow_size = (asize + bsize) - BRANCH;
            *overflow = pnode_new(depth);
            PVNode *node = pnode_new(depth);
            memcpy(node->children, a->children, asize * sizeof(void*));
            memcpy(node->children + asize, b->children, (BRANCH - asize) * sizeof(void*));
            memcpy(((PVNode*)(*overflow))->children, b->children + (BRANCH - asize), overflow_size * sizeof(void*));

            update_index_table(node);
            increment_children_refcount(node);
            update_index_table((PVNode*)(*overflow));
            increment_children_refcount((PVNode*)(*overflow));
            return node;
        }
    }
}

PVHead* pvector_combine_uint16(PVHead *a, PVHead *b) {
    // Construct the paths to the rightmost leaf of left value and leftmost leaf of right value
    void* patha[pvector_depth(a) + 1];
    void* pathb[pvector_depth(b) + 1];
    
    void *na = a->node;
    void *nb = b->node;
    uint8_t ia = 0;
    uint8_t ib = 0;

    while (((PVNode*)na)->header.depth > 0) {
        patha[ia++] = na;
        na = pvnode_right_child(na);
    }
    patha[ia] = na;

    while (((PVNode*)nb)->header.depth > 0) {
        pathb[ib++] = nb;
        nb = ((PVNode*)nb)->children[0];
    }
    pathb[ib] = nb;

    void *l = patha[ia];
    void *r = pathb[ib];
    while (ia > 0 && ib > 0) {
        ia--;
        ib--;
        l = patha[ia];
        r = pathb[ib];
        combine_level((PVNode**)&l, (PVNode**)&r);
    }
    if (ib == 0 && ia == 0) {
        void* overflow;
        PVNode *join = join_nodes(l, r, &overflow);
        PVHead *head = pvector_new();
        head->size = a->size + b->size;
        if (overflow == 0) {
            head->node = join;
           return head;
        } else {
            PVNode *node = pnode_new(pvector_depth(a) + 1);
            node->children[0] = join;
            node->children[1] = overflow;
            update_index_table(node);
            head->node = node;
            return head;
        }
    }
    fprintf(stderr, "UNEVEN MERGE NOT IMPLEMENTED!");
    exit(1);
}

uint8_t pleaf_equals(void *a, void *b, uint32_t leaf_size) {
    if (a == b) {
        return 1;
    }
    return !memcmp(a, b, leaf_size);
}

uint8_t pnode_equals(PVNode *a, PVNode *b, uint8_t leaf_size) {
    if (a == b) {
        return 1;
    }
    for(uint8_t i = 0; i < BRANCH; i++) {
        if (a->children[i] == 0) {
            return 1;
        }
        if (a->header.depth > 0) {
            if(!pnode_equals((PVNode*)a->children[i], (PVNode*)b->children[i], leaf_size)) {
                return 0;
            }
        } else {
            if(!pleaf_equals(a->children[i], b->children[i], leaf_size)) {
                return 0;
            }
        }
    }
    return 1;
}

uint8_t pvector_equals(PVHead *a, PVHead *b, uint32_t leaf_size) {
    if (a->size != b->size) {
        return 0;
    }
    if (a == b || a->size == 0) {
        return 1;
    }
    if (pvector_depth(a) > 0) {
        return pnode_equals(a->node, b->node, leaf_size);
    } else {
        return pleaf_equals(a->node, b->node, leaf_size);
    }
}