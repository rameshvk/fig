package match_test

import (
	"fmt"

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
	pattern := match.Equals(v)
	pattern = pattern.And([]interface{}{
		match.Number().And(match.Capture(&n)),
		match.Number().And(match.Capture(&boo)).Or(
			match.String().
				And(match.StringPrefix("hel")).
				And(match.Capture(&s)),
		),
		match.Number().Not().And(match.Bool().And(match.Capture(&b))),
		nil,
		match.List().And([]interface{}{1}).Not(),
	})

	pattern = pattern.
		And(match.ListFirst(match.Number(), match.List())).
		And(match.ListLast(match.List(), match.List())).
		And(match.None().Not().And(match.Any()))

	if !match.Equals(pattern).Match(v) {
		fmt.Println("Did not match")
	} else {
		fmt.Println("Got", n, boo, s, b)
	}

	// Output: Got 1 <nil> hello true
}
