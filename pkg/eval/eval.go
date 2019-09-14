// Package eval evaluates simple lisp-like S-expressions
//
// Valid input is:
//
//     Any string (value is the string itself)
//     Any numver (value is the number itself)
//     Any bool (value is the bool itself)
//     An array: [fn args...]
//         If fn is a string, it is looked up in the
//         function table and the equivalent Callable
//         is used.  The string can be of the form
//         "x:n" where n is the offset in the underlying
//         source code. In this case, x is used for the
//         lookup and n is used to report errors
//         If fn is an array, it is evaluated to find
//         the equivalent Callable.  Note that if the
//         evaluation returns a string, that is not
//         looked up.
//         The callable is then invoked with the raw
//         values of the args. This allows the callable
//         to lazy evaluate its args, if any.
//         If no callable can be found, the input is
//         considered invalid
package eval

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Encoded evaluates a json-encoded string.
func Encoded(s string, scope Scope) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, err
	}
	return Value(0, v, scope)
}

// Value evaluates an arbtirary value
func Value(offset int, v interface{}, scope Scope) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case float64:
		return v, nil
	case bool:
		return v, nil
	case []interface{}:
		return List(offset, v, scope)
	}
	return nil, fmt.Errorf("unexpected type at %d", offset)
}

// List evaluates a list of items
func List(offset int, v []interface{}, scope Scope) (interface{}, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf("empty array at %d", offset)
	}

	var fn interface{}
	var err error

	switch x := v[0].(type) {
	case string:
		offset, x = split(offset, x)
		fn, err = scope.Lookup(scope, offset, x)
	case []interface{}:
		fn, err = List(offset, x, scope)
	}
	if err != nil {
		return nil, err
	}
	if callable, ok := fn.(Callable); ok {
		return callable.Call(scope, offset, v[1:])
	}
	return nil, fmt.Errorf("not a function at %d", offset)
}

func split(offset int, x string) (int, string) {
	if parts := strings.Split(x, ":"); len(parts) == 2 {
		if off, err := strconv.Atoi(parts[1]); err == nil {
			return off, parts[0]
		}
	}
	return offset, x
}
