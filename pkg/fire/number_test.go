package fire_test

import (
	"context"
	"math"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
)

func TestNumber(t *testing.T) {
	ctx := context.Background()
	if c := fire.Number(1.5).Code(ctx); c != `1.5` {
		t.Fatal("code failed", c)
	}

	if c := fire.Number(math.Inf(+1)).Code(ctx); c != `math.Inf` {
		t.Fatal("code failed", c)
	}

	if c := fire.Number(math.Inf(-1)).Code(ctx); c != `-math.Inf` {
		t.Fatal("code failed", c)
	}

	if c := fire.Number(math.NaN()).Code(ctx); c != `math.NaN` {
		t.Fatal("code failed", c)
	}

	if fire.Number(1.5).HashCode() == fire.Number(2).HashCode() {
		t.Fatal("Hash codes were same")
	}

	v := fire.Number(1.5).Lookup(ctx, fire.String("len"))
	if err, ok := v.Error(ctx); !ok || err == nil {
		t.Fatal("string lookup succeeded", v)
	}

	if !fire.Number(1.5).Equals(ctx, fire.Number(1.5)) {
		t.Fatal("identity check failed")
	}

	if fire.Number(1.5).Equals(ctx, fire.Number(2)) {
		t.Fatal("identity check failed")
	}

	if n, ok := fire.Number(1.5).Number(ctx); !ok || n != 1.5 {
		t.Fatal("Number check failed", n, ok)
	}

	if _, ok := fire.Number(1.5).String(ctx); ok {
		t.Fatal("string check failed", ok)
	}

	if _, ok := fire.Number(1.5).Error(ctx); ok {
		t.Fatal("error check failed")
	}

	if _, ok := fire.Number(1.5).Bool(ctx); ok {
		t.Fatal("bool check failed")
	}
}
