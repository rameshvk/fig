package eval_test

import (
	"testing"

	"github.com/rameshvk/fig/pkg/eval"
)

var intmap = map[int]int{4: 24}
var strmap = map[string]interface{}{"hoo": "hoodoo"}
var scope = eval.ExtendScope(map[string]interface{}{
	"boo":    45.0,
	"hoo":    42,
	"boohoo": eval.ExtendScope(strmap,nil),
	"map": eval.Reflect(intmap, eval.Reflect(strmap, nil, nil), nil),
	"var": eval.CallableFunc(func(root eval.Scope, offset int, args []interface{}) (interface{}, error) {
		return root.Lookup(root, offset, args[0].(string))
	}),
}, eval.DefaultScope)

func TestExpressions(t *testing.T) {
	exprs := map[string]interface{}{
		`["+", ["var", "boo"], ["var", "hoo"]]`:  45.0 + 42,
		`["-", ["var", "boo"], ["var", "hoo"]]`:  45.0 - 42,
		`["*", ["var", "boo"], ["var", "hoo"]]`:  45.0 * 42,
		`["/", ["var", "boo"], ["var", "hoo"]]`:  45.0 / 42,
		`["<", ["var", "boo"], ["var", "hoo"]]`:  45.0 < 42,
		`[">", ["var", "boo"], ["var", "hoo"]]`:  45.0 > 42,
		`["<=", ["var", "boo"], ["var", "hoo"]]`: 45.0 <= 42,
		`[">=", ["var", "boo"], ["var", "hoo"]]`: 45.0 >= 42,
		`["==", ["var", "boo"], ["var", "hoo"]]`: 45.0 == 42,
		`["!=", ["var", "boo"], ["var", "hoo"]]`: 45.0 != 42,
		`["==", 4, 4]`:                           true,
		`["!=", 4, 4]`:                           false,
		`["==", "hello", "hello"]`:               true,
		`["!=", "hello", "hello"]`:               false,
		`["&&", true, false]`:                    false,
		`["&&", true, true]`:                     true,
		`["||", true, false]`:                    true,
		`["||", false, false]`:                   false,
		`["!", true]`:                            false,
		`["if", true, 1, 2]`:                     1.0,
		`["if", false, 1, 2]`:                    2.0,
		`[".", ["var", "boohoo"], "hoo"]`:        "hoodoo",
		`[".", ["var", "map"], "4"]`:           24,
		`[".", ["var", "map"], "hoo"]`:           "hoodoo",		
	}

	for k, v := range exprs {
		t.Run(k, func(t *testing.T) {
			got, err := eval.Encoded(k, scope)
			if err != nil || got != v {
				t.Fatal("failed", err, got, v)
			}
		})
	}
}
