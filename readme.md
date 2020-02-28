# Shine frontend for LLVM

Implements a LLVM frontend for a simple language

## Usage

Compile into executable
> cat examples/functions.sh | go run main.go | opt-9 -S | llvm-as-9 > test

Interpret
> cat examples/functions.sh | go run main.go | lli-9
