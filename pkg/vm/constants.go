package vm

// Constant tracks all program constant values
type Constants []interface{}

// Constant returns the value at the specified offset
func (c Constants) Constant(at uint32) interface{} {
	return c[at]
}
