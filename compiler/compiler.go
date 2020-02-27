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

func evalOpFactor(block *ir.Block, opf *grammar.OpFactor, left value.Value) value.Value {
	right := evalValue(block, opf.Right)
	switch *opf.Operation {
	case "*":
		return block.NewMul(left, right)
	case "/":
		return block.NewUDiv(left, right)
	default:
		panic("invalid opfactor: " + *opf.Operation)
	}
}

func evalTerm(block *ir.Block, term *grammar.Term) value.Value {
	v := evalValue(block, term.Left)
	for _, r := range term.Right {
		v = evalOpFactor(block, r, v)
	}
	return v
}

func evalOpTerm(block *ir.Block, opt *grammar.OpTerm, left value.Value) value.Value {
	right := evalTerm(block, opt.Right)
	switch *opt.Operation {
	case "+":
		return block.NewAdd(left, right)
	case "-":
		return block.NewSub(left, right)
	default:
		panic("invalid opterm: " + *opt.Operation)
	}
}

func evalExpression(block *ir.Block, prg *grammar.Expression) value.Value {
	v := evalTerm(block, prg.Left)
	for _, r := range prg.Right {
		v = evalOpTerm(block, r, v)
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
