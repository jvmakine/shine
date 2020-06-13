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

var LeafType = types.NewStruct(types.I32, types.NewArray(PV_BRANCH, types.I16))
var NodeType = types.NewStruct(types.I32, types.NewArray(PV_BRANCH, types.I8Ptr))

func (c *context) makeStringRefRoot(str string) value.Value {
	if c.global.strings[str] != nil {
		return c.global.strings[str]
	}
	stringID := "string%" + strconv.Itoa(len(c.global.strings))
	encoded := utf16.Encode([]rune(str))
	i := len(encoded)
	depth := 0
	for i > 0 {
		depth++
		i = i >> PV_BITS
	}

	zero := constant.NewInt(types.I32, 0)
	gd, count := makePVLeaves(c, encoded, stringID+"%"+strconv.Itoa(depth))
	arrayType := types.NewArray(uint64(count), LeafType)

	if depth <= 1 {
		depth++
		gd = makePVNode(c, arrayType, 0, gd, count, stringID+"%"+strconv.Itoa(depth))
		arrayType = types.NewArray(1, NodeType)
	} else {
		for count > 1 {
			depth++
			new := []constant.Constant{}
			n := 0
			for n < count {
				i := 0
				cs := make([]constant.Constant, PV_BRANCH)
				for i < PV_BRANCH {
					if n < count {
						ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, int64(n)))
						cs[i] = constant.NewBitCast(ptr, types.I8Ptr)
					} else {
						cs[i] = constant.NewNull(types.I8Ptr)
					}
					n++
					i++
				}
				arr := constant.NewArray(nil, cs...)
				l := constant.NewStruct(NodeType, zero, arr)
				new = append(new, l)
			}
			count = len(new)
			arrayType = types.NewArray(uint64(count), NodeType)
			array := constant.NewArray(nil, new...)
			id := stringID + "%" + strconv.Itoa(depth)
			gd = c.global.Module.NewGlobalDef(id, array)
		}
	}
	ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, 0))
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
		for i < PV_BRANCH {
			if n < len(elements) {
				cs[i] = constant.NewInt(types.I16, int64(elements[n]))
			} else {
				cs[i] = constant.NewInt(types.I16, 0)
			}
			n++
			i++
		}
		arr := constant.NewArray(nil, cs...)
		l := constant.NewStruct(LeafType, zero, arr)
		nodes = append(nodes, l)
	}

	array := constant.NewArray(nil, nodes...)
	return c.global.Module.NewGlobalDef(id, array), len(nodes)
}

func makePVNode(c *context, arrayType types.Type, offset int, gd constant.Constant, nodes int, id string) *ir.Global {
	zero := constant.NewInt(types.I32, 0)
	cs := make([]constant.Constant, PV_BRANCH)
	i := 0
	for i < PV_BRANCH {
		if i < nodes {
			ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, int64(i+offset)))
			cs[i] = constant.NewBitCast(ptr, types.I8Ptr)
		} else {
			cs[i] = constant.NewNull(types.I8Ptr)
		}
		i++
	}
	arr := constant.NewArray(nil, cs...)
	l := constant.NewStruct(NodeType, zero, arr)
	array := constant.NewArray(nil, l)
	return c.global.Module.NewGlobalDef(id, array)
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
