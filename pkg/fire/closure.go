package fire

import (
	"context"
)

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
