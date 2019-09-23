package parse_test

import (
	"github.com/rameshvk/fig/pkg/parse"
	"github.com/tvastar/test"

	"testing"
)

func TestSuite(t *testing.T) {
	cases := []string{
		"x + y + z",
		"f(x = y = z, a = b)",
		`"hel\"lo".length`,
		"f()",
		"f.g().h(5).k",
		"f(x = g(), where (g = h))",
		`f.g(x + y*3.2,"hell\"o").h + 1000`,
		"!(x == y)",
		"!(-1 < +x)",
		"x < y & y < z | boo",
		"!!x",
		"x.(y)",
		"{ x, y = 23 }",
		"f({ z.x, z = g() })",
		"{x}()",
	}

	results := map[string]interface{}{}
	for _, code := range cases {
		result, errs := parse.String(code)
		if len(errs) == 0 {
			results[code] = result
		} else {
			results[code] = map[string]interface{}{
				"result": result,
				"errors": errs,
			}
		}
	}
	test.Artifact(t.Error, "parse.json", results)
}
