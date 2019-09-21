package vm

// Locals implements the store of local values (scope)
type Locals struct {
	values []interface{}
	active []bool
	prev   *Locals
}

// Store stores a value at the specified offset.
func (l *Locals) Store(at uint32, v interface{}) {
	l.ensure(at)
	l.values[at] = v
	l.active[at] = true
}

// Fetch fetches a value at the specified offset.
//
// If the offset was not stored before, this will panic.
func (l *Locals) Fetch(at uint32) interface{} {
	if !l.active[at] {
		panic("invalid fetch")
	}
	return l.values[at]
}

// Delete inactivates the value at the offset.
func (l *Locals) Delete(at uint32) {
	if !l.active[at] {
		panic("invalid delete")
	}
	l.values[at] = nil
	l.active[at] = false
	for kk := at + 1; kk < uint32(len(l.active)); kk++ {
		if l.active[kk] {
			return
		}
	}
	kk := at
	for ; kk >= 0 && !l.active[kk]; kk-- {
	}
	l.active = l.active[:kk+1]
	l.values = l.values[:kk+1]
}

// PushLocals saves the current locals and creates a new set of locals
func (l *Locals) PushLocals() {
	prev := *l
	l.prev = &prev
	l.values = nil
	l.active = nil
}

// PopLocals restores the locals back
func (l *Locals) PopLocals() {
	prev := l.prev
	*l = *prev
}

func (l *Locals) ensure(idx uint32) {
	for uint32(len(l.values)) <= idx {
		l.values = append(l.values, nil)
		l.active = append(l.active, false)
	}
}
