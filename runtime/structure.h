typedef struct Structure {
    RefCount ref;
    uint16_t clscount;
    uint16_t strucount;
} Structure;

void free_structure(Structure *s);