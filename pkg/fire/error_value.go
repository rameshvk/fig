package fire

import (
	"context"
)

// Error converts a string into an error value
func Error(s string) Value {
	return errorValue(s)
}

type errorValue string

func (e errorValue) Code(ctx context.Context) string {
	return `error(` + stringValue(e).Code(ctx) + `)`
}

func (e errorValue) HashCode() interface{} {
	return e
}

func (e errorValue) Call(ctx context.Context, args ...Value) Value {
	return e
}

func (e errorValue) Lookup(ctx context.Context, field Value) Value {
	return e
}

func (e errorValue) Equals(ctx context.Context, other Value) bool {
	return e == other
}

func (e errorValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (e errorValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (e errorValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (e errorValue) Error(ctx context.Context) (error, bool) {
	return fireError(string(e)), true
}

type fireError string

func (f fireError) Error() string {
	return string(f)
}
