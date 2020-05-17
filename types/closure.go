package types

type ClosureParam struct {
	Name string
	Type Type
}

type Closure []ClosureParam

func (c *Closure) Copy(ctx *TypeCopyCtx) *Closure {
	if c == nil {
		return nil
	}
	params := make(Closure, len(*c))
	for i, p := range *c {
		params[i] = ClosureParam{
			Name: p.Name,
			Type: p.Type.Copy(ctx),
		}
	}
	return &params
}
