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
	var arg Value
	if len(args) == 1 && assignPattern.Match(args[0]) != nil {
		arg = Eval(ctx, args[0], newScope(scope, outer))
	} else {
		var errv Value
		arg, errv = evalArgument(ctx, args, newScope(scope, outer))
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
	for k := range inner {
		if _, ok := o[k]; ok {
			return nil, errorValue("duplicate name")
		}
	}
	fullScope := newScope(inner, scope)
	for k := range inner {
		o[k] = fullScope.Lookup(ctx, k)
	}
	return o, nil
}

func newScope(m map[Value]interface{}, parent Value) Value {
	s := &scope{Value: errorValue("internal error"), parent: parent}
	for name, unevaluated := range m {
		value := Value(nil)
		if val, ok := unevaluated.(Value); ok {
			value = val
		}
		s.nameValues = append(s.nameValues, &nameValue{
			name:        name,
			value:       value,
			unevaluated: unevaluated,
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
