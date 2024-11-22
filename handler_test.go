package yafasm_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"testing"
)

func TestMachine_Handler(t *testing.T) {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		Build(context.Background())

	type msg string
	const Message1 msg = "message1"
	const Message2 msg = "message2"

	const OnLocked = "onLocked"
	const OnClosed = "onClosed"

	counts := map[string]map[msg]int{OnLocked: {}, OnClosed: {}}

	handler := func(s string) func(ctx context.Context, m msg) error {
		return func(ctx context.Context, m msg) error {
			counts[s][m]++
			return nil
		}
	}

	h := yafasm.NewHandler[msg](door)
	h.On(Locked, handler(OnLocked))
	h.On(Closed, handler(OnClosed))

	ctx := context.Background()
	// initially locked, so Message1 should be called
	assert.NoError(t, h.Handle(ctx, Message1))
	assert.NoError(t, h.Handle(ctx, Message2))
	assert.Equal(t, 1, counts[OnLocked][Message1])
	assert.Equal(t, 1, counts[OnLocked][Message2])
	assert.Equal(t, 0, counts[OnClosed][Message1])
	assert.Equal(t, 0, counts[OnClosed][Message2])

	// unlock the door
	assert.NoError(t, door.Event(ctx, Unlock))
	// now the door is closed, so onClosed should be called
	assert.NoError(t, h.Handle(ctx, Message1))
	assert.NoError(t, h.Handle(ctx, Message2))
	assert.Equal(t, 1, counts[OnLocked][Message1])
	assert.Equal(t, 1, counts[OnLocked][Message2])
	assert.Equal(t, 1, counts[OnClosed][Message1])
	assert.Equal(t, 1, counts[OnClosed][Message2])
}
