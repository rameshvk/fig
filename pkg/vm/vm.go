// Package vm implements a stack machine
//
// A Machine has the following sections:
//
// - Code has the byte code.
// - Stack has the main computation stack.
// - Locals has the current set of locals (and also the return
//   address).
// - Constants has a set of arbitrary constant values. This is
//   where the runtime support functions are stored. This also
//   doubles up as the function table.
//
// There are six instructions:
// - PushLocal pushes a local at the specified index onto the stack
// - PopLocal pops the stack into thee local at the specified index
// - PushConnst pushes a const at the specified index onto the stack
// - Func call makes a function call to the address at the top of the
//   stack using the specified number of stack entries as args. This
//   pops the args from the stack into the locals. The old locals are
//   saved
// - Return leaves the stack alone but unwinds the old locals and also
//   jumps to the return address
// - DeleteLocal is an optional instruction that clears the local at
//   the specified index. This causes that local to be returned for
//   garbage collection
//
// Note that the VM does not have a memory model. That is completely
// up to the runtime methods.  The constants collection typically has
// runtime methods which could allocate/free memory as needed.
//
// Encoding
//
// The instructions are encoded as follows:
//
// - If the index/number needed as arg for each instruction is
//   less than 8, it is just stored in the lower 3 bits of the
//   instruction. If it is greater, then extra bytes (1-4) are
//   are used to store it in little-endian format.
//
// - The six instructions are numbered 0 - 5.
// - The nubmer of extra bytes can be 0 - 4.
// - The two are combined (first * 5 + extra_bytes_count) to form
//   the higher 5-bits of the first byte.
//
// - Note that the lower 3-bits are unused if there are any
//   extra bytes. Also, Return alwyas has zero extra bytes
package vm

// State represents the state of a virtual machine.
//
// All Machines start in Init state. They move to Running
// state with the first instruction.  They move to Stopped
// state once all the instructions are exhausted.
//
// The Paused state indicates debugger interactions. The Suspend state
// indicates the VM has been suspended.
type State int

const (
	Init State = iota
	Running
	Stopped
	Paused
	Suspended
)

// Machine wraps Code, Stack, Locals to form a virtual machine
type Machine struct {
	Code
	Stack
	Locals
	Constants

	State
	Address uint32

	defers []func(m *Machine)

	Stats
}

// Run executes one instruction or less (limited by the provided
// cost).  It returns the actual cost of the execution
func (m *Machine) Run(maxCost int) int {
	if m.State != Paused && m.State != Init {
		return 0
	}

	m.SetState(Running)
	cost := m.processDefers(0, maxCost)
	if instr := m.DecodeInstruction(m.Address); instr == nil {
		m.SetState(Stopped)
	} else if cost < maxCost && m.State == Running {
		cost++
		m.Stats.Instructions++
		m.Address += instr.Size()
		instr.Exec(m)
	}
	cost = m.processDefers(cost, maxCost)
	if m.State == Running {
		m.SetState(Paused)
	}
	return cost
}

func (m *Machine) processDefers(cost, maxCost int) int {
	for cost < maxCost && len(m.defers) > 0 && m.State == Running {
		cost++
		m.Stats.Defers++
		m.defers[0](m)
		m.defers = m.defers[1:]
	}
	return cost
}

// Defer adds a task to the Defers queue
func (m *Machine) Defer(fn func(m *Machine)) {
	m.defers = append(m.defers, fn)
}

// SetState sets the current state of the VM
func (m *Machine) SetState(newState State) {
	m.Stats.recordTime(m.State, newState)
	m.State = newState
}

// PushLocal returns an instruction which would push the value stored
// in the local at the specified index onto the top of the stack
func PushLocal(idx uint32) Instruction {
	return pushLocal(idx + 1)
}

// PopLocal returns an instruction which would pop the data stack
// and store it in the local at the specified index
func PopLocal(idx uint32) Instruction {
	return popLocal(idx + 1)
}

// DeleteLocal returns an instruction which deletes the local at the
// specified index.
func DeleteLocal(idx uint32) Instruction {
	return deleteLocal(idx + 1)
}

// PushConst returns an instruction which pushes the const at the
// specified index onto the top of the stack
func PushConst(idx uint32) Instruction {
	return pushConst(idx)
}

// FuncCall returns an instruction which executes a function call of
// the specified number of args
func FuncCall(nArgs uint32) Instruction {
	return callFunc(nArgs)
}

// RetCall returns an instruction which "returns"
func FuncRet() Instruction {
	return retFunc(0)
}

// CallableAddress converts an address to a callable.  This is
// requried for executing a function call at a specific address as
// function calls expect a callable on stack.
func CallableAddress(address uint32) Callable {
	return callableAddress(address)
}
