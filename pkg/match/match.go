// Package match is a tiny pattern matching and capture library
//
// This can be seen as equivalent to Erlang's pattern matching
// but in semi-idiomatic Go:
// http://erlang.org/doc/reference_manual/patterns.html
//
// Patterns are setup to be composed of other patterns using
// And, Or, Not and list primitives.
//
// Basic Usage
//
//    pattern := match.Pattern("string") // create matcher
//    err := pattern.Match(someValue)   // check if pattern matches
//
// Basic Types
//
// The package currently provides matches for JSON types: bool, string
// and number.  Note that the Number() type is bound to float64 (as
// that is how typical JSON numbers are parsed in Golang).
//
// Type checks
//
// The `Number()`, `String()`, `Bool()` matchers only check if the
// type matches (no value comparisons).  These are often useful for
// validations.
//
// Composition
//
// Boolean operators And/Or/Not are available for composition
//
//    pattern := match.Pattern("hello").Or("world")
//    // now pattern matches both hello and world
//
// List matching
//
// Lists are JSON arrays (typed `[]interface{}`).
//
// In addition to the `List()` type check method, `ListFirst` can be
// used to apply a condition on the first element and the rest of the
// list.
//
// Capturing
//
// A common need is to capture part of the input. For example, the
// check to see if the first  element of a list is "hello" and if so,
// to capture the rest of the list:
//
//    var rest []interface{}
//    pattern := match.ListFirst("hello", &rest)
//    if err := pattern.Match([]interface{}{"hello", "world"}); err != nil {
//       // now rest == []interface{}{"world"}
//    }
//
// Custom matchers
//
// A custom matcher can be built by using the `MatchFunc`
// helper. For instance to check if a list consists only of strings,
// one can do:
//
//     var isStringList match.Matcher
//     isStringList := match.MatchFunc(func(v interface{}) error {
//        pattern := match.ListEmpty().Or(
//          match.ListFirst(match.String(), isStringList)
//        )
//        return pattern.Match(v)
//     }
//
// Type safety
//
// While this package does not use reflection at all, it does allow a
// lot of flexible types.  This is intentional as working with JSON
// often involves heterogenous types.
//
// Other possible applications
//
// A typical situation where these types of pattern matching come in
// handy is working with code transformations: patterns allow catching
// specific AST sub-trees and rewriting them.
//
// Dual use values
//
// Erlang supports the ability for a placeholder ("capture") to act as
// both a unbound capture variable (on first use) and as a exact
// matcher (on subsequent use).
//
// These can be done using Auto().
package match

import (
	"errors"
	"strings"
)

// ErrNoMatch is returned if there is no match
var ErrNoMatch = errors.New("no match")

// ErrNotNumber is returned if a number was expected but input wasn't a number
var ErrNotNumber = errors.New("not a number")

// ErrNotBool is returned if a bool was expected but input wasn't boolean
var ErrNotBool = errors.New("not a bool")

// ErrNotString is returned if a string was expected but input wasn't a string
var ErrNotString = errors.New("not a string")

// ErrNotList is returned if a list was expected but input wasn't a list
var ErrNotList = errors.New("not a list")

// ErrEmptyList is returned if an empty list was found where a non-empty list was expected
var ErrEmptyList = errors.New("list is empty")

// Match is the interface implemented by matchers
type Matcher interface {
	Match(v interface{}) error

	// chainable helper methods
	And(v interface{}) Matcher
	Or(v interface{}) Matcher
	Not() Matcher
}

// MatchFunc converts a function into a matcher
type MatchFunc func(v interface{}) error

func (m MatchFunc) Match(v interface{}) error {
	return m(v)
}

// And succeeds if the values matches both matchers
func (m MatchFunc) And(x interface{}) Matcher {
	return MatchFunc(func(v interface{}) error {
		if err := m(v); err != nil {
			return err
		}
		return Pattern(x).Match(v)
	})
}

// Or succeeds if the values matches either matchers
func (m MatchFunc) Or(x interface{}) Matcher {
	return MatchFunc(func(v interface{}) error {
		if err := m(v); err == nil {
			return nil
		}
		return Pattern(x).Match(v)
	})
}

// Not inverts the match
func (m MatchFunc) Not() Matcher {
	return MatchFunc(func(v interface{}) error {
		if err := m(v); err != nil {
			return nil
		}
		return ErrNoMatch
	})
}

// Any matches anything
func Any() Matcher {
	return MatchFunc(func(v interface{}) error {
		return nil
	})
}

// None matches nothing
func None() Matcher {
	return MatchFunc(func(v interface{}) error {
		return ErrNoMatch
	})
}

// Number matches against a number.
func Number() Matcher {
	return MatchFunc(func(v interface{}) error {
		if _, ok := v.(float64); ok {
			return nil
		}
		return ErrNotNumber
	})
}

// Bool matches against a bool.
func Bool() Matcher {
	return MatchFunc(func(v interface{}) error {
		if _, ok := v.(bool); ok {
			return nil
		}
		return ErrNotBool
	})
}

// String matches against a string.
func String() Matcher {
	return MatchFunc(func(v interface{}) error {
		if _, ok := v.(string); ok {
			return nil
		}
		return ErrNotString
	})
}

// StringPrefix matches a any string for which `s` is a prefix
func StringPrefix(s string) Matcher {
	return MatchFunc(func(v interface{}) error {
		if ss, ok := v.(string); ok && strings.HasPrefix(ss, s) {
			return nil
		}
		return ErrNoMatch
	})
}
