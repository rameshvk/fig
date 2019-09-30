package fire

import (
	"context"
)

// ToNative unwraps a value into native Go types
func ToNative(ctx context.Context, v Value) interface{} {
	if s, ok := v.String(ctx); ok {
		return s
	}
	if b, ok := v.Bool(ctx); ok {
		return b
	}
	if n, ok := v.Number(ctx); ok {
		return n
	}
	if err, ok := v.Error(ctx); ok {
		return err
	}
	if o, ok := v.(obj); ok {
		result := map[interface{}]interface{}{}
		for k, v := range o {
			result[ToNative(ctx, k)] = ToNative(ctx, v)
		}
		return result
	}
	return nil
}

// FromNative wraps an interface into a Value
func FromNative(ctx context.Context, v interface{}) Value {
	switch v := v.(type) {
	case string:
		return String(v)
	case float64:
		return Number(v)
	case bool:
		return Bool(v)
	case error:
		return Error(v.Error())
	case map[interface{}]interface{}:
		result := map[Value]Value{}
		for k, val := range v {
			result[FromNative(ctx, k)] = FromNative(ctx, val)
		}
		return Object(result)
	}
	return Error("unknown native type")
}
