package types

type signatureContext struct {
	variableCount int
	variables     map[VariableID]string
	definingNamed map[string]bool
}

func Signature(t Type) string {
	varm := map[VariableID]string{}
	ds := map[string]bool{}
	ctx := signatureContext{variableCount: 0, variables: varm, definingNamed: ds}
	return t.signature(&ctx)
}
