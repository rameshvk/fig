package fire

import (
	"context"
)

// Function creates a function value
func Function(code func(ctx context.Context) string, fn func(ctx context.Context, args ...Value) Value) Value {
	return funcValue{code, fn}
}

type funcValue struct {
	code func(ctx context.Context) string
	run  func(ctx context.Context, args ...Value) Value
}

func (f funcValue) Code(ctx context.Context) string {
	return f.code(ctx)
}

func (f funcValue) Call(ctx context.Context, args ...Value) Value {
	return f.run(ctx, args...)
}

func (f funcValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add methods
	return errorValue("cannot lookup a function")
}

func (f funcValue) Equals(ctx context.Context, other Value) bool {
	return f.Code(ctx) == other.Code(ctx)
}

func (f funcValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (f funcValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (f funcValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (f funcValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
