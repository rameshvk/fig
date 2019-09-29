package fire

import (
	"context"
)

// Object creates an object
func Object(v map[Value]Value) Value {
	return obj(v)
}

type obj map[Value]Value

func (o obj) Code(ctx context.Context) string {
	panic("not yet implemented")
}

func (o obj) HashCode() interface{} {
	return "(obj)"
}

func (o obj) Call(ctx context.Context, args ...Value) Value {
	return errorValue("cannot call an object")
}

func (o obj) Lookup(ctx context.Context, field Value) Value {
	if v, ok := o[field]; ok {
		return v
	}
	return errorValue("field not found: " + field.Code(ctx))
}

func (o obj) Equals(ctx context.Context, other Value) bool {
	panic("not yet implemented")
}

func (o obj) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (o obj) String(ctx context.Context) (string, bool) {
	return "", false
}

func (o obj) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (o obj) Error(ctx context.Context) (error, bool) {
	return nil, false
}
