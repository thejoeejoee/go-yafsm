package yafasm

import "context"

type Notification func(ctx context.Context)

func (b *Builder[S, E]) AddNotification(e E, notification Notification) *Builder[S, E] {
	b.callbacks.events[e][afterEvent] = append(
		b.callbacks.events[e][afterEvent],
		func(ctx context.Context, event E) {
			notification(ctx)
		},
	)
	return b
}
