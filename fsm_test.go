package yafasm_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/thejoeejoee/go-yafsm"
	"testing"
)

type DoorState string

const (
	Locked DoorState = "Locked"
	Closed DoorState = "Closed"
	Opened DoorState = "Opened"
)

type DoorEvent string

const (
	Lock   DoorEvent = "Lock"
	Unlock DoorEvent = "Unlock"
	Open   DoorEvent = "Open"
	Close  DoorEvent = "Close"
)

var transitions = yafasm.Transitions[DoorState, DoorEvent]{
	Locked: {Unlock: Closed},
	Closed: {Lock: Locked, Open: Opened},
	Opened: {Close: Locked},
}

func TestMachine_Event(t *testing.T) {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		Build(context.Background())

	ctx := context.Background()

	assert.Equalf(t, Locked, door.State(), "door should be locked")

	assert.NoError(t, door.Event(ctx, Unlock))
	assert.Equalf(t, Closed, door.State(), "door should be closed")

	assert.NoError(t, door.Event(ctx, Open))
	assert.Equalf(t, Opened, door.State(), "door should be opened")

	assert.NoError(t, door.Event(ctx, Close))
	assert.Equalf(t, Locked, door.State(), "door should be locked")

	assert.NoError(t, door.Event(ctx, Lock))
	assert.Equalf(t, Locked, door.State(), "door should be locked")
}

func TestMachine_Can(t *testing.T) {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		Build(context.Background())

	assert.True(t, door.Can(Unlock), "should be able to unlock")
	assert.False(t, door.Can(Open), "should not be able to open")
	assert.False(t, door.Can(Close), "should not be able to close")
	assert.False(t, door.Can(Lock), "should not be able to lock")

	assert.NoError(t, door.Event(context.Background(), Unlock))

	assert.True(t, door.Can(Open), "should be able to open")
	assert.True(t, door.Can(Lock), "should be able to lock")
	assert.False(t, door.Can(Unlock), "should not be able to unlock")
	assert.False(t, door.Can(Close), "should not be able to close")
}

func TestMachine_AddNotification(t *testing.T) {
	called := false

	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		EventProcessed(Unlock, func(ctx context.Context, event DoorEvent) {
			assert.Equal(t, Unlock, event)
			assert.Equal(t, Locked, *yafasm.OriginStateFromCtx[DoorState](ctx))
			assert.Equal(t, Closed, *yafasm.TargetStateFromCtx[DoorState](ctx))
			called = true
		}).
		Build(context.Background())

	ctx := context.Background()
	assert.Equalf(t, Locked, door.State(), "door should be locked")
	assert.NoError(t, door.Event(ctx, Unlock))

	assert.True(t, called, "notification should be called")
}

func TestMachine_WithNotifications(t *testing.T) {
	called1 := false
	called2 := false

	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		AddNotification(Unlock, func(ctx context.Context) {
			assert.Equal(t, Locked, *yafasm.OriginStateFromCtx[DoorState](ctx))
			assert.Equal(t, Closed, *yafasm.TargetStateFromCtx[DoorState](ctx))
			assert.Equalf(t, Unlock, *yafasm.EventFromCtx[DoorEvent](ctx), "event should be unlock")
			called1 = true
		}).
		AddNotification(Unlock, func(ctx context.Context) {
			assert.Equal(t, Locked, *yafasm.OriginStateFromCtx[DoorState](ctx))
			assert.Equal(t, Closed, *yafasm.TargetStateFromCtx[DoorState](ctx))
			assert.Equalf(t, Unlock, *yafasm.EventFromCtx[DoorEvent](ctx), "event should be unlock")
			called2 = true
		}).
		Build(context.Background())

	ctx := context.Background()
	assert.Equalf(t, Locked, door.State(), "door should be locked")
	assert.NoError(t, door.Event(ctx, Unlock))

	assert.True(t, called1, "notification 1 should be called")
	assert.True(t, called2, "notification 2 should be called")
}

func TestMachine_WithCondition(t *testing.T) {
	var pinKey struct{}
	var invalidPin = errors.New("invalid pin")
	const LockPIN = "1234"

	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		AddCondition(Unlock, func(ctx context.Context) error {
			assert.Equal(t, Locked, *yafasm.OriginStateFromCtx[DoorState](ctx))
			assert.Equal(t, Closed, *yafasm.TargetStateFromCtx[DoorState](ctx))
			assert.Equalf(t, Unlock, *yafasm.EventFromCtx[DoorEvent](ctx), "event should be unlock")

			if p, ok := ctx.Value(pinKey).(string); ok && p == LockPIN {
				return nil
			}

			return invalidPin
		}).
		Build(context.Background())

	ctx := context.Background()

	assert.Equalf(t, Locked, door.State(), "door should be locked")
	err := door.Event(ctx, Unlock)
	assert.ErrorIs(t, err, yafasm.ErrConditionFailed, "condition failed")
	assert.ErrorIs(t, err, invalidPin, "condition failed")
	assert.Equalf(t, Locked, door.State(), "door should be locked")

	ctx = context.WithValue(ctx, pinKey, LockPIN)
	err = door.Event(ctx, Unlock)
	assert.NoError(t, err, "should be able to unlock")
	assert.Equalf(t, Closed, door.State(), "door should be closed")
}
