package yafasm

import (
	"context"
)

type eventCallbackType string

const (
	beforeEvent eventCallbackType = "beforeEvent"
	afterEvent  eventCallbackType = "afterEvent"
)

type stateCallbackType string

const (
	enterState stateCallbackType = "enterState"
	leaveState stateCallbackType = "leaveState"
)

type EventCallback[E comparable] func(ctx context.Context, event E)
type StateCallback[S comparable] func(ctx context.Context, state S)

type SingleEventCallbacks[E comparable] map[eventCallbackType][]EventCallback[E]
type SingleStateCallbacks[S comparable] map[stateCallbackType][]StateCallback[S]

type EventCallbacks[E comparable] map[E]SingleEventCallbacks[E]
type StateCallbacks[S comparable] map[S]SingleStateCallbacks[S]

type Callbacks[S, E comparable] struct {
	events EventCallbacks[E]
	states StateCallbacks[S]

	allEvents SingleEventCallbacks[E]
	allStates SingleStateCallbacks[S]
}

func (b *Builder[S, E]) BeforeEvent(e E, callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.events[e][beforeEvent] = append(b.callbacks.events[e][beforeEvent], callback)
	return b
}
func (b *Builder[S, E]) BeforeAll(callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.allEvents[beforeEvent] = append(b.callbacks.allEvents[beforeEvent], callback)
	return b
}
func (b *Builder[S, E]) AfterEvent(e E, callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.events[e][afterEvent] = append(b.callbacks.events[e][afterEvent], callback)
	return b
}
func (b *Builder[S, E]) AfterAll(callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.allEvents[afterEvent] = append(b.callbacks.allEvents[afterEvent], callback)
	return b
}

func (b *Builder[S, E]) EnterState(s S, callback StateCallback[S]) *Builder[S, E] {
	b.callbacks.states[s][enterState] = append(b.callbacks.states[s][enterState], callback)
	return b
}

func (b *Builder[S, E]) EnterAll(callback StateCallback[S]) *Builder[S, E] {
	b.callbacks.allStates[enterState] = append(b.callbacks.allStates[enterState], callback)
	return b
}

func (b *Builder[S, E]) LeaveState(s S, callback StateCallback[S]) *Builder[S, E] {
	b.callbacks.states[s][leaveState] = append(b.callbacks.states[s][leaveState], callback)
	return b
}

func (b *Builder[S, E]) LeaveAll(callback StateCallback[S]) *Builder[S, E] {
	b.callbacks.allStates[leaveState] = append(b.callbacks.allStates[leaveState], callback)
	return b
}

func (m *Machine[S, E]) fireStateCallbacks(ctx context.Context, t stateCallbackType, state S) {
	for _, callback := range m.callbacks.allStates[t] {
		callback(ctx, state)
	}

	for _, callback := range m.callbacks.states[state][t] {
		callback(ctx, state)
	}
}
func (m *Machine[S, E]) fireEventCallbacks(ctx context.Context, t eventCallbackType, event E) {
	//slog.Info("fireEventCallbacks", "event", event, "state", m.state, "t", t)

	for _, callback := range m.callbacks.allEvents[t] {
		callback(ctx, event)
	}

	for _, callback := range m.callbacks.events[event][t] {
		callback(ctx, event)
	}
}
