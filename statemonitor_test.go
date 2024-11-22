package yafasm_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	yafasm "github.com/thejoeejoee/go-yafsm"
	"testing"
	"time"
)

func TestStateMonitor(t *testing.T) {
	door, _ := yafasm.New[DoorState, DoorEvent]().
		WithInitial(Locked).
		WithTransitions(transitions).
		AnyReceived(func(ctx context.Context, event DoorEvent) {
			t.Logf("before %s", event)
		}).
		AnyProcessed(func(ctx context.Context, event DoorEvent) {
			t.Logf("after %s", event)
		}).
		Build(context.Background())

	stamp := &time.Time{}
	tick0 := time.Time{}

	tick := func(v int) time.Time {
		return tick0.Add(time.Duration(v) * time.Second)
	}

	now := func() time.Time {
		defer func() {
			ticked := stamp.Add(time.Second)
			stamp = &ticked
		}()
		t.Logf("now: %v", *stamp)
		return *stamp
	}

	m := yafasm.NewStateMonitor[DoorState, DoorEvent](
		door,
		yafasm.WithMonitorTime[DoorState, DoorEvent](now),
	)

	assert.Equal(t, Locked, door.State())
	assert.Equal(t, m.FirstEnterAt(Locked), tick(0))

	assert.NoError(t, door.Event(context.Background(), Unlock))
	assert.Equal(t, Closed, door.State())
	assert.Equal(t, m.FirstLeaveAt(Locked), tick(1))
	assert.Equal(t, m.LastLeaveAt(Locked), tick(1))
	// leave, tick, enter
	assert.Equal(t, m.FirstEnterAt(Closed), tick(2))
	assert.Equal(t, m.LastEnterAt(Closed), tick(2))

	assert.NoError(t, door.Event(context.Background(), Open)) // tick3
	assert.Equal(t, Opened, door.State())

	assert.Equal(t, m.FirstLeaveAt(Locked), tick(1))
	assert.Equal(t, m.LastLeaveAt(Locked), tick(1))
	assert.Equal(t, m.FirstLeaveAt(Closed), tick(3))
	assert.Equal(t, m.LastLeaveAt(Closed), tick(3))
	assert.Equal(t, m.FirstEnterAt(Opened), tick(4))
	assert.Equal(t, m.LastEnterAt(Opened), tick(4))

}
