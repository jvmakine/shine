// Package ast contains the definition of the inital program structure
// as parsed from the translation unit
package ast

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jvmakine/shine/types"
)

type Ast interface {
	Visit(
		before VisitFunc,
		after VisitFunc,
		crawl bool,
		rewrite RewriteFunc,
		ctx *VisitContext,
	) error
}

type Expression interface {
	Ast

	Type() types.Type

	CopyWithCtx(ctx *types.TypeCopyCtx) Expression

	Format(builder *strings.Builder, level int, options *FormatOptions)
}

// Expressions

type Op struct {
	Name string

	OpType types.Type
}

func (e *Op) Type() types.Type {
	return e.OpType
}

func (e *Op) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return &Op{
		Name:   e.Name,
		OpType: e.OpType.Copy(ctx),
	}
}

func (e *Op) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *Op) Format(builder *strings.Builder, level int, options *FormatOptions) {
	builder.WriteString(e.Name)
}

func NewOp(name string) *Op {
	return &Op{Name: name}
}

type Id struct {
	Name string

	IdType types.Type
}

func (e *Id) Type() types.Type {
	return e.IdType
}

func (e *Id) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return &Id{
		Name:   e.Name,
		IdType: e.IdType.Copy(ctx),
	}
}

func (e *Id) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	if crawl {
		if r, c := ctx.Resolve(e.Name); r != nil {
			if !ctx.global.visited[r] {
				ctx.global.visited[r] = true
				if err := r.Visit(before, after, crawl, rewrite, c.WithAssignment(e.Name)); err != nil {
					return err
				}
			}
		}
	}
	return after(e, ctx)
}

func (e *Id) Format(b *strings.Builder, level int, options *FormatOptions) {
	b.WriteString("\"")
	b.WriteString(e.Name)
	b.WriteString("\"")
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func NewId(name string) *Id {
	return &Id{Name: name}
}

type Const struct {
	Int    *int64
	Real   *float64
	Bool   *bool
	String *string

	ConstType types.Type
}

func (e *Const) Type() types.Type {
	return e.ConstType
}

func (e *Const) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return e
}

func (e *Const) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *Const) Format(b *strings.Builder, level int, options *FormatOptions) {
	if e.Int != nil {
		b.WriteString(strconv.FormatInt(*e.Int, 10))
	} else if e.Real != nil {
		b.WriteString(strconv.FormatFloat(*e.Real, 'f', -1, 64))
	} else if e.Bool != nil {
		if *e.Bool {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	} else if e.String != nil {
		b.WriteString("\"")
		b.WriteString(*e.String)
		b.WriteString("\"")
	} else {
		panic("invalid const")
	}
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func NewConst(v interface{}) *Const {
	if i, ok := v.(int); ok {
		i64 := int64(i)
		return &Const{Int: &i64}
	} else if f, ok := v.(float64); ok {
		return &Const{Real: &f}
	} else if s, ok := v.(string); ok {
		return &Const{String: &s}
	} else if b, ok := v.(bool); ok {
		return &Const{Bool: &b}
	}
	panic("illegal const")
}

type TypeDecl struct {
	Exp      Expression
	DeclType types.Type
}

func (e *TypeDecl) Type() types.Type {
	return e.DeclType
}

func (e *TypeDecl) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return &TypeDecl{
		Exp:      e.Exp.CopyWithCtx(ctx),
		DeclType: e.DeclType.Copy(ctx),
	}
}

func (e *TypeDecl) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	e.Exp = rewrite(e.Exp, ctx).(Expression)
	err = e.Exp.Visit(before, after, crawl, rewrite, ctx)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *TypeDecl) Format(b *strings.Builder, level int, options *FormatOptions) {
	e.Exp.Format(b, level, options)
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func NewTypeDecl(typ types.Type, exp Expression) *TypeDecl {
	return &TypeDecl{
		DeclType: typ,
		Exp:      exp,
	}
}

type FieldAccessor struct {
	Exp    Expression
	Field  string
	FAType types.Type
}

func (e *FieldAccessor) Type() types.Type {
	return e.FAType
}

func (e *FieldAccessor) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return &FieldAccessor{
		Exp:    e.Exp.CopyWithCtx(ctx),
		Field:  e.Field,
		FAType: e.FAType.Copy(ctx),
	}
}

func (e *FieldAccessor) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	e.Exp = rewrite(e.Exp, ctx).(Expression)
	err = e.Exp.Visit(before, after, crawl, rewrite, ctx)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *FieldAccessor) Format(b *strings.Builder, level int, options *FormatOptions) {
	e.Exp.Format(b, level, options)
	b.WriteString(".")
	b.WriteString(e.Field)
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func NewFieldAccessor(name string, exp Expression) *FieldAccessor {
	return &FieldAccessor{
		Exp:   exp,
		Field: name,
	}
}

// Functions

type FCall struct {
	Function Expression
	Params   []Expression

	CallType types.Type
}

func (e *FCall) Type() types.Type {
	return e.CallType
}

func (e *FCall) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	ps := make([]Expression, len(e.Params))
	for i, p := range e.Params {
		ps[i] = p.CopyWithCtx(ctx)
	}
	return &FCall{
		Function: e.Function.CopyWithCtx(ctx),
		Params:   ps,
		CallType: e.CallType.Copy(ctx),
	}
}

func (e *FCall) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	for _, p := range e.Params {
		err := p.Visit(before, after, crawl, rewrite, ctx)
		if err != nil {
			return err
		}
	}
	e.Function = rewrite(e.Function, ctx).(Expression)
	for i, p := range e.Params {
		e.Params[i] = rewrite(p, ctx).(Expression)
	}
	err = e.Function.Visit(before, after, crawl, rewrite, ctx)

	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *FCall) Format(b *strings.Builder, level int, options *FormatOptions) {
	e.Function.Format(b, level, options)
	b.WriteString("(")
	for i, p := range e.Params {
		p.Format(b, level, options)
		if i < len(e.Params)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString(")")
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func (a *FCall) RootFunc() Expression {
	if f, ok := a.Function.(*FCall); ok {
		return f.RootFunc()
	}
	return a.Function
}

func (call *FCall) MakeFunType() types.Type {
	funps := make([]types.Type, len(call.Params)+1)
	for i, p := range call.Params {
		funps[i] = p.Type()
	}
	funps[len(call.Params)] = call.Type()
	return types.MakeFunction(funps...)
}

func NewFCall(fun Expression, params ...Expression) *FCall {
	return &FCall{
		Function: fun,
		Params:   params,
	}
}

type FParam struct {
	Name      string
	ParamType types.Type
}

type FDef struct {
	Params []*FParam
	Body   Expression

	Closure *types.Structure
}

func (e *FDef) Type() types.Type {
	ts := make([]types.Type, len(e.Params)+1)
	for i, p := range e.Params {
		ts[i] = p.ParamType
	}
	ts[len(e.Params)] = e.Body.Type()
	return types.MakeFunction(ts...)
}

func (a *FDef) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	pc := make([]*FParam, len(a.Params))
	for i, p := range a.Params {
		pc[i] = &FParam{
			ParamType: p.ParamType.Copy(ctx),
			Name:      p.Name,
		}
	}
	return &FDef{
		Params:  pc,
		Body:    a.Body.CopyWithCtx(ctx),
		Closure: a.Closure.Copy(ctx),
	}
}

func (e *FDef) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	e.Body = rewrite(e.Body, ctx).(Expression)
	err = e.Body.Visit(before, after, crawl, rewrite, ctx.WithDef(e))
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *FDef) Format(b *strings.Builder, level int, options *FormatOptions) {
	b.WriteString("(")
	for i, p := range e.Params {
		b.WriteString(p.Name)
		if options.Types {
			b.WriteString(":")
			b.WriteString(p.ParamType.Signature())
		}
		if i < len(e.Params)-1 || e.HasClosure() {
			b.WriteString(",")
		}
	}
	if e.HasClosure() {
		b.WriteString("[")
		for i, p := range e.Closure.Fields {
			b.WriteString(p.Name)
			if options.Types {
				b.WriteString(":")
				b.WriteString(p.Type.Signature())
			}
			if i < len(e.Closure.Fields)-1 {
				b.WriteString(",")
			}
		}
		b.WriteString("]")
	}
	b.WriteString(") => ")
	e.Body.Format(b, level, options)
}

func (a *FDef) ParamOf(name string) *FParam {
	for _, p := range a.Params {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (a *FDef) HasClosure() bool {
	return a.Closure != nil && len(a.Closure.Fields) > 0
}

func NewFDef(body Expression, params ...interface{}) *FDef {
	pars := make([]*FParam, len(params))
	for i, p := range params {
		if s, ok := p.(string); ok {
			pars[i] = &FParam{Name: s}
		} else if f, ok := p.(*FParam); ok {
			pars[i] = f
		} else {
			panic("invalid param")
		}
	}
	return &FDef{
		Body:   body,
		Params: pars,
	}
}

// Blocks

type Block struct {
	Def   *Definitions
	Value Expression
	ID    int
}

func (e *Block) Type() types.Type {
	return e.Value.Type()
}

func (e *Block) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	return &Block{
		Def:   e.Def.copy(ctx),
		Value: e.Value.CopyWithCtx(ctx),
		ID:    e.ID,
	}
}

func (e *Block) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	sub := ctx.WithBlock(e)
	if !crawl {
		err := e.Def.Visit(before, after, crawl, rewrite, sub)
		if err != nil {
			return err
		}
	}
	e.Value = rewrite(e.Value, ctx).(Expression)
	err = e.Value.Visit(before, after, crawl, rewrite, sub)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (e *Block) Format(b *strings.Builder, level int, options *FormatOptions) {
	b.WriteString("{")
	newline(b, level+2)
	for k, p := range e.Def.Assignments {
		b.WriteString(k)
		b.WriteString(" = ")
		p.Format(b, level+2, options)
		newline(b, level+2)
	}
	e.Value.Format(b, level+2, options)
	newline(b, level)
	b.WriteString("}")
	if options.Types {
		b.WriteString(":")
		b.WriteString(e.Type().Signature())
	}
}

func (b *Block) CheckValueCycles() error {
	names := map[string]Expression{}

	type ToDo struct {
		id   string
		path []string
	}
	todo := []ToDo{}

	for k, a := range b.Def.Assignments {
		names[k] = a
		todo = append(todo, ToDo{id: k, path: []string{}})
	}

	for len(todo) > 0 {
		i := todo[0]
		todo = todo[1:]
		for _, p := range i.path {
			if p == i.id {
				return errors.New("recursive value: " + cycleToStr(i.path, i.id))
			}
		}
		exp := b.Def.Assignments[i.id]
		if _, ok := exp.(*FDef); !ok {
			ids := CollectIds(exp)
			for _, id := range ids {
				if names[id] != nil {
					todo = append(todo, ToDo{id: id, path: append(i.path, i.id)})
				}
			}
		}
	}
	return nil
}

func NewDefinitions() *Definitions {
	return &Definitions{
		Assignments: map[string]Expression{},
		Interfaces:  map[types.Type][]*Interface{},
	}
}

func (e *Definitions) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {

	for n, a := range e.Assignments {
		e.Assignments[n] = rewrite(a, ctx).(Expression)
		err := a.Visit(before, after, crawl, rewrite, ctx.WithAssignment(n))
		if err != nil {
			return err
		}
	}
	for _, i := range e.Interfaces {
		for _, in := range i {
			err := in.Definitions.Visit(before, after, crawl, rewrite, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewBlock(body Expression) *Block {
	return &Block{
		ID:    0,
		Value: body,
		Def:   NewDefinitions(),
	}
}

func (b *Block) WithAssignment(name string, value interface{}) *Block {
	dc := b.Def.WithAssignment(name, value)
	return &Block{
		ID:    b.ID,
		Value: b.Value,
		Def:   dc,
	}
}

func (b *Block) WithInterface(typ types.Type, defs *Definitions) *Block {
	dc := b.Def.WithInterface(typ, defs)
	return &Block{
		ID:    b.ID,
		Value: b.Value,
		Def:   dc,
	}
}

func (d *Definitions) WithAssignment(name string, value interface{}) *Definitions {
	dc := d.shallowCopy()
	if ex, ok := value.(Expression); ok {
		dc.Assignments[name] = ex
	} else {
		panic("invalid assignment")
	}
	return dc
}

func (d *Definitions) WithInterface(typ types.Type, defs *Definitions) *Definitions {
	dc := d.shallowCopy()
	dc.Interfaces[typ] = append(dc.Interfaces[typ], &Interface{
		Definitions: defs,
	})
	return dc
}

type Definitions struct {
	Assignments map[string]Expression
	Interfaces  map[types.Type][]*Interface
}

func (a *Definitions) shallowCopy() *Definitions {
	ac := map[string]Expression{}
	for k, v := range a.Assignments {
		ac[k] = v
	}
	ic := map[types.Type][]*Interface{}
	for k, v := range a.Interfaces {
		ic[k] = v
	}
	return &Definitions{
		Assignments: ac,
		Interfaces:  ic,
	}
}

func (a *Definitions) copy(ctx *types.TypeCopyCtx) *Definitions {
	ac := map[string]Expression{}
	for k, v := range a.Assignments {
		ac[k] = v.CopyWithCtx(ctx)
	}
	ic := map[types.Type][]*Interface{}
	for k, lst := range a.Interfaces {
		is := make([]*Interface, len(lst))
		for i, in := range lst {
			is[i] = in.CopyWithCtx(ctx)
		}
		ic[k] = is
	}
	return &Definitions{
		Assignments: ac,
		Interfaces:  ic,
	}
}

type Interface struct {
	Definitions *Definitions
}

func (i *Interface) CopyWithCtx(ctx *types.TypeCopyCtx) *Interface {
	return &Interface{
		Definitions: i.Definitions.copy(ctx),
	}
}

// Types

type StructField struct {
	Name string
	Type types.Type
}

type Struct struct {
	Fields     []*StructField
	StructType types.Type
}

func NewStruct(fields ...StructField) *Struct {
	fs := make([]*StructField, len(fields))
	for i, f := range fields {
		ff := f
		fs[i] = &ff
	}
	return &Struct{
		Fields: fs,
	}
}

func (e *Struct) Visit(before VisitFunc, after VisitFunc, crawl bool, rewrite RewriteFunc, ctx *VisitContext) error {
	err := before(e, ctx)
	if err != nil {
		return err
	}
	return after(e, ctx)
}

func (s *Struct) CopyWithCtx(ctx *types.TypeCopyCtx) Expression {
	fs := make([]*StructField, len(s.Fields))
	for i, f := range s.Fields {
		fs[i] = &StructField{
			Name: f.Name,
			Type: f.Type.Copy(ctx),
		}
	}
	return &Struct{
		Fields:     fs,
		StructType: s.StructType.Copy(ctx),
	}
}

func (s *Struct) Type() types.Type {
	ts := make([]types.Type, len(s.Fields)+1)
	for i, t := range s.Fields {
		ts[i] = t.Type
	}
	ts[len(s.Fields)] = s.StructType
	return types.MakeFunction(ts...)
}

func (s *Struct) Format(builder *strings.Builder, level int, options *FormatOptions) {
	builder.WriteString("<structure>")
}

func CollectIds(exp Ast) []string {
	ids := map[string]bool{}
	exp.Visit(func(v Ast, c *VisitContext) error {
		if id, ok := v.(*Id); ok {
			name := id.Name
			ids[name] = true
		}
		return nil
	}, NullFun, false, IdRewrite, NewVisitCtx())
	result := []string{}
	for k := range ids {
		result = append(result, k)
	}
	return result
}

func cycleToStr(arr []string, v string) string {
	res := ""
	for _, a := range arr {
		res = res + a + " -> "
	}
	return res + v
}
