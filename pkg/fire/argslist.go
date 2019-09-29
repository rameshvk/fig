package fire

import (
	"github.com/rameshvk/fig/pkg/match"
)

func filterArgs(args []interface{}) ([]interface{}, map[Value]interface{}, Value) {
	rest := []interface{}{}
	where := map[Value]interface{}{}
	for _, arg := range args {
		if wargs := whereArgs(arg); wargs != nil {
			for _, warg := range wargs {
				if err := accumulateAssign(where, warg); err != nil {
					return nil, nil, err
				}
			}
		} else {
			rest = append(rest, arg)
		}
	}
	return rest, where, nil
}

func whereArgs(v interface{}) []interface{} {
	var args []interface{}
	err := match.ListFirst(
		match.StringPrefix("call"),
		match.ListFirst(
			[]interface{}{match.StringPrefix("name"), "where"},
			&args,
		)).Match(v)
	if err != nil {
		return nil
	}
	return args
}

func accumulateAssign(result map[Value]interface{}, arg interface{}) Value {
	names := []string{}
	for assignPattern.Match(arg) == nil {
		var name string
		err := match.Pattern([]interface{}{
			match.StringPrefix("="),
			[]interface{}{match.StringPrefix("string"), &name},
			&arg,
		}).Match(arg)
		if err != nil {
			return errorValue(err.Error())
		}
		names = append(names, name)
	}

	if len(names) == 0 {
		return errorValue("missing = in where")
	}
	for kk, name := range names {
		key := stringValue(name)
		if _, ok := result[key]; ok {
			return errorValue("duplicate name: " + name)
		}
		if kk == 0 {
			result[key] = arg
		} else {
			result[key] = stringValue(names[kk-1])
		}
	}
	return nil
}

var assignPattern = match.Pattern([]interface{}{match.StringPrefix("="), match.Any(), match.Any()})
