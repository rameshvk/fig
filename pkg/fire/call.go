package fire

import (
	"context"
)

func call(ctx context.Context, args []interface{}, outer Value) Value {
	fn := Eval(ctx, args[0], outer)
	args, scope, err := filterArgs(args[1:])
	if err != nil {
		return err
	}
	if len(args) != 1 {
		return errorValue("multiple args nyi")
	}
	arg := Eval(ctx, args[0], newScope(scope, outer))
	return fn.Call(ctx, arg)
}

func newScope(m map[Value]interface{}, parent Value) Value {
	s := &scope{Value: errorValue("internal error"), parent: parent}
	for k, v := range m {
		s.nameValues = append(s.nameValues, &nameValue{
			name:        k,
			unevaluated: v,
		})
	}
	return s
}

// scope implements a scope lookup but only implements the Lookup
// part of it.  For everything else, the underlying Value is used
type scope struct {
	Value
	parent     Value
	nameValues []*nameValue
}

type nameValue struct {
	name, value Value
	unevaluated interface{}
	inProgress  bool
}

func (s *scope) Lookup(ctx context.Context, field Value) Value {
	for _, entry := range s.nameValues {
		// TODO: make name calculatable as well
		if !entry.name.Equals(ctx, field) {
			continue
		}
		if entry.value == nil && !entry.inProgress {
			entry.inProgress = true
			entry.value = Eval(ctx, entry.unevaluated, s)
			entry.inProgress = false
		}

		if entry.inProgress {
			return errorValue("recursion detected: " + field.Code(ctx))
		}
		return entry.value
	}
	if s.parent == nil {
		return errorValue("name not found: " + field.Code(ctx))
	}
	return s.parent.Lookup(ctx, field)
}
