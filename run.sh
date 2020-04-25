#!/bin/sh

OPT_ARGS="--globaldce --die --dce --tailcallelim"

cat $1 | go run main.go | opt-9 -S $OPT_ARGS | lli-9