package fire

import (
	"context"
)

// Scope creates a new scope with the provided key-value pairs
//
// The parent scope is inherited
func Scope(parent Value, pairs ...[2]Value) Value {
	// TODO: get rid of the map here in favor of a better structure
	mm := map[Value]interface{}{}
	for _, pair := range pairs {
		mm[pair[0]] = pair[1]
	}
	return newScope(mm, parent)
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
