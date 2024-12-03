package yafasm_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"testing"
)

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

func TestMachine_WithConditionErrCallback(t *testing.T) {
	var pinKey struct{}
	var invalidPin = errors.New("invalid pin")
	const LockPIN = "1234"

	var catchErr error

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
		OnConditionErr(func(ctx context.Context, event DoorEvent) {
			catchErr = *yafasm.ConditionErrFromCtx(ctx)
		}).
		Build(context.Background())

	ctx := context.Background()

	assert.Equalf(t, Locked, door.State(), "door should be locked")
	err := door.Event(ctx, Unlock)
	assert.ErrorIs(t, err, yafasm.ErrConditionFailed, "condition failed")
	assert.ErrorIs(t, err, invalidPin, "condition failed")
	assert.Equalf(t, Locked, door.State(), "door should be locked")

	assert.ErrorIs(t, catchErr, invalidPin, "condition failed")
	catchErr = nil

	ctx = context.WithValue(ctx, pinKey, LockPIN)
	err = door.Event(ctx, Unlock)
	assert.NoError(t, err, "should be able to unlock")
	assert.Equalf(t, Closed, door.State(), "door should be closed")
	assert.Nil(t, catchErr, "condition not failed")
}
