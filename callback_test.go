package yafasm_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	yafsm "github.com/thejoeejoee/go-yafsm"
	"testing"
)

func TestEnterStateCallbackForInitial(t *testing.T) {
	entered := false

	_, _ = yafsm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		EnterState(Locked, func(ctx context.Context, state DoorState) {
			entered = true
		}).
		Build(context.Background())

	assert.True(t, entered)

}
