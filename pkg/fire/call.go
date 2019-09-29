package fire

import (
	"context"

	"github.com/rameshvk/fig/pkg/match"
)

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
	if native, ok := fn.(NativeCallable); ok {
		return native.NativeCall(ctx, args, s)
	}
	var arg Value
	if len(args) == 1 && assignPattern.Match(args[0]) != nil {
		arg = Eval(ctx, args[0], s)
	} else {
		var errv Value
		arg, errv = evalArgument(ctx, args, s)
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
