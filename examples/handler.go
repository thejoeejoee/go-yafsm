package main

import (
	"context"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"log/slog"
)

func handler() {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		WithEventLog(fsmLog).
		WithStateLog(fsmLog).
		Build(context.Background())

	type message string

	handle := yafasm.NewHandler[message](door)

	handle.On(Locked, func(ctx context.Context, m message) error {
		slog.Info("message when locked", slog.String("message", string(m)))
		return nil
	})
	handle.On(Closed, func(ctx context.Context, m message) error {
		slog.Info("message when closed", slog.String("message", string(m)))
		return nil
	})

	ctx := context.Background()
	_ = handle.Handle(ctx, "before unlock 1")
	_ = handle.Handle(ctx, "before unlock 2")

	_ = door.Event(ctx, Unlock)

	_ = handle.Handle(ctx, "after unlock 1")
	_ = handle.Handle(ctx, "after unlock 2")
}
