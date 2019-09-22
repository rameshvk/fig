package parse_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/rameshvk/fig/pkg/parse"

	"testing"
)

func name(s, loc string) []interface{} {
	return []interface{}{"name", loc, s}
}

func number(f float64, loc string) []interface{} {
	return []interface{}{"number", loc, f}
}

func TestSuccess(t *testing.T) {
	code := `f.g(x + y*3.2,"hell\"o").h + 1000`
	xplusy32 := []interface{}{
		"+",
		"6:7",
		name("x", "4:5"),
		[]interface{}{"*", "9:10", name("y", "8:9"), number(3.2, "10:13")},
	}
	args := []interface{}{
		",",
		"13:14",
		xplusy32,
		[]interface{}{"string", "14:23", `hell"o`},
	}
	gcall := []interface{}{
		"",
		"3:3",
		name("g", "2:3"),
		[]interface{}{"()", "3:24", args},
	}
	fdotg := []interface{}{".", "1:2", name("f", "0:1"), gcall}
	fdotgdoth := []interface{}{".", "24:25", fdotg, name("h", "25:26")}
	expected := []interface{}{"+", "27:28", fdotgdoth, number(1000, "29:33")}

	got, errs := parse.String(code)
	if len(errs) > 0 {
		t.Fatal("Failed", errs)
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatal("Diff", diff)
	}
}
