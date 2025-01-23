package yafasm

import "context"

type Condition func(ctx context.Context) error

func (b *Builder[S, E]) AnyCondition(condition Condition) *Builder[S, E] {
	b.callbacks.anyCondition = append(b.callbacks.anyCondition, condition)
	return b
}

func (b *Builder[S, E]) AddCondition(e E, condition Condition) *Builder[S, E] {
	b.conditions[e] = append(b.conditions[e], condition)
	return b
}

func (b *Builder[S, E]) OnConditionErr(f EventCallback[E]) *Builder[S, E] {
	b.callbacks.conditionErr = append(b.callbacks.conditionErr, f)

	return b
}
