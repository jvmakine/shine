package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func expression(module *ir.Module, block *ir.Block, exp *grammar.Expression) value.Value {
	e := exp
	var v value.Value
	v = constant.NewInt(types.I32, int64(*exp.Value))

	for e.Op != nil {
		if e.Op.Add != nil {
			v = block.NewAdd(v, constant.NewInt(types.I32, int64(*e.Op.Add.Value)))
			e = e.Op.Add
		} else if e.Op.Sub != nil {
			v = block.NewSub(v, constant.NewInt(types.I32, int64(*e.Op.Sub.Value)))
			e = e.Op.Sub
		}
	}

	return v
}

func Compile(prg *grammar.Program) *ir.Module {
	module := ir.NewModule()

	msg := module.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n"))
	printf := module.NewFunc("printf", types.I32, ir.NewParam("msg", types.I8Ptr))
	printf.Sig.Variadic = true

	mainfun := module.NewFunc("main", types.I32)
	entry := mainfun.NewBlock("")
	v := expression(module, entry, prg.Exp)
	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module
}
