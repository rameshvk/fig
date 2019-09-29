package fire_test

import (
	"context"
	"math"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
	"github.com/rameshvk/fig/pkg/parse"
)

func TestSimple(t *testing.T) {
	sampleObject := fire.Object(map[fire.Value]fire.Value{
		fire.String("five"):       fire.Number(5),
		fire.String("square"):     fire.Function(nocode, squaref),
		fire.String("hypotenuse"): fire.Function(nocode, hypotenuse),
	})
	globals := fire.Scope(
		fire.Globals(),
		[2]fire.Value{fire.String("x"), fire.String("hello")},
		[2]fire.Value{fire.String("o"), sampleObject},
	)

	ctx := context.Background()
	suite := map[string]fire.Value{
		"1 + 2":              fire.Number(3),
		"1 < 5":              fire.Bool(true),
		"1 <= 1":             fire.Bool(true),
		"1 == 1":             fire.Bool(true),
		"1 > 1":              fire.Bool(false),
		"2 >= 1":             fire.Bool(true),
		"1 - 5":              fire.Number(-4),
		"5 / 5":              fire.Number(1),
		"5 * 5":              fire.Number(25),
		"4 != 4":             fire.Bool(false),
		"4 + (-4)":           fire.Number(0),
		"true == true":       fire.Bool(true),
		"true == 5":          fire.Bool(false),
		`"hello" == "hello"`: fire.Bool(true),
		"true & true":        fire.Bool(true),
		"false & x()":        fire.Bool(false),
		"true & false":       fire.Bool(false),
		"!(1 < 5)":           fire.Bool(false),
		"false | true":       fire.Bool(true),
		"true | x()":         fire.Bool(true),
		"false | false":      fire.Bool(false),
		"x":                  fire.String("hello"),
		"o.five + 10":        fire.Number(15),
		"o.square(5)+2":      fire.Number(27),
		"o.square(x, where(x = 1 + y, y = z = 4))": fire.Number(25),
		"o.hypotenuse(x=3, y=4)":                   fire.Number(5),
		"o.hypotenuse(x, y, where(x=3, y=4))":      fire.Number(5),
		"o.hypotenuse(x, y=z, where(x=3, z=4))":    fire.Number(5),
		"{ it }(2)":                                fire.Number(2),
		"{ it }(x = 5, y = 10).x":                  fire.Number(5),
	}

	for k, v := range suite {
		parsed, err := parse.String(k)
		if err != nil {
			t.Fatal("Unexpected parse error", err)
		}
		got := fire.Eval(ctx, parsed, globals)
		if !got.Equals(ctx, v) {
			t.Fatal("Mismatched values", k, v, got)
		}
	}
}

func nocode(ctx context.Context) string {
	panic("no code for builtin functions?")
}

func squaref(ctx context.Context, args ...fire.Value) fire.Value {
	if f, ok := args[0].Number(ctx); ok {
		return fire.Number(f * f)
	}
	return fire.Error("arg not a number")
}

func hypotenuse(ctx context.Context, args ...fire.Value) fire.Value {
	x, okx := args[0].Lookup(ctx, fire.String("x")).Number(ctx)
	y, oky := args[0].Lookup(ctx, fire.String("y")).Number(ctx)
	if !okx || !oky {
		return fire.Error("invalid args")
	}
	return fire.Number(math.Sqrt(x*x + y*y))
}
