package match

// Auto captures a value on first use and matches on subsequence use
//
// For example, to check if the first and second elements of a list
// have the same value, one can do:
//
//      var v interface{}
//      first := match.Auto(&v)
//      match.Equals([]interface{first, first})
func Auto(ptr interface{}) Matcher {
	matched := false
	var pattern interface{}
	return MatchFunc(func(v interface{}) error {
		if matched {
			return Equals(pattern).Match(v)
		}
		switch ptr.(type) {
		case *bool:
			_, matched = v.(bool)
		case *float64:
			_, matched = v.(float64)
		case *string:
			_, matched = v.(string)
		case *[]interface{}:
			_, matched = v.([]interface{})
		case *interface{}:
			_, matched = v.(interface{})
		}
		if matched {
			pattern = v
			return nil
		}
		return ErrNoMatch
	})
}
