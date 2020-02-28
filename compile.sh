#!/bin/sh

OPT_ARGS="--globaldce --die --dce"

cat $1 | go run main.go | opt-9 -S $OPT_ARGS | llvm-as-9