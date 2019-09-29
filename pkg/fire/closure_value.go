package fire

import (
	"context"
)

// Closure creates a closure of an expression and associated scope
func Closure(expression interface{}, scope Value) Value {
	return closureValue{expression, scope}
}

type closureValue struct {
	expression interface{}
	scope      Value
}

func (c closureValue) Code(ctx context.Context) string {
	panic("not yet implemented")
}

func (c closureValue) Call(ctx context.Context, args ...Value) Value {
	if len(args) != 1 {
		return errorValue("closures always take one arg")
	}
	it := map[Value]interface{}{stringValue("it"): args[0]}
	return Eval(ctx, c.expression, newScope(it, c.scope))
}

func (c closureValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add closure methods
	return errorValue("cannot lookup a closure")
}

func (c closureValue) Equals(ctx context.Context, other Value) bool {
	// TODO: simplify this by using hash keys
	return c.Code(ctx) == other.Code(ctx)
}

func (c closureValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (c closureValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (c closureValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (c closureValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
