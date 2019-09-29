package fire

import (
	"bytes"
	"context"
)

// String creates a string value
func String(s string) Value {
	return stringValue(s)
}

type stringValue string

func (s stringValue) Code(ctx context.Context) string {
	check := func(n int, err error) {
		if err != nil {
			panic(err)
		}
	}
	var buf bytes.Buffer
	check(buf.WriteRune('"'))
	for _, r := range string(s) {
		if r == '"' {
			check(buf.WriteRune('\\'))
		}
		check(buf.WriteRune(r))
	}
	check(buf.WriteRune('"'))
	return buf.String()
}

func (s stringValue) Call(ctx context.Context, args ...Value) Value {
	return errorValue("cannot call a string")
}

func (s stringValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add string methods
	return errorValue("cannot lookup a string")
}

func (s stringValue) Equals(ctx context.Context, other Value) bool {
	return s == other
}

func (s stringValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (s stringValue) String(ctx context.Context) (string, bool) {
	return string(s), true
}

func (s stringValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (s stringValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
