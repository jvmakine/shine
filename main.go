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

	parsed, err := grammar.Parse(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}

	module, err := compiler.Compile(parsed.ToAst())
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(2)
	}
	fmt.Println(module)
}
