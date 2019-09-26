package match

// Equals matches the input with the provided value
//
// String, Bool and Float64 are matched exactly
// Arrays must match in size and the corresponding
// elements must match.
//
// If the argument is a matcher, that is invoked instead.
// Similarly, if an array element is a  matcher, that is invoked
// instead.
func Equals(pattern interface{}) Matcher {
	return MatchFunc(func(v interface{}) bool {
		switch pattern := pattern.(type) {
		case nil:
			return v == nil
		case Matcher:
			return pattern.Match(v)
		case string:
			return pattern == v
		case float64:
			return pattern == v
		case bool:
			return pattern == v
		case []interface{}:
			items, ok := v.([]interface{})
			if !ok || len(items) != len(pattern) {
				return false
			}

			for kk, elt := range pattern {
				if !Equals(elt).Match(items[kk]) {
					return false
				}
			}
			return true
		}
		return false
	})
}
