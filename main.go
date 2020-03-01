package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString(0)

	ast, err := grammar.Parse(text)
	if err != nil {
		panic(err)
	}

	module, err := compiler.Compile(ast)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
	} else {
		fmt.Println(module)
	}
}
