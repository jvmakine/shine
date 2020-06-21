#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "pvector.h"

#define RRB_ERROR 2

void printf_uint16_node(PVH *node) {
    uint8_t d = node->depth;
    printf("{");
    if (d == 0) {
        printf("data:[");
        PVLeaf_uint16 *leaf = (PVLeaf_uint16*)node;
        for (uint8_t i = 0; i < node->size; ++i) {
            printf("%d,", leaf->data[i]);
        }
        
    } else {
        PVNode *n = (PVNode*)node;
        if (n->indextable) {
            printf("it:[");
            for (uint8_t i = 0; i < BRANCH; ++i) {
                if (n->indextable[i]) {
                    printf("%d,", n->indextable[i]);
                }
            }
            printf("],");
        }
        printf("children:[");
        for (uint8_t i = 0; i < BRANCH; ++i) {
            if (n->children[i]) {
                printf_uint16_node(n->children[i]);
                printf(",");
            }
        }
    }
    printf("]");
    printf("}");
}

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

PVNode* pnode_new(uint8_t depth) {
    PVNode* node = heap_calloc(1, sizeof(PVNode));
    node->header.refcount = 1;
    node->header.depth = depth;
    node->header.size = 0;
    node->indextable = 0;
    return node;
}

void *pleaf_new(uint32_t leaf_size) {
    PVLeaf_uint16* leaf = heap_calloc(1, leaf_size);
    leaf->header.refcount = 1;
    return (void*)leaf;
}

void pleaf_free(PVH *leaf) {
    // TODO: release references properly
    uint32_t rc = leaf->refcount;
    if (rc == 0) {
        return;
    }
    if (rc > 1) {
        leaf->refcount = rc - 1;
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
                pnode_free((PVNode*)node->children[i]);
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
            pnode_free((PVNode*)vector->node);
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

void *copy_pleaf(PVH *leaf, uint32_t leaf_size) {
    PVH *res = heap_malloc(leaf_size);
    memcpy(res, leaf, leaf_size);
    res->refcount = 1;
    return res;
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
    PVH *node = 0;
    if (new_node || vector->node == 0) {
        if (depth > 0) {
            node = (PVH*)pnode_new(depth);
            PVNode *vn = (PVNode*)vector->node;
            if (vn) {
                ((PVNode*)node)->children[0] = (PVH*)vn;
                 ((PVNode*)node)->header.size = vn->header.size;
                increment_children_refcount((PVNode*)node);
            }
        } else {
            node = pleaf_new(leaf_size);
        }
    } else {
        if (depth > 0) {
            node = (PVH*)copy_pnode((PVNode*)vector->node);
         } else {
             node = (PVH*)copy_pleaf(vector->node, leaf_size);
         }
    }

    PVHead *head = pvector_new();
    head->node = node;
    head->size = new_size;
    ((PVH*)node)->size++;
    while (depth) {
        uint8_t shift = depth*BITS;
        depth--;
        uint32_t key = (old_size >> shift) & MASK;
        PVH** children = ((PVNode*)node)->children;
        increment_children_refcount((PVNode*)node);
        if (depth) {
            if (children[key] == 0) {
                node = (PVH*)pnode_new(depth);
            } else {
                node = (PVH*)copy_pnode((PVNode*)children[key]);
            }
        } else {
            if (children[key] == 0) {
                node = (PVH*)pleaf_new(leaf_size);
            } else {
                node = (PVH*)copy_pleaf(children[key], leaf_size);
            }
        }
        children[key] = node;
        ((PVH*)node)->size++;
    }
    *retval = node;
    return head;
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

PVHead* pvector_append_uint16(PVHead *vector, uint16_t value) {
    PVLeaf_uint16 *leaf;
    PVHead *head = pvector_append_leaf(vector, sizeof(PVLeaf_uint16), (void*)&leaf);
    leaf->data[vector->size & MASK] = value;
    return head;
}

uint8_t pvnode_right_child_index(PVNode *n) {
    int8_t i = BRANCH - 1;
    while(i >= 0 && n->children[i--] == 0);
    return i + 1;
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
            uint32_t full = 1 << BITS;
            for (uint8_t d = n->header.depth; d > 0; d--) {
                full = full << BITS;
            }
            if (size > 0 && size < full && i < BRANCH - 1 && n->children[i + 1] != 0) {
                needed = 1;
            }
            size = ((PVH*)(n->children[i]))->size;
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
    if (depth == 0) {
        return ((PVH*)node)->size;
    }
    for(uint8_t i = 0; i < BRANCH; ++i) {
        if (!node->children[i]) {
            break;
        }
        if (depth > 1) {
            p += pvnode_branches((PVNode*)node->children[i]);
        } else {
            p += ((PVH*)node->children[i])->size;
        }
    }
    return p;
}

uint8_t needs_rebalancing(PVNode* left, PVNode* right) {
    if (left == 0 || right == 0) return 0;
    uint32_t p = branching_sum(left) + branching_sum(right);
    uint32_t a = pvnode_branches(left) + pvnode_branches(right);
    uint32_t e = a - ((p - 1) >> BITS) - 1;
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
        PVLeaf_uint16 *leaf = pleaf_new(sizeof(PVLeaf_uint16));
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
        *overflow = pleaf_new(sizeof(PVLeaf_uint16));
        PVLeaf_uint16 *leaf = pleaf_new(sizeof(PVLeaf_uint16));
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
        uint32_t asize = pvnode_branches(a);
        uint32_t bsize = pvnode_branches(b);
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
            *overflow = (PVH*)pnode_new(depth);
            PVNode *node = pnode_new(depth);
            uint32_t node_size = 0;
            uint32_t overflow_size = 0;
            for (uint8_t i = 0; i < asize; ++i) {
                node->children[i] = a->children[i];
                node_size +=  ((PVH*)a->children[i])->size;
            }
            for (uint8_t i = 0; i < (BRANCH - asize); ++i) {
                node->children[i + asize] = b->children[i];
                node_size +=  ((PVH*)b->children[i])->size;
            }
            node->header.size = node_size;
             for (uint8_t i = 0; i < overflow_branches; ++i) {
                ((PVNode*)(*overflow))->children[i] = b->children[i + (BRANCH - asize)];
                overflow_size +=  ((PVH*)b->children[i + (BRANCH - asize)])->size;
            }
            ((PVH*)*overflow)->size = overflow_size;

            update_index_table(node);
            increment_children_refcount(node);
            update_index_table((PVNode*)(*overflow));
            increment_children_refcount((PVNode*)(*overflow));
            return (PVH*)node;
        }
    }
}

PVNode *pnode_replace_child(PVNode *node, uint8_t index, PVH* new_child) {
    PVNode *n = copy_pnode(node);
    uint32_t os = n->children[index]->size;
    uint32_t ns = new_child->size;
    n->children[index] = 0;
    increment_children_refcount(n);
    n->children[index] = new_child;
    n->header.size -= os;
    n->header.size += ns;
    return n;
}

uint8_t can_join(PVH *l, PVH *r) {
    if (l->depth == 0) {
        if (l->size + r->size < BRANCH) {
            return 1;
        } else {
            return 0;
        }
    } else {
        PVNode *a = (PVNode*)l;
        PVNode *b = (PVNode*)r;
        if (pvnode_branches(a) + pvnode_branches(b) < BRANCH) {
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
    // Construct the paths to the rightmost leaf of left value and leftmost leaf of right value
    PVH* patha[pvector_depth(a) + 1];
    PVH* pathb[pvector_depth(b) + 1];
    
    PVH* na = a->node;
    PVH* nb = b->node;

    /*printf("l=");
    printf_uint16_node(na);
    printf(" r=");
    printf_uint16_node(nb);
    printf("\n");*/

    uint8_t ia = 0;
    uint8_t ib = 0;

    while (((PVNode*)na)->header.depth > 0) {
        patha[ia++] = na;
        na = pvnode_right_child((PVNode*)na);
    }
    patha[ia] = na;

    while (((PVNode*)nb)->header.depth > 0) {
        pathb[ib++] = nb;
        nb = ((PVNode*)nb)->children[0];
    }
    pathb[ib] = nb;

    PVH* l = patha[ia];
    PVH* r = pathb[ib];
    while (ia > 0 || ib > 0) {
        uint8_t balanced = 0;
        if (l && r && l->depth > 0 && needs_rebalancing((PVNode*)l, (PVNode*)r)) {
            balance_level((PVNode**)&l, (PVNode**)&r);
            balanced = 1;
        }
        
        if (l && r) {
            if (ia == 0) {
                if (can_join(l, r)) {
                    l = 0;
                    PVNode *n = (PVNode*)pathb[ib - 1];
                    r = (PVH*)pnode_replace_child(n, 0, join_nodes(l, r, 0));
                } else {
                    l = (PVH*)make_parent_node(l);
                }
            } 
            if (ib == 0) {
                if (can_join(l, r)) {
                    PVNode *n = (PVNode*)patha[ia - 1];
                    uint8_t index = pvnode_right_child_index(n);
                    PVNode *joined = (PVNode*)join_nodes(l, r, 0);
                    l = (PVH*)pnode_replace_child(n, index, (PVH*)joined);
                    r = 0;
                } else {
                    r = (PVH*)make_parent_node(r);
                }
            } 
            if (ib > 0 && r == pathb[ib]) {
                r = pathb[ib - 1];
            }
            if (ia > 0 && l == patha[ia]) {
                l = patha[ia - 1];
            }
        } else if (l) {
            PVNode *n = (PVNode*)patha[ia - 1];
            l = (PVH*)pnode_replace_child(n, pvnode_right_child_index(n), l);
        } else {
            PVNode *n = (PVNode*)pathb[ib - 1];
            r = (PVH*)pnode_replace_child(n, 0, r);
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
    PVHead *head = pvector_new();
    head->size = a->size + b->size;
    head->node = result;
    /*printf("result=");
    printf_uint16_node(result);
    printf("\n");*/
    return head;
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

uint8_t pvector_equals_uint16(PVHead *a, PVHead *b) {
    if (a->size != b->size) {
        return 0;
    }
    if (a == b || a->size == 0) {
        return 1;
    }
    for (uint32_t i = 0; i < a->size; ++i) {
        if (pvector_get_uint16(a, i) != pvector_get_uint16(b, i)) {
            return 0;
        }
    }
    return 1;
}