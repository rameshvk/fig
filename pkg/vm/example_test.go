package vm_test

import (
	"fmt"

	"github.com/rameshvk/fig/pkg/vm"
)

func Example() {
	// constant locations
	c0 := uint32(7)
	c1 := uint32(8)
	c2 := uint32(256)
	c3 := uint32(256 * 266)
	c4 := uint32(256 * 256 * 256)
	addf := uint32(0)
	sumf := uint32(1)
	exitf := uint32(2)

	consts := make(vm.Constants, c4+1)
	consts[int(addf)] = vm.Function(func(args ...interface{}) interface{} {
		return args[0].(int) + args[1].(int)
	})
	consts[int(c0)] = 1
	consts[int(c1)] = 2
	consts[int(c2)] = 3
	consts[int(c3)] = 4
	consts[int(c4)] = 5

	code := vm.Code(nil).AppendInstructions(
		vm.PushConst(c0),
		vm.PushConst(c1),
		vm.PushConst(c2),
		vm.PushConst(c3),
		vm.PushConst(c4),
		vm.PushConst(sumf),
		vm.FuncCall(5), // 5 args
		vm.PushConst(exitf),
		vm.FuncCall(0), // exit!
	)

	// the actual function being called is here
	consts[int(sumf)] = vm.CallableAddress(uint32(len(code)))
	code = code.AppendInstructions(
		vm.PushLocal(0),
		vm.DeleteLocal(0),
		vm.PushLocal(1),
		vm.DeleteLocal(1),
		vm.PushConst(addf), // adder
		vm.FuncCall(2),     // local0 + local1

		// the following three are not really necessary
		// but it tests that the part of the code
		vm.PopLocal(0),
		vm.PushLocal(0),
		vm.DeleteLocal(0),

		vm.PushLocal(2),
		vm.DeleteLocal(2),
		vm.PushConst(addf), // adder
		vm.FuncCall(2),     // last + local2

		vm.PushLocal(3),
		vm.DeleteLocal(3),
		vm.PushConst(addf),
		vm.FuncCall(2), // last + local3

		vm.PushLocal(4),
		vm.DeleteLocal(4),
		vm.PushConst(addf),
		vm.FuncCall(2), // last + local4

		vm.FuncRet(),
	)
	consts[exitf] = vm.CallableAddress(uint32(len(code)))

	m := vm.Machine{Code: code, Constants: consts}
	cost := 0
	for m.State != vm.Stopped {
		cost += m.Run(1000)
	}

	fmt.Printf("Result: %v (cost = %d)\n", m.Pop(), cost)

	// Output:
	// Result: 15 (cost = 31)
}
