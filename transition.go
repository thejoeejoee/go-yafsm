package yafasm

import (
	"context"
	"fmt"
)

// Transitions defines all possible transitions between states.
//
// The key is the current state, and the value is a map of events to the next state.
//
// S is the type of the state, and E is the type of the event.
type Transitions[S, E comparable] map[S]map[E]S

// The Event tries to transition to the next state based on the event.
//
// ErrConditionFailed is returned when a condition fails.
func (m *Machine[S, E]) Event(ctx context.Context, e E) error {
	ctx = withOriginState(ctx, m.state)
	ctx = withEvent(ctx, e)

	m.fireEventCallbacks(ctx, eventReceived, e)
	defer func() {
		m.fireEventCallbacks(ctx, eventProcessed, e)
	}()

	if _, ok := m.transitions[m.state][e]; ok {
		to := m.transitions[m.state][e]

		ctx = withTargetState(ctx, to)

		err := m.checkTransition(ctx, e, m.state, to)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrConditionFailed, err)
		}

		m.fireStateCallbacks(ctx, leaveState, m.state)
		defer func() {
			m.fireStateCallbacks(ctx, enterState, m.state)
		}()

		m.state = to
	}

	return nil
}

func (m *Machine[S, E]) Can(e E) bool {
	_, ok := m.transitions[m.state][e]
	return ok
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
		return fmt.Errorf("%w %w", ErrStoreSaveFailed, err)
	}

	return nil
}
