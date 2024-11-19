# Yet Another Golang Finite State Machine

* finite state machine
* Go Generics support
* conditions support
* callbacks support

## Usage

```go
package main

import (
    "context"
    "github.com/thejoeejoee/go-yafsm"
)

type DoorState string
type DoorEvent string

const (
	Locked DoorState = "Locked"
	Closed DoorState = "Closed"
	Opened DoorState = "Opened"
)

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

func main() {
    door, _ := yafasm.New[DoorState, DoorEvent]().
    WithInitial(Locked).
    WithTransitions(transitions).
    Build(context.Background())
    
    door.Event(Unlock)
}
```

## Author 

https://github.com/thejoeejoee