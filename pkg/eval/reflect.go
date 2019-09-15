package eval

import (
	"encoding/json"
	"reflect"
)

// Reflect creates a scope from any golang value.
//
// If the value is a map, it exposes a Fieldable
// interface
//
// If no match is found, the base scope is used for lookups
//
// The wrap function can be used to recursively wrap interfaces if
// desired.  Note that when defaulting to the base scope, wrap is
// not called
//
// Both base and wrap can be nil
func Reflect(v interface{}, base Scope, wrap func(interface{}) interface{}) Scope {
	if base == nil {
		base = ExtendScope(nil, nil)
	}
	if wrap == nil {
		wrap = func(x interface{}) interface{} { return x }
	}

	if s, ok := v.(Scope); ok {
		// TODO: wrap scope!
		return s
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Map:
		return &reflected_map{val, base, wrap}
	}
	return base
}

type reflected_map struct {
	reflect.Value
	base Scope
	wrap func(interface{}) interface{}
}

func (r *reflected_map) Lookup(root Scope, offset int, field string) (interface{}, error) {
	keyType := r.Type().Key()
	key := reflect.ValueOf(field)
	if keyType.Kind() != reflect.String {
		ptr := reflect.New(keyType)
		if err := json.Unmarshal([]byte(field), ptr.Interface()); err != nil {
			return r.base.Lookup(root, offset, field)
		}
		key = ptr.Elem()
	}
	if val := r.MapIndex(key); val.IsValid() {
		return r.wrap(val.Interface()), nil
	}
	return r.base.Lookup(root, offset, field)
}
