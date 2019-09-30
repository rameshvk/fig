package fire

import (
	"context"
	"strings"

	"github.com/rameshvk/fig/pkg/match"
)

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
	case "-":
		if len(args) == 1 {
			args = []interface{}{zero, args[0]}
		}
		fallthrough // return numericOperator(ctx, args, scope, sub)
	default:
		fn := scope.Lookup(ctx, stringValue(s))
		if native, ok := fn.(NativeCallable); ok {
			return native.NativeCall(ctx, args, scope)
		}
		values := make([]Value, len(args))
		for kk, arg := range args {
			values[kk] = Eval(ctx, arg, scope)
		}
		return fn.Call(ctx, values...)
	}

	return errorValue("invalid expression")
}

func call(ctx context.Context, args []interface{}, outer Value) Value {
	fn := Eval(ctx, args[0], outer)
	args, scope, err := filterArgs(args[1:])
	if err != nil {
		return err
	}
	s := newScope(outer)
	for k, v := range scope {
		s.add(ctx, k, v)
	}
	return callValue(ctx, fn, args, s)
}

func callValue(ctx context.Context, fn Value, args []interface{}, scope Value) Value {
	if native, ok := fn.(NativeCallable); ok {
		return native.NativeCall(ctx, args, scope)
	}
	var arg Value
	if len(args) == 1 && assignPattern.Match(args[0]) != nil {
		arg = Eval(ctx, args[0], scope)
	} else {
		var errv Value
		arg, errv = evalArgument(ctx, args, scope)
		if errv != nil {
			return errv
		}
	}
	return fn.Call(ctx, arg)
}

func evalArgument(ctx context.Context, args []interface{}, scope Value) (Value, Value) {
	o := obj{}
	inner := map[Value]interface{}{}
	for _, arg := range args {
		var name string
		p := match.Pattern([]interface{}{match.StringPrefix("name"), &name})
		if err := p.Match(arg); err == nil {
			key := stringValue(name)
			if _, ok := o[key]; ok {
				return nil, errorValue("duplicate name")
			}
			o[key] = scope.Lookup(ctx, key)
		} else if err := accumulateAssign(inner, arg); err != nil {
			return nil, err
		}
	}
	s := newScope(scope)
	for k, v := range inner {
		if err := s.add(ctx, k, v); err != nil {
			return nil, err
		}
	}
	for k := range inner {
		o[k] = s.Lookup(ctx, k)
	}
	return o, nil
}

func closure(ctx context.Context, args []interface{}, outer Value) Value {
	args, scope, err := filterArgs(args)
	if err != nil {
		return err
	}
	var result interface{}
	for _, arg := range args {
		if result == nil && assignPattern.Match(arg) != nil {
			result = arg
		} else if err := accumulateAssign(scope, arg); err != nil {
			return err
		}
	}
	if result == nil {
		return errorValue("no expression provided")
	}
	if _, ok := scope[stringValue("it")]; ok {
		return errorValue("cannot define value for it")
	}

	s := newScope(outer)
	for k, v := range scope {
		s.add(ctx, k, v)
	}
	return closureValue{result, s}
}

var zero = []interface{}{"number:0:1", float64(0)}
