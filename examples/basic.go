package main

import (
	"context"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"log/slog"
)

func basic() {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		WithEventLog(fsmLog).
		WithStateLog(fsmLog).
		Build(context.Background())

	ctx := context.Background()

	slog.Info("initial state", slog.String("state", string(door.State())))

	_ = door.Event(ctx, Unlock)

	slog.Info("final state", slog.String("state", string(door.State())))
}
