package main

import (
	"context"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"log/slog"
)

func stateMonitor() {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		WithEventLog(fsmLog).
		WithStateLog(fsmLog).
		Build(context.Background())

	mon := yafasm.NewStateMonitor(door)

	slog.Info("locked state entered", slog.Time("at", mon.LastEnterAt(door.State())))

	_ = door.Event(context.Background(), Unlock)

	slog.Info("closed state entered", slog.Time("at", mon.LastEnterAt(door.State())))

	_ = door.Event(context.Background(), Lock)

	slog.Info("locked state entered", slog.Time("at", mon.LastEnterAt(door.State())))
}
