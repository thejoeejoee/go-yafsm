package yafasm

import (
	"context"
	"errors"
)

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

func (m *Machine[S, E]) checkTransition(ctx context.Context, event E, origin S, target S) error {
	for _, condition := range m.conditions[event] {
		err := condition(ctx)
		if err != nil {
			return err
		}
	}

	err := m.store.Save(ctx, target)
	if err != nil {
		return errors.Join(ErrStoreSaveFailed, err)
	}

	return nil
}
