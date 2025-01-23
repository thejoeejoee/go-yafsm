package yafasm

import (
	"context"
)

type eventCallbackType string

const (
	eventReceived  eventCallbackType = "eventReceived"
	eventProcessed eventCallbackType = "eventProcessed"
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

	conditionErr []EventCallback[E]
	anyCondition []Condition
}

func (b *Builder[S, E]) EventReceived(e E, callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.events[e][eventReceived] = append(b.callbacks.events[e][eventReceived], callback)
	return b
}
func (b *Builder[S, E]) AnyReceived(callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.allEvents[eventReceived] = append(b.callbacks.allEvents[eventReceived], callback)
	return b
}
func (b *Builder[S, E]) EventProcessed(e E, callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.events[e][eventProcessed] = append(b.callbacks.events[e][eventProcessed], callback)
	return b
}
func (b *Builder[S, E]) AnyProcessed(callback EventCallback[E]) *Builder[S, E] {
	b.callbacks.allEvents[eventProcessed] = append(b.callbacks.allEvents[eventProcessed], callback)
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
