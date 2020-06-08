
struct String {
    uint32_t refcount;
    uint16_t clscount;
    uint16_t strucount;
    struct String *next;
    char *base;
};