package main

import (
	yafasm "github.com/thejoeejoee/go-yafsm"
	"log/slog"
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

var fsmLog = slog.Default().WithGroup("fsm")
