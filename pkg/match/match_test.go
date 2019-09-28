package match_test

import (
	"fmt"
	"testing"

	"github.com/rameshvk/fig/pkg/match"
)

func Example() {
	v := []interface{}{
		float64(1.0),
		"hello",
		true,
		nil,
		[]interface{}{},
	}
	var n, boo, s, b interface{}
	var nx float64
	var sx string
	var bx bool

	// must match against itself
	pattern := match.Equals(v)

	// must much with matchers for each element
	pattern = pattern.And([]interface{}{
		match.Number().And(&n).And(&nx),
		match.Number().And(&boo).Or(
			match.String().And(match.StringPrefix("hel")).And(&s).And(&sx),
		),
		match.Number().Not().And(match.Bool().And(&b).And(&bx)),
		nil,
		match.List().And([]interface{}{1}).Not(),
	})

	var minusHead, minusTail []interface{}

	// also do some more generic match patterns
	pattern = pattern.
		And(match.ListFirst(match.Number(), &minusHead)).
		And(match.ListLast(match.List(), &minusTail)).
		And(match.None().Not().And(match.Any()))

	if err := match.Equals(pattern).Match(v); err != nil {
		fmt.Println("Did not match", err)
	} else {
		fmt.Println("Got", n, boo, s, b, nx, sx, bx)
	}

	if err := match.Equals(minusHead).Match(v[1:]); err != nil {
		fmt.Println("minus head failed", err)
	}

	if err := match.Equals(minusTail).Match(v[:len(v)-1]); err != nil {
		fmt.Println("minus tail failed", err)
	}

	// Output: Got 1 <nil> hello true 1 hello true
}

func TestAuto(t *testing.T) {
	var b bool
	var f float64
	var s string
	var l []interface{}
	var v interface{}

	autob := match.Auto(&b)
	if err := match.Equals([]interface{}{autob, autob}).Match([]interface{}{true, true}); err != nil {
		t.Fatal("did not match as expected", err)
	}

	autof := match.Auto(&f)
	if err := match.Equals([]interface{}{autof, autof}).Match([]interface{}{1.0, 1.0}); err != nil {
		t.Fatal("did not match as expected", err)
	}

	autos := match.Auto(&s)
	if err := match.Equals([]interface{}{autos, autos}).Match([]interface{}{"ok", "ok"}); err != nil {
		t.Fatal("did not match as expected", err)
	}

	autol := match.Auto(&l)
	ll := []interface{}{"boo"}
	if err := match.Equals([]interface{}{autol, autol}).Match([]interface{}{ll, ll}); err != nil {
		t.Fatal("did not match as expected", err)
	}

	autov := match.Auto(&v)
	if err := match.Equals([]interface{}{autov, autov}).Match([]interface{}{1.0, 1.0}); err != nil {
		t.Fatal("did not match as expected", err)
	}

	if err := match.Equals([]interface{}{autov, autov}).Match([]interface{}{1.0, 2.0}); err == nil {
		t.Fatal("did not match as expected", err)
	}

	if err := match.Auto(t).Match(t); err == nil {
		t.Fatal("Unexpected success", err)
	}
}

func TestEdgeCases(t *testing.T) {
	if err := match.Equals(t).Match(t); err == nil {
		t.Fatal("Unexpected match")
	}

	var input = []interface{}{"x"}
	var pattern = []interface{}{true}
	if err := match.Equals(pattern).Match(input); err == nil {
		t.Fatal("Unexpected match")
	}

	if err := match.List().Match(true); err != match.ErrNotList {
		t.Fatal("Unexpected match", err)
	}

	if err := match.ListFirst(nil, nil).Match([]interface{}{}); err != match.ErrEmptyList {
		t.Fatal("Unexpected match", err)
	}

	if err := match.ListFirst(5, nil).Match([]interface{}{10}); err != match.ErrNoMatch {
		t.Fatal("Unexpected match", err)
	}

	if err := match.ListLast(nil, nil).Match([]interface{}{}); err != match.ErrEmptyList {
		t.Fatal("Unexpected match", err)
	}

	if err := match.ListLast(5, nil).Match([]interface{}{10}); err != match.ErrNoMatch {
		t.Fatal("Unexpected match", err)
	}

	if err := match.String().Match(5.0); err != match.ErrNotString {
		t.Fatal("Unexpected match", err)
	}

	if err := match.Bool().Match(5.0); err != match.ErrNotBool {
		t.Fatal("Unexpected match", err)
	}

	if err := match.Number().Match("boo"); err != match.ErrNotNumber {
		t.Fatal("Unexpected match", err)
	}

	if err := match.String().Not().Match("boo"); err != match.ErrNoMatch {
		t.Fatal("unexpected match", err)
	}

	if err := match.String().Or(match.Number()).Match("s"); err != nil {
		t.Fatal("unexpected match", err)
	}

	if err := match.String().Or(match.Number()).Match(5.0); err != nil {
		t.Fatal("unexpected match", err)
	}

	if err := match.StringPrefix("boo").Match("Boo"); err != match.ErrNoMatch {
		t.Fatal("Unexpected match", err)
	}
}
