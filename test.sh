#!/bin/sh

rtest() {
    res=$(./run.sh $1)
    if [ "$res" != "$2" ]; then
        echo "FAILED $1: $res != $2"
        exit 1
    else
        echo "OK $1"
    fi
}

rtest "examples/features/integers.shi" "3"
rtest "examples/features/sequential_fn.shi" "9"
#rtest "examples/features/real_numbers.shi" "102334155.000000"
rtest "examples/features/multi_types.shi" "2"
rtest "examples/features/functions_as_types.shi" "13"
rtest "examples/features/closures.shi" "15"
#rtest "examples/features/memory_management.shi" "3"
rtest "examples/features/booleans.shi" "false"
rtest "examples/features/explicit_types.shi" "2.000000"
#rtest "examples/features/structures.shi" "40.094891"
#rtest "examples/features/strings.shi" "Lorem ipsum dolor sit amet, consectetur adipiscing elit. 日本語 ab is ab"
#rtest "examples/features/interfaces.shi" "12"

#rtest "examples/fibonacci.shi" "102334155"

#rtest "examples/euler/001.shi" "233168"
#rtest "examples/euler/002.shi" "4613732"
#rtest "examples/euler/003.shi" "6857"
#rtest "examples/euler/004.shi" "906609"
#rtest "examples/euler/005.shi" "232792560"
#rtest "examples/euler/006.shi" "25164150"
#rtest "examples/euler/007.shi" "104743"