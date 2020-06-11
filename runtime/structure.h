typedef struct Structure {
    RefCount ref;
    uint16_t strucount;
} Structure;

typedef struct Closure {
    RefCount ref;
    void* fnptr;
    uint16_t strucount;
} Closure;

void free_structure(Structure *s);

void free_closure(Closure *s);