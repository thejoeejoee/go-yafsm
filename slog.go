package yafasm

import (
	"context"
	"log/slog"
)

func (b *Builder[S, E]) WithEventLog(logger *slog.Logger) *Builder[S, E] {
	b.callbacks.allEvents[eventReceived] = append(b.callbacks.allEvents[eventReceived], func(ctx context.Context, e E) {
		logger.Info("event received", slog.Any("event", e))
	})

	b.callbacks.allEvents[eventProcessed] = append(b.callbacks.allEvents[eventProcessed], func(ctx context.Context, e E) {
		logger.Info("event processed", slog.Any("event", e))
	})

	return b
}

func (b *Builder[S, E]) WithStateLog(logger *slog.Logger) *Builder[S, E] {
	b.callbacks.allStates[enterState] = append(b.callbacks.allStates[enterState], func(ctx context.Context, s S) {
		logger.Info("enter state", slog.Any("state", s))
	})

	b.callbacks.allStates[leaveState] = append(b.callbacks.allStates[leaveState], func(ctx context.Context, s S) {
		logger.Info("leave state", slog.Any("state", s))
	})

	return b
}
