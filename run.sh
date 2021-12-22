#!/bin/sh

OPT_ARGS="--globaldce --die --dce --tailcallelim"

cat $1 | go run main.go | opt -S $OPT_ARGS | lli -load lib/runtime.so
