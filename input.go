package yafasm

import "context"

// Handler proxies f to the machine and calls the appropriate function based on the current state.
type Handler[S, E comparable, V any] struct {
	machine *Machine[S, E]

	f map[S]func(context.Context, V) error
}

// NewHandler creates a new Handler.
func NewHandler[S, E comparable, V any](m *Machine[S, E]) *Handler[S, E, V] {
	return &Handler[S, E, V]{
		machine: m,
		f:       make(map[S]func(context.Context, V) error),
	}
}

// On adds a new input to the machine.
func (i *Handler[S, E, V]) On(s S, f func(context.Context, V) error) *Handler[S, E, V] {
	i.f[s] = f
	return i
}

// Drop removes an input from the machine.
func (i *Handler[S, E, V]) Drop(s S) *Handler[S, E, V] {
	delete(i.f, s)
	return i
}

// Handle calls the appropriate function based on the current state.
func (i *Handler[S, E, V]) Handle(ctx context.Context, v V) error {
	f, ok := i.f[i.machine.State()]
	if !ok {
		return nil
	}
	return f(ctx, v)
}
