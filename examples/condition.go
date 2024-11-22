package main

import (
	"context"
	"errors"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"log/slog"
)

type pinCtxKey struct{}

const doorPin = "4321"

func checkPin(ctx context.Context) error {
	// get pin from context
	pin, ok := ctx.Value(pinCtxKey{}).(string)
	if !ok {
		return errors.New("pin not found in context")
	}

	// check pin
	if pin != doorPin {
		return errors.New("invalid pin")
	}

	// pin is correct
	return nil
}

func condition() {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		WithEventLog(fsmLog).
		WithStateLog(fsmLog).
		AddCondition(Unlock, checkPin).
		Build(context.Background())

	noCtx := context.Background()
	wrongCtx := context.WithValue(noCtx, pinCtxKey{}, "1234")
	correctCtx := context.WithValue(noCtx, pinCtxKey{}, "4321")

	slog.Info("initial state", slog.String("state", string(door.State())))

	err := door.Event(noCtx, Unlock)

	slog.Info("unlock error", slog.Any("error", err))

	err = door.Event(wrongCtx, Unlock)

	slog.Info("unlock error", slog.Any("error", err))

	err = door.Event(correctCtx, Unlock)

	slog.Info("final state", slog.String("state", string(door.State())))

}
