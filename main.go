package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/llir/llvm/ir"

	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/inferer"
	"github.com/jvmakine/shine/inferer/callresolver"
	"github.com/jvmakine/shine/optimisation"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString(0)

	module, err := Compile(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
	fmt.Println(module)
}

func Compile(text string) (*ir.Module, error) {
	parsed, err := grammar.Parse(text)
	if err != nil {
		return nil, err
	}
	ast := parsed.ToAst()
	err = inferer.Infer(ast)
	if err != nil {
		return nil, err
	}
	optimisation.Optimise(ast)
	fcat := callresolver.Resolve(ast)
	module := compiler.Compile(ast, fcat)
	return module, nil
}
