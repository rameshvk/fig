package match

// List matches any list
func List() Matcher {
	return MatchFunc(func(v interface{}) bool {
		_, ok := v.([]interface{})
		return ok
	})
}

// List matches a non-empty list whose first element
// matches the provided first element and the rest match
// the provided rest matcher
func ListFirst(first, rest Matcher) Matcher {
	return MatchFunc(func(v interface{}) bool {
		x, ok := v.([]interface{})
		if !ok || len(x) == 0 {
			return false
		}
		return first.Match(x[0]) && rest.Match(x[1:])
	})
}

// ListLast is like List except last matches the
// last element while rest matches the list of all other
// elements.
func ListLast(last, rest Matcher) Matcher {
	return MatchFunc(func(v interface{}) bool {
		x, ok := v.([]interface{})
		if !ok || len(x) == 0 {
			return false
		}
		l := len(x)
		return last.Match(x[l-1]) && rest.Match(x[:l-1])
	})
}
