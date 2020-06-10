.DEFAULT_GOAL := build
build:
	clang-9 -shared -o lib/runtime.so runtime/memory.c runtime/io.c runtime/strings.c runtime/pvector.c
	clang-9 -S -emit-llvm runtime/memory.c runtime/io.c runtime/strings.c runtime/pvector.c
	llvm-link-9 -S memory.ll io.ll strings.ll pvector.ll > lib/runtime.ll
	rm *.ll

test-pvector:
	clang-9 runtime/pvector_test.c runtime/pvector.c runtime/memory.c -o build/pvector_test
	build/pvector_test