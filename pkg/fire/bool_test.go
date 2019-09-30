package fire_test

import (
	"context"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
)

func TestBool(t *testing.T) {
	ctx := context.Background()
	if c := fire.Bool(true).Code(ctx); c != "true" {
		t.Fatal("code failed", c)
	}

	if c := fire.Bool(false).Code(ctx); c != "false" {
		t.Fatal("code failed", c)
	}

	if fire.Bool(true).HashCode() == fire.Bool(false).HashCode() {
		t.Fatal("Hash codes were same")
	}

	v := fire.Bool(true).Lookup(ctx, fire.String("len"))
	if err, ok := v.Error(ctx); !ok || err == nil {
		t.Fatal("bool lookup succeeded")
	}

	if !fire.Bool(true).Equals(ctx, fire.Bool(true)) {
		t.Fatal("identity check failed")
	}

	if fire.Bool(true).Equals(ctx, fire.Bool(false)) {
		t.Fatal("identity check failed")
	}

	if _, ok := fire.Bool(true).Number(ctx); ok {
		t.Fatal("Number check failed")
	}

	if b, ok := fire.Bool(true).Bool(ctx); !ok || !b {
		t.Fatal("Number check failed", ok, b)
	}

	if _, ok := fire.Bool(true).String(ctx); ok {
		t.Fatal("string check failed")
	}

	if _, ok := fire.Bool(true).Error(ctx); ok {
		t.Fatal("error check failed")
	}
}
