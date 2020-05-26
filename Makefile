.DEFAULT_GOAL := build
build:
	clang-9 -shared -o lib/runtime.so runtime/memory.c
	clang-9 -S -emit-llvm runtime/memory.c
	mv memory.ll lib/runtime.ll