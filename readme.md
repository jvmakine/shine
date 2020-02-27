# Shine frontend for LLVM

## Usage

Compile into executable
> echo "1 + 2 + 3 - 4 + 10" | go run main.go | opt-9 -S | llvm-as-9 > test

Interpret
> echo "1 + 2 + 3 - 4 + 10" | go run main.go | lli-9
