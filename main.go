package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/jvmakine/shine/compiler"
	"github.com/jvmakine/shine/grammar"
	"github.com/jvmakine/shine/passes/callresolver"
	"github.com/jvmakine/shine/passes/closureresolver"
	"github.com/jvmakine/shine/passes/optimisation"
	"github.com/jvmakine/shine/passes/typeinference"
	"github.com/llir/llvm/ir"
)

//go:embed lib/runtime.ll
var runtime string

type Build struct {
	File string `arg:"" name:"file" help:"Source code file"`
}

type Compile struct {
	File string `arg:"" name:"file" help:"Source code file"`
}

type Run struct {
	File string `arg:"" name:"file" help:"Source code file"`
}

type CLI struct {
	Build   Build   `cmd:"" help:"Build executable"`
	Compile Compile `cmd:"" help:"Compile to LLVM IR"`
	Run     Run     `cmd:"" help:"Compile to LLVM IR"`
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli)
	err := ctx.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
}

func (cmd *Compile) Run() error {
	text, err := ioutil.ReadFile(cmd.File)
	if err != nil {
		return err
	}
	module, err := compileModule(string(text))
	if err != nil {
		return err
	}
	fmt.Println(module)
	return nil
}

func (cmd *Run) Run() error {
	_, err := exec.LookPath("lli")
	if err != nil {
		return errors.New("lli could not be found. Is LLVM installed?")
	}
	_, err = exec.LookPath("llvm-link")
	if err != nil {
		return errors.New("llvm-link could not be found. Is LLVM installed?")
	}
	text, err := ioutil.ReadFile(cmd.File)
	if err != nil {
		return err
	}
	module, err := compileModule(string(text))
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir(".", "run")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	input := filepath.Join(dir, "input.ll")
	rt := filepath.Join(dir, "runtime.ll")
	all := filepath.Join(dir, "all.ll")

	if err := ioutil.WriteFile(input, []byte(module.String()), 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(rt, []byte(runtime), 0600); err != nil {
		return err
	}
	combined, err := exec.Command("llvm-link", "-S", input, rt).Output()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(all, combined, 0600); err != nil {
		return err
	}

	c := exec.Command("lli", all, rt)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		return err
	}

	return nil
}

func compileModule(text string) (*ir.Module, error) {
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
	optimisation.ClosureRemoval(ast)
	optimisation.DeadCodeElimination(ast)

	fcat := callresolver.Collect(ast)
	module := compiler.Compile(ast, &fcat)
	return module, nil
}
