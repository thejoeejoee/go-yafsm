package yafasm

import (
	"context"
	"time"
)

// StateMonitor records the time of the first and last enter/leave of each state.
// The time is recorded for leave and enter separately.
type StateMonitor[S, E comparable] struct {
	machine *Machine[S, E]

	first map[stateCallbackType]map[S]time.Time
	last  map[stateCallbackType]map[S]time.Time

	now func() time.Time
}

type StateMonitorOption[S, E comparable] func(*StateMonitor[S, E])

func WithMonitorTime[S, E comparable](now func() time.Time) StateMonitorOption[S, E] {
	return func(m *StateMonitor[S, E]) {
		m.now = now
	}
}

func NewStateMonitor[S, E comparable](
	m *Machine[S, E],
	opts ...StateMonitorOption[S, E],
) *StateMonitor[S, E] {
	mon := &StateMonitor[S, E]{
		machine: m,

		first: map[stateCallbackType]map[S]time.Time{
			enterState: {},
			leaveState: {},
		},
		last: map[stateCallbackType]map[S]time.Time{
			enterState: {},
			leaveState: {},
		},
		now: time.Now,
	}

	m.callbacks.allStates[enterState] = append(
		m.callbacks.allStates[enterState],
		mon.callback(enterState),
	)
	m.callbacks.allStates[leaveState] = append(
		m.callbacks.allStates[leaveState],
		mon.callback(leaveState),
	)

	for _, opt := range opts {
		opt(mon)
	}

	// record the initial state
	mon.callback(enterState)(context.Background(), m.State())

	return mon
}

func (m *StateMonitor[S, E]) callback(when stateCallbackType) func(_ context.Context, state S) {
	return func(_ context.Context, state S) {
		t := m.now()

		if _, ok := m.first[when][state]; !ok {
			m.first[when][state] = t
		}
		m.last[when][state] = t
	}
}

// FirstEnterAt returns the time when the state was first entered.
func (m *StateMonitor[S, E]) FirstEnterAt(state S) time.Time {
	return m.first[enterState][state]
}

// LastEnterAt returns the time when the state was entered the last time.
func (m *StateMonitor[S, E]) LastEnterAt(state S) time.Time {
	return m.last[enterState][state]
}

// FirstLeaveAt returns the time when the state was first left.
func (m *StateMonitor[S, E]) FirstLeaveAt(state S) time.Time {
	return m.first[leaveState][state]
}

// LastLeaveAt returns the time when the state was left the last time.
func (m *StateMonitor[S, E]) LastLeaveAt(state S) time.Time {
	return m.last[leaveState][state]
}
