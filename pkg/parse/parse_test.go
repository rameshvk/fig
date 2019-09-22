package parse_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/rameshvk/fig/pkg/parse"

	"testing"
)

func name(s, loc string) []interface{} {
	return []interface{}{"name:" + loc, s}
}

func number(f float64, loc string) []interface{} {
	return []interface{}{"number:" + loc, f}
}

func TestComplex(t *testing.T) {
	code := `f.g(x + y*3.2,"hell\"o").h + 1000`
	xplusy32 := []interface{}{
		"+:6:7",
		name("x", "4:5"),
		[]interface{}{"*:9:10", name("y", "8:9"), number(3.2, "10:13")},
	}
	gcall := []interface{}{
		"call:3:3",
		name("g", "2:3"),
		xplusy32,
		[]interface{}{"string:14:23", `hell"o`},
	}
	fdotg := []interface{}{".:1:2", name("f", "0:1"), gcall}
	fdotgdoth := []interface{}{".:24:25", fdotg, name("h", "25:26")}
	expected := []interface{}{"+:27:28", fdotgdoth, number(1000, "29:33")}

	got, errs := parse.String(code)
	if len(errs) > 0 {
		t.Fatal("Failed", errs)
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatal("Diff", diff)
	}
}

func TestFunctionNoArgs(t *testing.T) {
	code := "f()"
	expected := []interface{}{"call:1:1", name("f", "0:1")}

	got, errs := parse.String(code)
	if len(errs) > 0 {
		t.Fatal("Failed", errs)
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatal("Diff", diff)
	}
}

func TestUnaryOps(t *testing.T) {
	code := "!(-1 < +x)"
	expected := []interface{}{
		"!:0:1",
		[]interface{}{
			"<:5:6",
			[]interface{}{"-:2:3", number(1, "3:4")},
			[]interface{}{"+:7:8", name("x", "8:9")},
		},
	}

	got, errs := parse.String(code)
	if len(errs) > 0 {
		t.Fatal("Failed", errs)
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatal("Diff", diff)
	}
}

func TestAssignment(t *testing.T) {
	code := "f(x = y = a + b, p = q)"
	firstArg := []interface{}{
		"=:4:5",
		name("x", "2:3"),
		[]interface{}{
			"=:8:9",
			name("y", "6:7"),
			[]interface{}{"+:12:13", name("a", "10:11"), name("b", "14:15")},
		},
	}

	expected := []interface{}{
		"call:1:1",
		name("f", "0:1"),
		firstArg,
		[]interface{}{"=:19:20", name("p", "17:18"), name("q", "21:22")},
	}

	got, errs := parse.String(code)
	if len(errs) > 0 {
		t.Fatal("Failed", errs)
	}
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatal("Diff", diff)
	}
}
