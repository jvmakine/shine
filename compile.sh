#!/bin/sh

OPT_ARGS="--globaldce --die --dce --tailcallelim"

cat $1 | go run main.go | opt-9 -S $OPT_ARGS > tmp.ll
llvm-link-9 tmp.ll lib/runtime.ll > $2
rm tmp.ll
chmod u+x $2