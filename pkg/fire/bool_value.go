package fire

import (
	"context"
)

type boolValue bool

func (b boolValue) Code(ctx context.Context) string {
	if bool(b) {
		return "true"
	}
	return "false"
}

func (b boolValue) Call(ctx context.Context, args ...Value) Value {
	return errorValue("cannot call a bool")
}

func (b boolValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add string methods
	return errorValue("cannot lookup a bool")
}

func (b boolValue) Equals(ctx context.Context, other Value) bool {
	return b == other
}

func (b boolValue) Number(ctx context.Context) (float64, bool) {
	return 0, false
}

func (b boolValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (b boolValue) Bool(ctx context.Context) (bool, bool) {
	return bool(b), true
}

func (b boolValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
