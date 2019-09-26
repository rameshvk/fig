// Package match has helpful function to match JSON
package match

import "strings"

// Match is the interface implemented by matchers
type Matcher interface {
	Match(v interface{}) bool

	// chainable helper methods
	And(v interface{}) Matcher
	Or(v interface{}) Matcher
	Not() Matcher
}

// MatchFunc converts a function into a matcher
type MatchFunc func(v interface{}) bool

func (m MatchFunc) Match(v interface{}) bool {
	return m(v)
}

// And succeeds if the values matches both matchers
func (m MatchFunc) And(x interface{}) Matcher {
	return MatchFunc(func(v interface{}) bool {
		return m(v) && Equals(x).Match(v)
	})
}

// Or succeeds if the values matches either matchers
func (m MatchFunc) Or(x interface{}) Matcher {
	return MatchFunc(func(v interface{}) bool {
		return m(v) || Equals(x).Match(v)
	})
}

// Not inverts the match
func (m MatchFunc) Not() Matcher {
	return MatchFunc(func(v interface{}) bool {
		return !m(v)
	})
}

// Any matches anything
func Any() Matcher {
	return MatchFunc(func(v interface{}) bool {
		return true
	})
}

// None matches nothing
func None() Matcher {
	return MatchFunc(func(v interface{}) bool {
		return false
	})
}

// Number matches against a number.
func Number() Matcher {
	return MatchFunc(func(v interface{}) bool {
		_, ok := v.(float64)
		return ok
	})
}

// Bool matches against a bool.
func Bool() Matcher {
	return MatchFunc(func(v interface{}) bool {
		_, ok := v.(bool)
		return ok
	})
}

// String matches against a string.
func String() Matcher {
	return MatchFunc(func(v interface{}) bool {
		_, ok := v.(string)
		return ok
	})
}

// StringPrefix matches a any string for which `s` is a prefix
func StringPrefix(s string) Matcher {
	return MatchFunc(func(v interface{}) bool {
		ss, ok := v.(string)
		return ok && strings.HasPrefix(ss, s)
	})
}

// Capture matches any thing but also copies the value
func Capture(pv *interface{}) Matcher {
	return MatchFunc(func(v interface{}) bool {
		*pv = v
		return true
	})
}
