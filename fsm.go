package yafasm

// Machine is a finite state machine.
// S is the type of the states.
// E is the type of the events.
type Machine[S, E comparable] struct {
	state       S
	transitions Transitions[S, E]

	conditions map[E][]Condition

	store Store[S, E]

	callbacks Callbacks[S, E]
}

// State returns the current state of the machine.
func (m *Machine[S, E]) State() S {
	return m.state
}
