package match

// Pattern matches the input with the provided value
//
// string, bool and float64 values are matched exactly
// Arrays must match in size and the individual elements
// are matched by recurively calling Pattern on them.
// nil values are also matched exactly (i.e. input must be nil)
//
// *string, *bool and *float64 are considered "capture"
// values -- the input type must match the base type and if
// so the values are copied into the provided pointers.
// *[]interface{} and generic *interface{} are also considered
// pointer types.
//
// If the argument is a matcher, that is invoked instead.
//
// Any other value provided to Pattern results in a ErrNoMatch
func Pattern(pattern interface{}) Matcher {
	if m, ok := pattern.(Matcher); ok {
		return m
	}

	check := func(v bool) error {
		if v {
			return nil
		}
		return ErrNoMatch
	}
	return MatchFunc(func(v interface{}) error {
		switch pattern := pattern.(type) {
		case nil:
			return check(v == nil)
		case string:
			return check(pattern == v)
		case *string:
			if s, ok := v.(string); ok {
				*pattern = s
				return nil
			}
		case float64:
			return check(pattern == v)
		case *float64:
			if f, ok := v.(float64); ok {
				*pattern = f
				return nil
			}
		case bool:
			return check(pattern == v)
		case *bool:
			if b, ok := v.(bool); ok {
				*pattern = b
				return nil
			}
		case []interface{}:
			items, ok := v.([]interface{})
			if !ok || len(items) != len(pattern) {
				return ErrNoMatch
			}

			for kk, elt := range pattern {
				if err := Pattern(elt).Match(items[kk]); err != nil {
					return err
				}
			}
			return nil
		case *[]interface{}:
			if l, ok := v.([]interface{}); ok {
				*pattern = l
				return nil
			}
		case *interface{}:
			*pattern = v
			return nil
		}
		return ErrNoMatch
	})
}
