package fire

import (
	"context"
)

// NativeCallable is the interface used for functions that can be
// called with expression rather than evaluated values.
type NativeCallable interface {
	Value
	NativeCall(ctx context.Context, args []interface{}, scope Value) Value
}

// NativeFunction represents a native function whose arguments are not
// pre-evaluated values but instead are expressions
//
// This is predominantly used for implementing things like `list` or `if`
// which take non-named args. Also, applicable if the function has
// lazy evaluation semantics as with `if`
func NativeFunction(code func(ctx context.Context) string, fn func(ctx context.Context, args []interface{}, scope Value) Value) NativeCallable {
	return nativeFnValue{code, fn}
}

type nativeFnValue struct {
	code func(ctx context.Context) string
	fn   func(ctx context.Context, args []interface{}, scope Value) Value
}

func (f nativeFnValue) Code(ctx context.Context) string {
	return f.code(ctx)
}

func (f nativeFnValue) HashCode() interface{} {
	return "(native)"
}

func (f nativeFnValue) NativeCall(ctx context.Context, args []interface{}, scope Value) Value {
	return f.fn(ctx, args, scope)
}

func (f nativeFnValue) Call(ctx context.Context, args ...Value) Value {
	panic("cannot call a native function directly")
}

func (f nativeFnValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add methods
	return errorValue("cannot lookup a function")
}

func (f nativeFnValue) Equals(ctx context.Context, other Value) bool {
	return f.Code(ctx) == other.Code(ctx)
}

func (f nativeFnValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (f nativeFnValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (f nativeFnValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (f nativeFnValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
