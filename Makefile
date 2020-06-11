.DEFAULT_GOAL := runtime
.PHONY: runtime

SRC = runtime/memory.c runtime/io.c runtime/strings.c runtime/pvector.c runtime/structure.c
LLS = memory.ll io.ll strings.ll pvector.ll structure.ll

runtime:
	clang-9 -shared -fPIC -o lib/runtime.so $(SRC)
	clang-9 -S -emit-llvm $(SRC)
	llvm-link-9 -S $(LLS) > lib/runtime.ll
	rm *.ll

test-pvector:
	clang-9 runtime/pvector_test.c $(SRC) -o build/pvector_test
	build/pvector_test