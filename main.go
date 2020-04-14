package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/typeinferer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString(0)

	parsed, err := grammar.Parse(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
	ast := parsed.ToAst()
	err = typeinferer.Infer(ast)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(2)
	}
	module := compiler.Compile(ast)
	fmt.Println(module)
}
