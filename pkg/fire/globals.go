package fire

import (
	"context"
	"math"
)

// Globals returns the standard globals
//
// This includes:
//
//   error(string)
//   if(condition, then, else)
//   object(key: value, ....)
//   math.Inf
func Globals() Value {
	return Object(map[Value]Value{
		String("error"):  Function(errorCode, errorf),
		String("if"):     NativeFunction(ifCode, nativeIf),
		String("object"): Function(objectCode, objectf),
		String("math"): Object(map[Value]Value{
			String("Inf"): Number(math.Inf(+1)),
		}),
	})
}

func errorCode(ctx context.Context) string {
	return "error"
}

func errorf(ctx context.Context, args ...Value) Value {
	if len(args) == 1 {
		if s, ok := args[0].String(ctx); ok {
			return errorValue(s)
		}
	}

	return errorValue("error() takes one string only")
}

func objectCode(ctx context.Context) string {
	return "object"
}

func objectf(ctx context.Context, args ...Value) Value {
	if len(args) != 1 {
		return errorValue("object() takes one arg only")
	}
	return args[0]
}

func ifCode(ctx context.Context) string {
	return "if"
}

func nativeIf(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 3 {
		return errorValue("if requires 3 unnamed args")
	}
	condition := Eval(ctx, args[0], scope)
	if b, ok := condition.Bool(ctx); b && ok {
		return Eval(ctx, args[1], scope)
	} else {
		return Eval(ctx, args[2], scope)
	}
}
