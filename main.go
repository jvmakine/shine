package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func main() {
	module := ir.NewModule()
	msg := module.NewGlobalDef("msg", constant.NewCharArrayFromString("Hello world!"))
	puts := module.NewFunc("puts", types.I32, ir.NewParam("msg", types.I8Ptr))

	mainfun := module.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")
	ptr := entry.NewGetElementPtr(types.NewArray(12, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(puts, ptr)
	entry.NewRet(constant.NewInt(types.I32, 0))

	fmt.Println(module)
}
