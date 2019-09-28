package match

// List matches any list
func List() Matcher {
	return MatchFunc(func(v interface{}) error {
		if _, ok := v.([]interface{}); ok {
			return nil
		}
		return ErrNotList
	})
}

// List matches a non-empty list whose first element
// matches the provided first element and the rest match
// the provided rest matcher
func ListFirst(f, r interface{}) Matcher {
	first, rest := Pattern(f), Pattern(r)
	return MatchFunc(func(v interface{}) error {
		x, ok := v.([]interface{})
		if !ok || len(x) == 0 {
			return ErrEmptyList
		}
		if err := first.Match(x[0]); err != nil {
			return err
		}
		return rest.Match(x[1:])
	})
}

// ListLast is like List except last matches the
// last element while rest matches the list of all other
// elements.
func ListLast(l, r interface{}) Matcher {
	last, rest := Pattern(l), Pattern(r)
	return MatchFunc(func(v interface{}) error {
		x, ok := v.([]interface{})
		if !ok || len(x) == 0 {
			return ErrEmptyList
		}
		l := len(x)
		if err := last.Match(x[l-1]); err != nil {
			return err
		}
		return rest.Match(x[:l-1])
	})
}
