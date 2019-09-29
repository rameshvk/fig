package fire

import (
	"context"
)

// Scope creates a new scope with the provided key-value pairs
//
// The parent scope is inherited
func Scope(ctx context.Context, parent Value, pairs ...[2]Value) Value {
	s := newScope(parent)
	for _, pair := range pairs {
		s.add(ctx, pair[0], pair[1])
	}
	return s
}

func newScope(parent Value) *scope {
	err := errorValue("internal error")
	return &scope{err, parent, map[interface{}][]*nameValue{}}
}

// scope implements a scope lookup but only implements the Lookup
// part of it.  For everything else, the underlying Value is used
type scope struct {
	Value
	parent Value
	pairs  map[interface{}][]*nameValue
}

func (s *scope) add(ctx context.Context, name Value, value interface{}) Value {
	hash := name.HashCode()
	if _, ok := s.pairs[hash]; !ok {
		s.pairs[hash] = []*nameValue{}
	}
	for _, entry := range s.pairs[hash] {
		if entry.name.Equals(ctx, name) {
			return errorValue("duplicate name: " + name.Code(ctx))
		}
	}
	v, _ := value.(Value)
	entry := &nameValue{name, v, value, false}
	s.pairs[hash] = append(s.pairs[hash], entry)
	return nil
}

func (s *scope) Lookup(ctx context.Context, field Value) Value {
	for _, entry := range s.pairs[field.HashCode()] {
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
	return s.parent.Lookup(ctx, field)
}

type nameValue struct {
	name, value Value
	unevaluated interface{}
	inProgress  bool
}
