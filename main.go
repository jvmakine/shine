package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func main() {
	message := constant.NewCharArrayFromString("Hello world!")

	m := ir.NewModule()
	msg := m.NewGlobalDef("msg", message)

	puts := m.NewFunc("puts", types.I32, ir.NewParam("msg", types.I8Ptr))

	mainfun := m.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")

	ptr := entry.NewGetElementPtr(types.NewArray(12, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	ptr.Typ = types.I8Ptr
	entry.NewCall(puts, ptr)
	entry.NewRet(constant.NewInt(types.I32, 0))

	fmt.Println(m)
}
