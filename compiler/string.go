package compiler

import (
	"strconv"
	"unicode/utf16"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

const PV_BITS = 5
const PV_BRANCH = 1 << PV_BITS

var LeafType = types.NewStruct(types.I32, types.I8, types.NewArray(PV_BRANCH, types.I16))
var NodeType = types.NewStruct(types.I32, types.I8Ptr, types.NewArray(PV_BRANCH, types.I8Ptr))

func (c *context) makeStringRefRoot(str string) value.Value {
	if c.global.strings[str] != nil {
		return c.global.strings[str]
	}
	stringID := "string%" + strconv.Itoa(len(c.global.strings))
	encoded := utf16.Encode([]rune(str))
	depth := 0
	if len(encoded) > 0 {
		i := (len(encoded) - 1) >> PV_BITS
		for i > 0 {
			depth++
			i = i >> PV_BITS
		}
	}

	gd, count := makePVLeaves(c, encoded, stringID+"%0")
	arrayType := types.NewArray(uint64(count), LeafType)

	if depth > 0 {
		gd, count = makePVNodes(c, count, arrayType, gd, stringID)
		arrayType = types.NewArray(uint64(count), NodeType)
	}
	ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	res := makePVHead(c, ptr, len(encoded), stringID)
	c.global.strings[str] = res
	return res
}

func makePVLeaves(c *context, elements []uint16, id string) (*ir.Global, int) {
	zero := constant.NewInt(types.I32, 0)
	n := 0
	nodes := []constant.Constant{}
	for n < len(elements) {
		cs := make([]constant.Constant, PV_BRANCH)
		i := 0
		w := 0
		for i < PV_BRANCH {
			if n < len(elements) {
				cs[i] = constant.NewInt(types.I16, int64(elements[n]))
				w++
			} else {
				cs[i] = constant.NewInt(types.I16, 0)
			}
			n++
			i++
		}
		arr := constant.NewArray(nil, cs...)
		l := constant.NewStruct(LeafType, zero, constant.NewInt(types.I8, int64(w)), arr)
		nodes = append(nodes, l)
	}

	array := constant.NewArray(nil, nodes...)
	return c.global.Module.NewGlobalDef(id, array), len(nodes)
}

func makePVNodes(c *context, count int, arrayType types.Type, src *ir.Global, id string) (*ir.Global, int) {
	depth := 0
	zero := constant.NewInt(types.I32, 0)
	for count > 1 {
		depth++
		nodes := []constant.Constant{}
		n := 0
		for n < count {
			i := 0
			cs := make([]constant.Constant, PV_BRANCH)
			for i < PV_BRANCH {
				if n < count {
					ptr := constant.NewGetElementPtr(arrayType, src, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(n)))
					cs[i] = constant.NewBitCast(ptr, types.I8Ptr)
				} else {
					cs[i] = constant.NewNull(types.I8Ptr)
				}
				n++
				i++
			}
			arr := constant.NewArray(nil, cs...)
			l := constant.NewStruct(NodeType, zero, constant.NewNull(types.I8Ptr), arr)
			nodes = append(nodes, l)
		}
		count = len(nodes)
		array := constant.NewArray(nil, nodes...)
		src = c.global.Module.NewGlobalDef(id+"%"+strconv.Itoa(depth), array)
		arrayType = types.NewArray(uint64(count), NodeType)
	}
	return src, count
}

func makePVHead(c *context, node constant.Constant, length int, id string) constant.Constant {
	headType := types.NewStruct(types.I8, types.I32, types.I32, types.I8Ptr)
	head := constant.NewStruct(
		headType,
		constant.NewInt(types.I8, 3),
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(length)),
		constant.NewBitCast(node, types.I8Ptr))

	gd := c.global.Module.NewGlobalDef(id, head)

	return constant.NewBitCast(gd, types.I8Ptr)
}
