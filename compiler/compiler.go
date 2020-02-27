package compiler

import (
	"github.com/jvmakine/shine/grammar"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func evalValue(block *ir.Block, val *grammar.Value) value.Value {
	if val.Int != nil {
		return constant.NewInt(types.I32, int64(*val.Int))
	} else if val.Sub != nil {
		return evalExpression(block, val.Sub)
	}
	panic("invalid value")
}

func evalExpression(block *ir.Block, exp *grammar.Expression) value.Value {
	e := exp
	var v value.Value
	v = evalValue(block, exp.Value)

	for e.Op != nil {
		if e.Op.Add != nil {
			subVal := evalValue(block, e.Op.Add.Value)
			v = block.NewAdd(v, subVal)
			e = e.Op.Add
		} else if e.Op.Sub != nil {
			subVal := evalValue(block, e.Op.Sub.Value)
			v = block.NewSub(v, subVal)
			e = e.Op.Sub
		} else if e.Op.Mul != nil {
			subVal := evalValue(block, e.Op.Mul.Value)
			v = block.NewMul(v, subVal)
			e = e.Op.Mul
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
	v := evalExpression(entry, prg.Exp)
	ptr := entry.NewGetElementPtr(types.NewArray(3, types.I8), msg, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	entry.NewCall(printf, ptr, v)
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module
}
