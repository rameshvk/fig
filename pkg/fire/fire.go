// Package fire is the fig runtime
package fire

import (
	"context"
)

type Value interface {
	// generates the fig language representation of the value.
	Code(ctx context.Context) string

	// HashCode is any code useful for avoiding collisions when
	// storing the value as a key. Two values that are equal must
	// have the same hash code.
	HashCode() interface{}

	// Call with the provided args.
	Call(ctx context.Context, args ...Value) Value

	// Lookup looks up the field/key
	Lookup(ctx context.Context, field Value) Value

	// Equals compares one value against another
	Equals(ctx context.Context, other Value) bool

	// Number converts the value to a float64 if possible
	Number(ctx context.Context) (float64, bool)

	// String converts the value to a string if possible
	String(ctx context.Context) (string, bool)

	// Bool converts the value to a bool if possible
	Bool(ctx context.Context) (bool, bool)

	// ToError converts the value to an error if possible
	Error(ctx context.Context) (error, bool)
}
