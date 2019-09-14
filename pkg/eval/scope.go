package eval

import (
	"fmt"
)

// Callable is any value that can be called
type Callable interface {
	Call(root Scope, offset int, args []interface{}) (interface{}, error)
}

// Fieldable is any value that allows "dot" access to fields
type Fieldable interface {
	Field(root Scope, offset int, field string) (interface{}, error)
}

// Scope is any object which a lookup of values in the scope
type Scope interface {
	Lookup(root Scope, offset int, name string) (interface{}, error)
}

// ExtendScope creates a new scope which uses the provided map
// for the initial lookup and failing that, switching to the base
// scope
func ExtendScope(v map[string]interface{}, base Scope) Scope {
	return scopeFunc(func(root Scope, offset int, name string) (interface{}, error) {
		if val, ok := v[name]; ok {
			return val, nil
		}
		if base != nil {
			return base.Lookup(root, offset, name)
		}
		return nil, fmt.Errorf("%s not found at %d", name, offset)
	})
}

// DefaultScope is the default scope which implements standard
// operators and such
var DefaultScope = ExtendScope(map[string]interface{}{
	"+":  numOperator(func(l, r float64) float64 { return l + r }),
	"-":  numOperator(func(l, r float64) float64 { return l - r }),
	"*":  numOperator(func(l, r float64) float64 { return l * r }),
	"/":  numOperator(func(l, r float64) float64 { return l / r }),
	"<":  numCmpOperator(func(l, r float64) bool { return l < r }),
	">":  numCmpOperator(func(l, r float64) bool { return l > r }),
	"<=": numCmpOperator(func(l, r float64) bool { return l <= r }),
	">=": numCmpOperator(func(l, r float64) bool { return l >= r }),
	"==": CallableFunc(equals),
	"!=": CallableFunc(notEquals),
	"&&": boolOperator(func(l, r bool) bool { return l && r }),
	"||": boolOperator(func(l, r bool) bool { return l || r }),
	"!":  CallableFunc(not),
	".":  CallableFunc(dot),
	"if": CallableFunc(ifFunc),
}, nil)

func dot(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("dot has incorrect args at %d", offset)
	}
	left, err := Value(offset, args[0], root)
	if err != nil {
		return nil, err
	}
	if fieldable, ok := left.(Fieldable); ok {
		field, err := toString(root, offset, args[1])
		if err != nil {
			return nil, err
		}
		return fieldable.Field(root, offset, field)
	}
	return nil, fmt.Errorf("dot is not valid at %d", offset)
}

func ifFunc(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("dot has incorrect args at %d", offset)
	}
	cond, err := toBool(root, offset, args[0])
	if err != nil {
		return nil, err
	}
	if cond {
		return Value(offset, args[1], root)
	}
	return Value(offset, args[2], root)
}

type numOperator func(l, r float64) float64

func (n numOperator) Call(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("op has incorrect args at %d", offset)
	}
	left, err := toNumber(root, offset, args[0])
	if err != nil {
		return nil, err
	}
	right, err := toNumber(root, offset, args[1])
	if err != nil {
		return nil, err
	}
	return n(left, right), nil
}

type numCmpOperator func(l, r float64) bool

func (n numCmpOperator) Call(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("op has incorrect args at %d", offset)
	}
	left, err := toNumber(root, offset, args[0])
	if err != nil {
		return nil, err
	}
	right, err := toNumber(root, offset, args[1])
	if err != nil {
		return nil, err
	}
	return n(left, right), nil
}

type boolOperator func(l, r bool) bool

func (b boolOperator) Call(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("op has incorrect args at %d", offset)
	}
	left, err := toBool(root, offset, args[0])
	if err != nil {
		return nil, err
	}
	right, err := toBool(root, offset, args[1])
	if err != nil {
		return nil, err
	}
	return b(left, right), nil
}

func not(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("! has incorrect args at %d", offset)
	}
	left, err := toBool(root, offset, args[0])
	if err != nil {
		return nil, err
	}
	return !left, nil
}

func equals(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("== has incorrect args at %d", offset)
	}
	left, err := Value(offset, args[0], root)
	if err != nil {
		return nil, err
	}
	right, err := Value(offset, args[1], root)
	if err != nil {
		return nil, err
	}
	return left == right, nil
}

func notEquals(root Scope, offset int, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("!= has incorrect args at %d", offset)
	}
	left, err := Value(offset, args[0], root)
	if err != nil {
		return nil, err
	}
	right, err := Value(offset, args[1], root)
	if err != nil {
		return nil, err
	}
	return left != right, nil
}

func toNumber(root Scope, offset int, arg interface{}) (float64, error) {
	v, err := Value(offset, arg, root)
	if err != nil {
		return 0, err
	}
	switch v := v.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, fmt.Errorf("not a number at %d", offset)
}

func toString(root Scope, offset int, arg interface{}) (string, error) {
	v, err := Value(offset, arg, root)
	if err != nil {
		return "", err
	}
	if s, ok := v.(string); ok {
		return s, nil
	}

	return "", fmt.Errorf("not a string at %d", offset)
}

func toBool(root Scope, offset int, arg interface{}) (bool, error) {
	v, err := Value(offset, arg, root)
	if err != nil {
		return false, err
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("not a bool at %d", offset)
}

type scopeFunc func(Scope, int, string) (interface{}, error)

func (s scopeFunc) Lookup(root Scope, offset int, name string) (interface{}, error) {
	return s(root, offset, name)
}

// CallableFunc converts a function into a callable
type CallableFunc func(Scope, int, []interface{}) (interface{}, error)

func (c CallableFunc) Call(root Scope, offset int, args []interface{}) (interface{}, error) {
	return c(root, offset, args)
}
