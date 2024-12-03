package yafasm

import "context"

type Builder[S, E comparable] struct {
	store Store[S, E]

	transitions Transitions[S, E]
	conditions  map[E][]Condition

	callbacks Callbacks[S, E]
}

func New[S, E comparable]() *Builder[S, E] {
	return &Builder[S, E]{
		conditions: map[E][]Condition{},
		callbacks: Callbacks[S, E]{
			events:       map[E]SingleEventCallbacks[E]{},
			states:       map[S]SingleStateCallbacks[S]{},
			allEvents:    map[eventCallbackType][]EventCallback[E]{},
			allStates:    map[stateCallbackType][]StateCallback[S]{},
			conditionErr: []EventCallback[E]{},
		},
	}
}

func (b *Builder[S, E]) WithInitial(initial S) *Builder[S, E] {
	b.store = &dummyStore[S, E]{initial: initial}
	return b
}

func (b *Builder[S, E]) WithTransitions(transitions Transitions[S, E]) *Builder[S, E] {
	b.transitions = transitions

	for _, to := range transitions {
		for event := range to {
			if _, ok := b.conditions[event]; !ok {
				b.conditions[event] = []Condition{}
			}
		}
	}

	// fill callbacks.events with empty maps
	for _, to := range transitions {
		for event := range to {
			if _, ok := b.callbacks.events[event]; !ok {
				b.callbacks.events[event] = map[eventCallbackType][]EventCallback[E]{}
			}
		}
	}
	// fill callbacks.states with empty maps
	for _, to := range transitions {
		for _, from := range to {
			if _, ok := b.callbacks.states[from]; !ok {
				b.callbacks.states[from] = map[stateCallbackType][]StateCallback[S]{}
			}
		}
	}

	return b
}

func (b *Builder[S, E]) WithStore(store Store[S, E]) *Builder[S, E] {
	b.store = store
	return b
}

func (b *Builder[S, E]) Build(ctx context.Context) (*Machine[S, E], error) {
	initial, err := b.store.Load(ctx)

	if err != nil {
		return nil, err
	}

	m := &Machine[S, E]{
		state:       initial,
		transitions: b.transitions,
		conditions:  b.conditions,
		store:       b.store,
		callbacks:   b.callbacks,
	}

	// run enter state for initial state
	m.fireStateCallbacks(ctx, enterState, initial)

	return m, nil
}
