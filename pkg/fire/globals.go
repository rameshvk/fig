package fire

import (
	"context"
	"math"
)

// Globals returns the standard globals
//
// This includes:
//
//   all standard operators
//   error(string)
//   if(condition, then, else)
//   object(key: value, ....)
//   math.Inf
//
func Globals() Value {
	code := func(c string) func(ctx context.Context) string {
		return func(ctx context.Context) string {
			return c
		}
	}

	return Object(map[Value]Value{
		String("+"):      Function(code("+"), add),
		String("-"):      Function(code("-"), sub),
		String("*"):      Function(code("*"), mul),
		String("/"):      Function(code("/"), div),
		String("<"):      Function(code("<"), less),
		String("<="):     Function(code("<="), notGreater),
		String(">"):      Function(code(">"), greater),
		String(">="):     Function(code(">="), notLess),
		String("=="):     Function(code("=="), equals),
		String("!="):     Function(code("=="), notEquals),
		String("&"):      NativeFunction(code("&"), and),
		String("|"):      NativeFunction(code("|"), or),
		String("!"):      Function(code("!"), not),
		String("."):      Function(code("."), field),
		String("{}"):     NativeFunction(code("{}"), closure),
		String("call"):   NativeFunction(code("()"), call),
		String("error"):  Function(code("error"), errorf),
		String("if"):     NativeFunction(code("if"), nativeIf),
		String("object"): Function(code("object"), objectf),
		String("math"): Object(map[Value]Value{
			String("Inf"): Number(math.Inf(+1)),
		}),
	})
}

func and(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 2 {
		return errorValue("& requires two args")
	}

	left := Eval(ctx, args[0], scope)
	if _, ok := left.Error(ctx); ok {
		return left
	}
	if b, ok := left.Bool(ctx); ok && !b {
		return boolValue(b)
	}
	return Eval(ctx, args[1], scope)
}

func or(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 2 {
		return errorValue("| requires two args")
	}

	left := Eval(ctx, args[0], scope)
	if _, ok := left.Error(ctx); ok {
		return left
	}
	if b, ok := left.Bool(ctx); !ok || b {
		if !ok {
			return left
		}
		return boolValue(b)
	}
	return Eval(ctx, args[1], scope)
}

func not(ctx context.Context, args ...Value) Value {
	if len(args) != 1 {
		return errorValue("operator requires one arg")
	}

	if _, ok := args[0].Error(ctx); ok {
		return args[0]
	}

	if b, ok := args[0].Bool(ctx); !ok || b {
		return boolValue(false)
	}

	return boolValue(true)
}

func field(ctx context.Context, args ...Value) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
	}

	return args[0].Lookup(ctx, args[1])
}

func numericArgs(ctx context.Context, args []Value) (float64, float64, Value) {
	if len(args) != 2 {
		return 0, 0, errorValue("operator requires two args")
	}
	if _, ok := args[0].Error(ctx); ok {
		return 0, 0, args[0]
	}
	if _, ok := args[1].Error(ctx); ok {
		return 0, 0, args[1]
	}
	f1, ok := args[0].Number(ctx)
	if !ok {
		return 0, 0, errorValue("not a number")
	}
	f2, ok := args[1].Number(ctx)
	if !ok {
		return 0, 0, errorValue("not a number")
	}

	return f1, f2, nil
}

func add(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(numberValue(f1+f2), err)
}

func sub(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(numberValue(f1-f2), err)
}

func mul(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(numberValue(f1*f2), err)
}

func div(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(numberValue(f1/f2), err)
}

func less(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(boolValue(f1 < f2), err)
}

func notLess(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(boolValue(f1 >= f2), err)
}

func greater(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(boolValue(f1 > f2), err)
}

func notGreater(ctx context.Context, args ...Value) Value {
	f1, f2, err := numericArgs(ctx, args)
	return checkError(boolValue(f1 <= f2), err)
}

func equals(ctx context.Context, args ...Value) Value {
	return boolValue(args[0].Equals(ctx, args[1]))
}

func notEquals(ctx context.Context, args ...Value) Value {
	return boolValue(!args[0].Equals(ctx, args[1]))
}

func errorf(ctx context.Context, args ...Value) Value {
	if len(args) == 1 {
		if s, ok := args[0].String(ctx); ok {
			return errorValue(s)
		}
	}

	return errorValue("error() takes one string only")
}

func objectf(ctx context.Context, args ...Value) Value {
	if len(args) != 1 {
		return errorValue("object() takes one arg only")
	}
	return args[0]
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

func checkError(good Value, err Value) Value {
	if err != nil {
		return err
	}
	return good
}
