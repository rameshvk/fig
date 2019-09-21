package vm

// Code manages the code section
type Code []byte

// DecodeInstruction decodes the instruction at the specified offset
func (c Code) DecodeInstruction(at uint32) Instruction {
	if at >= uint32(len(c)) {
		return nil
	}

	val := uint32(c[at] % 8)
	top := c[at] >> 3
	offset, extra := top/5, top%5
	if extra > 0 {
		val = 0
	}
	for kk := uint32(0); kk < uint32(extra); kk++ {
		val += uint32(c[at+kk+1]) << (kk * 8)
	}
	return instruction(val).decode(offset)
}

// AppendInstructions appends a bunch of instructions
func (c Code) AppendInstructions(instr ...Instruction) Code {
	c = append([]byte(nil), c...)
	buf := make([]byte, 5)
	for _, i := range instr {
		i.Encode(buf)
		c = append(c, buf[:i.Size()]...)
	}
	return c
}

// Instruction is the interface to execute instructions
type Instruction interface {
	Exec(m *Machine)
	Size() uint32
	Encode(buf []byte)
}

type instruction uint32

func (i instruction) size() uint32 {
	if i < 8 {
		return 1
	}
	if i < 256 {
		return 2
	}
	if i < 256*256 {
		return 3
	}
	if i < 256*256*256 {
		return 4
	}
	return 5
}
func (i instruction) encode(offset uint32, buf []byte) {
	extra := i.size() - 1
	if extra == 0 {
		buf[0] = byte(offset*5<<3 + uint32(i))
		return
	}
	buf[0] = byte(offset*5+extra) << 3
	val := uint32(i)
	for kk := uint32(0); kk < extra; kk++ {
		buf[kk+1] = byte(val % 256)
		val = val >> 8
	}
}

func (i instruction) decode(offset byte) Instruction {
	switch offset {
	case 0:
		return pushLocal(i)
	case 1:
		return popLocal(i)
	case 2:
		return deleteLocal(i)
	case 3:
		return pushConst(i)
	case 4:
		return callFunc(i)
	}
	return retFunc(i)
}

type pushLocal instruction

func (instr pushLocal) Exec(m *Machine) {
	m.Push(m.Fetch(uint32(instr)))
}
func (instr pushLocal) Encode(buf []byte) {
	instruction(instr).encode(0, buf)
}
func (instr pushLocal) Size() uint32 {
	return instruction(instr).size()
}

type popLocal instruction

func (instr popLocal) Exec(m *Machine) {
	m.Store(uint32(instr), m.Pop())
}
func (instr popLocal) Encode(buf []byte) {
	instruction(instr).encode(1, buf)
}
func (instr popLocal) Size() uint32 {
	return instruction(instr).size()
}

type deleteLocal instruction

func (instr deleteLocal) Exec(m *Machine) {
	m.Delete(uint32(instr))
}
func (instr deleteLocal) Encode(buf []byte) {
	instruction(instr).encode(2, buf)
}
func (instr deleteLocal) Size() uint32 {
	return instruction(instr).size()
}

type pushConst instruction

func (instr pushConst) Exec(m *Machine) {
	m.Push(m.Constant(uint32(instr)))
}
func (instr pushConst) Encode(buf []byte) {
	instruction(instr).encode(3, buf)
}
func (instr pushConst) Size() uint32 {
	return instruction(instr).size()
}

type callFunc instruction

func (instr callFunc) Exec(m *Machine) {
	m.Pop().(Callable).Call(m, m.Address, uint32(instr))
}
func (instr callFunc) Encode(buf []byte) {
	instruction(instr).encode(4, buf)
}
func (instr callFunc) Size() uint32 {
	return instruction(instr).size()
}

type retFunc instruction

func (instr retFunc) Exec(m *Machine) {
	m.Address = m.Fetch(0).(uint32)
	m.PopLocals()
}
func (instr retFunc) Encode(buf []byte) {
	instruction(instr).encode(5, buf)
}
func (instr retFunc) Size() uint32 {
	return instruction(instr).size()
}

type callableAddress uint32

func (c callableAddress) Call(m *Machine, returnAddress, nargs uint32) {
	m.PushLocals()
	m.Store(0, returnAddress)
	for kk := uint32(1); kk <= nargs; kk++ {
		m.Store(kk, m.Pop())
	}
	m.Address = uint32(c)
}

// Callable wraps the Call method.
type Callable interface {
	Call(m *Machine, returnAddress, nArgs uint32)
}

// Function is a simple function wrapper
type Function func(args ...interface{}) interface{}

func (f Function) Call(m *Machine, returnAddress, nArgs uint32) {
	args := make([]interface{}, nArgs)
	for kk := 0; kk < int(nArgs); kk++ {
		args[kk] = m.Pop()
	}
	m.Push(f(args...))
}
