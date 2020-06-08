# Shine frontend for LLVM

An LLVM frontend for a simple language

## Pre-requirements

LLVM 9 needs to be installed to compile or interpret the IR generated by this compiler.

## Building the runtime

To build the runtime for the language, call

```
make
```

## Usage

Interpret
> ./run.sh examples/fibonacci.shi

Compile into executable and run it
> ./compile.sh examples/fibonacci.shi fibonacci

> ./fibonacci

## Language

Every script needs to end in an expression. The value of this expression is printed as the output of the program.

Example script:
```
fibonacci = (n) => if (n <= 2) 1 else fibonacci(n-2) + fibonacci(n-1)

fibonacci(40)
```
This prints out `102334155`

### Types

The language supports full type inference. Explicitly declaring the types is completely optional.

#### Primitive types

The primitive types supported are

integer values supporting arithmetic operations
```
a: int = 1 + 2
```

real numbers supporting arithmetic operations
```
a: real = 1.0 + 2.0
```

boolean values that can be used as conditions in if expressions. Supports and and or operations
```
a: bool = true || false
if (a) 1 else 2
```

strings can be used as output values
```
a: string = "Hello world!"
a
```

#### Functions

Functions are composed or parameter types and the return type. Functions can be used as values themselves, and can capture values from the point of definition in their closures.
```
c = 3
f = (a, b) => a + b + c
f(1, 2)
```
prints out 6

Functions don't need to be named and anonymous functions can be defined within an expression.

```
op = (a, b, f) => f(a, b)
op(1, 2, (x, y) => { x + y })
```
prints out 3

Function parameters can be declared is several blocks to support currying.

```
f = (a) => (b) => a + b
x = f(1)
x(2)
```
prints out 3

#### Structures

Structures are used to combine several valuesinto a single value. These values can then be accessed with the dot operator. The structure is declared as a function parameter list without a body.

```
Person = (age: int, height, weight)
p = Person(38, 1.73, 60.0)
p.weight
```

prints out 60.0

Explicitly specifying that an expression is of a specific structure, the structure name can be used.
```
getAge = (p: Person) => p.age
```
If explicit structure is not defined, functions can take any structure as in input with required fields.
```
Person = (age, height)
Animal = (age, weight)

getAge = (x) => x.age

getAge(Person(30, 1.70)) + getAge(Animal(5, 0.3))
```
prints out 35
