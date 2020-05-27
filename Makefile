.DEFAULT_GOAL := build
build:
	clang-9 -shared -o lib/runtime.so runtime/memory.c runtime/io.c
	clang-9 -S -emit-llvm runtime/memory.c runtime/io.c
	llvm-link-9 -S memory.ll io.ll > lib/runtime.ll
	rm *.ll