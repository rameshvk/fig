package parse

import (
	"strconv"
	"strings"
)

func normalize(v interface{}, errs *[]error) interface{} {
	if err, ok := v.(error); ok {
		*errs = append(*errs, err)
		return nil
	}

	list, ok := v.([]interface{})
	if !ok {
		return v
	}
	fn := list[0].(string)
	loc := list[1].(string)
	var left, right interface{}
	if len(list) > 2 {
		left = list[2]
	}
	if len(list) > 3 {
		right = list[3]
	}

	switch fn {
	case "":
		l := []interface{}{"call:" + loc, normalize(left, errs)}
		return appendArgs(l, loc, right, errs)
	case "()":
		return normalize(left, errs)
	case "{}":
		panic("NYI") // need to allow comma and equals here
	case ",":
		panic("NYI") // need error code here
	case "=":
		panic("NYI") // need error code here
	case "string", "name", "number":
		return []interface{}{fn + ":" + loc, left}
	case "+", "-", "!":
		if _, ok := left.(MissingTermError); ok {
			// accept as unary operator
			return []interface{}{fn + ":" + loc, normalize(right, errs)}
		}
	}

	return []interface{}{fn + ":" + loc, normalize(left, errs), normalize(right, errs)}
}

func appendArgs(list []interface{}, loc string, args interface{}, errs *[]error) interface{} {
	l, ok := args.([]interface{})
	if !ok || l[0] != "()" {
		parts := strings.Split(loc, ":")
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		panic(args)
		*errs = append(*errs, MissingOperatorError(end))
		return nil
	}

	if _, ok := l[2].(MissingTermError); ok {
		// no args, that is ok
		return list
	}

	return appendCommas(list, l[2], errs)
}

func appendCommas(list []interface{}, l interface{}, errs *[]error) interface{} {
	if comma, ok := l.([]interface{}); ok && comma[0] == "," {
		list = append(list, normalizeArg(comma[2], errs))
		return appendCommas(list, comma[3], errs)
	}
	return append(list, normalizeArg(l, errs))
}

func normalizeArg(v interface{}, errs *[]error) interface{} {
	// equals is valid as an arg though lhs has constraints
	if assign, ok := v.([]interface{}); ok && assign[0] == "=" {
		loc := assign[1].(string)
		return []interface{}{"=:" + loc, normalizeScopeLHS(assign[2], errs), normalize(assign[3], errs)}
	}

	return normalize(v, errs)
}

func normalizeScopeLHS(v interface{}, errs *[]error) interface{} {
	if n, ok := v.([]interface{}); ok && n[0] == "name" {
		return normalize(v, errs)
	}

	if assign, ok := v.([]interface{}); ok && assign[0] == "=" {
		loc := assign[1].(string)
		return []interface{}{"=:" + loc, normalizeScopeLHS(assign[2], errs), normalize(assign[3], errs)}
	}

	panic("NYI") // allow .x type expressions too but nothing else
}
