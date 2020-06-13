package compiler

import (
	"strconv"
	"unicode/utf16"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

const PV_BITS = 5
const PV_BRANCH = 1 << PV_BITS

func (c *context) makeStringRefRoot(str string) value.Value {
	if c.global.strings[str] != nil {
		return c.global.strings[str]
	}
	leafType := types.NewStruct(types.I32, types.NewArray(PV_BRANCH, types.I16))
	nodeType := types.NewStruct(types.I32, types.NewArray(PV_BRANCH, types.I8Ptr))
	headType := types.NewStruct(types.I8, types.I32, types.I32, types.I8Ptr)

	encoded := utf16.Encode([]rune(str))
	nodes := []constant.Constant{}
	i := len(encoded)
	depth := 0
	for i > 0 {
		depth++
		i = i >> PV_BITS
	}
	n := 0

	zero := constant.NewInt(types.I32, 0)
	for n < len(encoded) {
		cs := make([]constant.Constant, PV_BRANCH)
		i := 0
		for i < PV_BRANCH {
			if n < len(encoded) {
				cs[i] = constant.NewInt(types.I16, int64(encoded[n]))
			} else {
				cs[i] = constant.NewInt(types.I16, 0)
			}
			n++
			i++
		}
		arr := constant.NewArray(nil, cs...)
		l := constant.NewStruct(leafType, zero, arr)
		nodes = append(nodes, l)
	}

	depth = 0
	arrayType := types.NewArray(uint64(len(nodes)), leafType)
	array := constant.NewArray(nil, nodes...)
	id := "string%" + strconv.Itoa(depth) + "%" + strconv.Itoa(len(c.global.strings))
	gd := c.global.Module.NewGlobalDef(id, array)

	if depth <= 1 {
		depth++
		i := 0
		n := 0
		cs := make([]constant.Constant, PV_BRANCH)
		for i < PV_BRANCH {
			if n < len(nodes) {
				ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, int64(n)))
				cs[i] = constant.NewBitCast(ptr, types.I8Ptr)
			} else {
				cs[i] = constant.NewNull(types.I8Ptr)
			}
			n++
			i++
		}
		arr := constant.NewArray(nil, cs...)
		l := constant.NewStruct(nodeType, zero, arr)
		arrayType = types.NewArray(1, nodeType)
		array = constant.NewArray(nil, l)
		id = "string%" + strconv.Itoa(depth) + "%" + strconv.Itoa(len(c.global.strings))
		gd = c.global.Module.NewGlobalDef(id, array)
	} else {
		for len(nodes) > 1 {
			depth++
			new := []constant.Constant{}
			n := 0
			for n < len(nodes) {
				i := 0
				cs := make([]constant.Constant, PV_BRANCH)
				for i < PV_BRANCH {
					if n < len(nodes) {
						ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, int64(n)))
						cs[i] = constant.NewBitCast(ptr, types.I8Ptr)
					} else {
						cs[i] = constant.NewNull(types.I8Ptr)
					}
					n++
					i++
				}
				arr := constant.NewArray(nil, cs...)
				l := constant.NewStruct(nodeType, zero, arr)
				new = append(new, l)
			}
			nodes = new
			arrayType = types.NewArray(uint64(len(nodes)), nodeType)
			array = constant.NewArray(nil, nodes...)
			id = "string%" + strconv.Itoa(depth) + "%" + strconv.Itoa(len(c.global.strings))
			gd = c.global.Module.NewGlobalDef(id, array)
		}
	}
	ptr := constant.NewGetElementPtr(arrayType, gd, constant.NewInt(types.I32, 0))

	head := constant.NewStruct(
		headType,
		constant.NewInt(types.I8, 3),
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(len(encoded))),
		constant.NewBitCast(ptr, types.I8Ptr))

	id = "string%" + strconv.Itoa(len(c.global.strings))
	gd = c.global.Module.NewGlobalDef(id, head)

	res := constant.NewBitCast(gd, types.I8Ptr)
	c.global.strings[str] = res
	return res
}
