package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/llir/llvm/ir"

	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/jvmakine/shine/passes/closureresolver"
	"github.com/jvmakine/shine/passes/optimisation"
	"github.com/jvmakine/shine/passes/typeinference"
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
	err = typeinference.Infer(ast)
	if err != nil {
		return nil, err
	}
	optimisation.SequentialFunctionPass(ast)
	callresolver.ResolveFunctions(ast)
	closureresolver.CollectClosures(ast)
	optimisation.DeadCodeElimination(ast)
	fcat := callresolver.Collect(ast)
	module := compiler.Compile(ast, &fcat)
	return module, nil
}
