package fire_test

import (
	"context"
	"testing"

	"github.com/rameshvk/fig/pkg/fire"
)

func TestError(t *testing.T) {
	ctx := context.Background()
	if c := fire.Error("ok").Code(ctx); c != `error("ok")` {
		t.Fatal("code failed", c)
	}

	if c := fire.Error("hello\"").Code(ctx); c != `error("hello\"")` {
		t.Fatal("code failed", c)
	}

	if fire.Error("ok").HashCode() == fire.Error("yo").HashCode() {
		t.Fatal("Hash codes were same")
	}

	v := fire.Error("hello").Lookup(ctx, fire.String("len"))
	if v != fire.Error("hello") {
		t.Fatal("string lookup succeeded", v)
	}

	if !fire.Error("ok").Equals(ctx, fire.Error("ok")) {
		t.Fatal("identity check failed")
	}

	if fire.Error("ok").Equals(ctx, fire.Error("yo")) {
		t.Fatal("identity check failed")
	}

	if _, ok := fire.Error("hello").Number(ctx); ok {
		t.Fatal("Number check failed")
	}

	if _, ok := fire.Error("hello").String(ctx); ok {
		t.Fatal("string check failed", ok)
	}

	if err, ok := fire.Error("hello").Error(ctx); !ok || err.Error() != "hello" {
		t.Fatal("bool check failed", err, ok)
	}

	if _, ok := fire.Error("hello").Bool(ctx); ok {
		t.Fatal("error check failed")
	}
}
