package yafasm

import "context"

type Condition func(ctx context.Context) error

func (b *Builder[S, E]) AddCondition(e E, condition Condition) *Builder[S, E] {
	b.conditions[e] = append(b.conditions[e], condition)
	return b
}
