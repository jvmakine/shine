package compiler

import (
	"strconv"

	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func add(exp *grammar.Expression) int {
	if exp.Add != nil {
		return *exp.Value + add(exp.Add)
	}
	return *exp.Value
}

func Compile(prg *grammar.Program) *ir.Module {
	module := ir.NewModule()
	str := strconv.Itoa(add(prg.Exp))
	msg := module.NewGlobalDef("msg", constant.NewCharArrayFromString(str))
	puts := module.NewFunc("puts", types.I32, ir.NewParam("msg", types.I8Ptr))

	mainfun := module.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")
	ptr := entry.NewGetElementPtr(types.NewArray(uint64(len(str)), types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(puts, ptr)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module
}
