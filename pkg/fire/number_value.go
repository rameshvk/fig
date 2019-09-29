package fire

import (
	"context"
	"math"
	"strconv"
)

// Number converts a float64 into a fire Value
func Number(f float64) Value {
	return numberValue(f)
}

type numberValue float64

func (n numberValue) Code(ctx context.Context) string {
	if math.IsNaN(float64(n)) {
		return "math.NaN()"
	}
	if math.IsInf(float64(n), 1) {
		return "math.Inf()"
	}
	if math.IsInf(float64(n), -1) {
		return "-math.Inf()"
	}
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

func (n numberValue) Call(ctx context.Context, args ...Value) Value {
	return errorValue("cannot call a number")
}

func (n numberValue) Lookup(ctx context.Context, field Value) Value {
	// TODO: add string methods
	return errorValue("cannot lookup a number")
}

func (n numberValue) Equals(ctx context.Context, other Value) bool {
	return n == other
}

func (n numberValue) Number(ctx context.Context) (float64, bool) {
	return float64(n), true
}

func (n numberValue) String(ctx context.Context) (string, bool) {
	return "", false
}

func (n numberValue) Bool(ctx context.Context) (bool, bool) {
	return false, false
}

func (n numberValue) Error(ctx context.Context) (error, bool) {
	return nil, false
}
