package yafasm

import "context"

// Handler proxies f to the machine and calls the appropriate function based on the current state.
type Handler[V any, S, E comparable] struct {
	machine *Machine[S, E]

	f map[S]func(context.Context, V) error
}

// NewHandler creates a new Handler.
func NewHandler[V any, S, E comparable](m *Machine[S, E]) *Handler[V, S, E] {
	return &Handler[V, S, E]{
		machine: m,
		f:       make(map[S]func(context.Context, V) error),
	}
}

// On sets the function to be called when the machine is in state s.
func (i *Handler[V, S, E]) On(s S, f func(context.Context, V) error) *Handler[V, S, E] {
	i.f[s] = f
	return i
}

// Drop removes the function to be called when the machine is in state s.
func (i *Handler[V, S, E]) Drop(s S) *Handler[V, S, E] {
	delete(i.f, s)
	return i
}

// Handle calls the appropriate function based on the current state.
func (i *Handler[V, S, E]) Handle(ctx context.Context, v V) error {
	f, ok := i.f[i.machine.State()]
	if !ok {
		return nil
	}
	return f(ctx, v)
}
