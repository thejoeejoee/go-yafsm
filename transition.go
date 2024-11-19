package yafasm

import (
	"context"
	"errors"
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

	m.fireEventCallbacks(ctx, beforeEvent, e)
	defer func() {
		m.fireEventCallbacks(ctx, afterEvent, e)
	}()

	if _, ok := m.transitions[m.state][e]; ok {
		to := m.transitions[m.state][e]

		ctx = withTargetState(ctx, to)

		err := m.checkTransition(ctx, e, m.state, to)
		if err != nil {
			return errors.Join(ErrConditionFailed, err)
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
