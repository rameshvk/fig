package fire_test

import (
	"context"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
)

func TestString(t *testing.T) {
	ctx := context.Background()
	if c := fire.String("hello").Code(ctx); c != `"hello"` {
		t.Fatal("code failed", c)
	}

	if c := fire.String("hello\"").Code(ctx); c != `"hello\""` {
		t.Fatal("code failed", c)
	}

	if fire.String("ok").HashCode() == fire.String("yo").HashCode() {
		t.Fatal("Hash codes were same")
	}

	v := fire.String("hello").Lookup(ctx, fire.String("len"))
	if err, ok := v.Error(ctx); !ok || err == nil {
		t.Fatal("string lookup succeeded")
	}

	if !fire.String("ok").Equals(ctx, fire.String("ok")) {
		t.Fatal("identity check failed")
	}

	if fire.String("ok").Equals(ctx, fire.String("yo")) {
		t.Fatal("identity check failed")
	}

	if _, ok := fire.String("hello").Number(ctx); ok {
		t.Fatal("Number check failed")
	}

	if s, ok := fire.String("hello").String(ctx); !ok || s != "hello" {
		t.Fatal("string check failed", ok, s)
	}

	if _, ok := fire.String("hello").Bool(ctx); ok {
		t.Fatal("bool check failed")
	}

	if _, ok := fire.String("hello").Error(ctx); ok {
		t.Fatal("error check failed")
	}
}
