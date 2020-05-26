.DEFAULT_GOAL := build
build:
	clang-9 -shared -o lib/runtime.so runtime/memory.c
	clang-9 -c -o lib/runtime.o runtime/memory.c