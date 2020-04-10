package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/types"
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
	err = types.Infer(ast)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(2)
	}
	module, err := compiler.Compile(ast)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(3)
	}
	fmt.Println(module)
}
