// Package fire is the fig runtime
package fire

import (
	"context"
	"strings"
)

type Value interface {
	// generates the fig language representation of the value.
	Code(ctx context.Context) string

	// Call with the provided args.
	Call(ctx context.Context, args ...Value) Value

	// Lookup looks up the field/key
	Lookup(ctx context.Context, field Value) Value

	// Equals compares one value against another
	Equals(ctx context.Context, other Value) bool

	// Number converts the value to a float64 if possible
	Number(ctx context.Context) (float64, bool)

	// String converts the value to a string if possible
	String(ctx context.Context) (string, bool)

	// Bool converts the value to a bool if possible
	Bool(ctx context.Context) (bool, bool)

	// ToError converts the value to an error if possible
	Error(ctx context.Context) (error, bool)
}

// Eval evaluates a fig list expression
func Eval(ctx context.Context, v interface{}, scope Value) Value {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return errorValue("invalid list expression")
	}

	s := list[0].(string)
	parts := strings.Split(s, ":")

	// TODO: save the location into ctx for error reporting
	return builtin(ctx, parts[0], list[1:], scope)
}

func builtin(ctx context.Context, s string, args []interface{}, scope Value) Value {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}
	switch s {
	case "bool":
		if b, ok := arg.(bool); ok {
			return boolValue(b)
		}
	case "string":
		if s, ok := arg.(string); ok {
			return stringValue(s)
		}
	case "number":
		if f, ok := arg.(float64); ok {
			return numberValue(f)
		}
	case "name":
		if s, ok := arg.(string); ok {
			return scope.Lookup(ctx, stringValue(s))
		}
	case "+":
		return numericOperator(ctx, args, scope, add)
	case "-":
		if len(args) == 1 {
			args = []interface{}{zero, args[0]}
		}
		return numericOperator(ctx, args, scope, sub)
	case "*":
		return numericOperator(ctx, args, scope, mul)
	case "/":
		return numericOperator(ctx, args, scope, div)
	case "<":
		return numericComparisonOperator(ctx, args, scope, less)
	case "<=":
		return numericComparisonOperator(ctx, args, scope, notgreater)
	case ">":
		return numericComparisonOperator(ctx, args, scope, greater)
	case ">=":
		return numericComparisonOperator(ctx, args, scope, notless)
	case "&":
		return and(ctx, args, scope)
	case "|":
		return or(ctx, args, scope)
	case "!":
		return not(ctx, args, scope)
	case ".":
		return field(ctx, args, scope)
	case "==":
		return comparisonOperator(ctx, args, scope, equals)
	case "!=":
		return comparisonOperator(ctx, args, scope, notEquals)
	case "{}":
		return closure(ctx, args, scope)
	case "call":
		return call(ctx, args, scope)
	}
	return errorValue("invalid list expression")
}

func numericOperator(ctx context.Context, args []interface{}, scope Value, f func(f1, f2 float64) float64) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
	}
	left := Eval(ctx, args[0], scope)
	if _, ok := left.Error(ctx); ok {
		return left
	}
	f1, ok := left.Number(ctx)
	if !ok {
		return errorValue("not a number")
	}
	right := Eval(ctx, args[1], scope)
	if _, ok := right.Error(ctx); ok {
		return right
	}
	f2, ok := right.Number(ctx)
	if !ok {
		return errorValue("not a number")
	}

	return numberValue(f(f1, f2))
}

func numericComparisonOperator(ctx context.Context, args []interface{}, scope Value, f func(f1, f2 float64) bool) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
	}
	left := Eval(ctx, args[0], scope)
	if _, ok := left.Error(ctx); ok {
		return left
	}
	f1, ok := left.Number(ctx)
	if !ok {
		return errorValue("not a number")
	}
	right := Eval(ctx, args[1], scope)
	if _, ok := right.Error(ctx); ok {
		return right
	}
	f2, ok := right.Number(ctx)
	if !ok {
		return errorValue("not a number")
	}

	return boolValue(f(f1, f2))
}

func comparisonOperator(ctx context.Context, args []interface{}, scope Value, f func(ctx context.Context, v1, v2 Value) Value) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
	}
	left := Eval(ctx, args[0], scope)
	if _, ok := left.Error(ctx); ok {
		return left
	}
	right := Eval(ctx, args[1], scope)
	if _, ok := right.Error(ctx); ok {
		return right
	}
	return f(ctx, left, right)
}

func and(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
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
		return errorValue("operator requires two args")
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

func not(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 1 {
		return errorValue("operator requires one arg")
	}

	arg := Eval(ctx, args[0], scope)
	if _, ok := arg.Error(ctx); ok {
		return arg
	}
	if b, ok := arg.Bool(ctx); !ok || b {
		return boolValue(false)
	}
	return boolValue(true)
}

func field(ctx context.Context, args []interface{}, scope Value) Value {
	if len(args) != 2 {
		return errorValue("operator requires two args")
	}

	return Eval(ctx, args[0], scope).Lookup(ctx, Eval(ctx, args[1], scope))
}

func add(f1, f2 float64) float64 {
	return f1 + f2
}

func sub(f1, f2 float64) float64 {
	return f1 - f2
}

func mul(f1, f2 float64) float64 {
	return f1 * f2
}

func div(f1, f2 float64) float64 {
	return f1 / f2
}

func less(f1, f2 float64) bool {
	return f1 < f2
}

func notless(f1, f2 float64) bool {
	return f1 >= f2
}

func greater(f1, f2 float64) bool {
	return f1 > f2
}

func notgreater(f1, f2 float64) bool {
	return f1 <= f2
}

func equals(ctx context.Context, l, r Value) Value {
	return boolValue(l.Equals(ctx, r))
}

func notEquals(ctx context.Context, l, r Value) Value {
	return boolValue(!l.Equals(ctx, r))
}

var zero = []interface{}{"number:0:1", float64(0)}
