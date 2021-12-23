.DEFAULT_GOAL := runtime

SRC = runtime/memory.c runtime/io.c runtime/pvector.c runtime/structure.c
LLS = memory.ll io.ll pvector.ll structure.ll

.PHONY: runtime
runtime:
	clang -S -emit-llvm $(SRC)
	llvm-link -S $(LLS) > lib/runtime.ll
	rm *.ll

test-pvector:
	clang runtime/pvector_test.c $(SRC) -o build/pvector_test
	build/pvector_test