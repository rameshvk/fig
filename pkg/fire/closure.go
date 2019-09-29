package fire

import (
	"context"
)

func closure(ctx context.Context, args []interface{}, global Value) Value {
	return errorValue("closure not yet implemented")
}
