package yafasm

import "context"

// Store is responsible for persisting the state of the machine.
type Store[S, E comparable] interface {
	Load(context.Context) (S, error)
	Save(context.Context, S) error
}

var _ Store[any, any] = &dummyStore[any, any]{}

type dummyStore[S, E any] struct {
	initial S
}

func (p *dummyStore[S, E]) Load(ctx context.Context) (S, error) {
	return p.initial, nil
}

func (p *dummyStore[S, E]) Save(context.Context, S) error {
	return nil
}
