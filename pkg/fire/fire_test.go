package fire_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
	"github.com/rameshvk/fig/pkg/parse"
)

func TestSimple(t *testing.T) {
	ctx := context.Background()
	sampleObject := fire.Object(map[fire.Value]fire.Value{
		fire.String("five"):       fire.Number(5),
		fire.String("square"):     fire.Function(nocode, squaref),
		fire.String("hypotenuse"): fire.Function(nocode, hypotenuse),
	})
	globals := fire.Scope(
		ctx,
		fire.Globals(),
		[2]fire.Value{fire.String("x"), fire.String("hello")},
		[2]fire.Value{fire.String("o"), sampleObject},
		[2]fire.Value{fire.String("panic"), fire.Function(nocode, panicf)},
	)

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

		// test that if short-circuits properly
		`if(5 < 10, "good", panic())`:       fire.String("good"),
		`if(error("boo"), panic(), "good")`: fire.String("good"),

		// test short circuited & and |
		`error("boo") & panic()`: fire.Error("boo"),
		`"hello" | panic()`:      fire.String("hello"),

		// test object and error
		`object(x = 5).x`: fire.Number(5),
		`error("hello")`:  fire.Error("hello"),
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

func TestWrap(t *testing.T) {
	values := []interface{}{
		"hello",
		true,
		5.0,
		errors.New("some error"),
		map[interface{}]interface{}{"hello": 5.0},
	}
	ctx := context.Background()
	for _, value := range values {
		dupe := fire.ToNative(ctx, fire.FromNative(ctx, value))
		errdupe, ok := dupe.(error)
		errval, ok2 := value.(error)
		if ok || ok2 {
			if errdupe.Error() != errval.Error() {
				t.Fatal("Mismatched error types", errdupe, errval)
			}
		} else if !reflect.DeepEqual(value, dupe) {
			t.Fatal("Mismatched", value, dupe)
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

func panicf(ctx context.Context, args ...fire.Value) fire.Value {
	panic("unexpected")
}
