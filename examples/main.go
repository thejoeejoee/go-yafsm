package main

import (
	"fmt"
	"log/slog"
)

func main() {

	cases := []struct {
		name string
		f    func()
	}{
		{"basic", basic},
		{"condition", condition},
		{"handler", handler},
		{"state_monitor", stateMonitor},
	}

	for _, c := range cases {
		slog.Info(fmt.Sprintf("==== %s ====", c.name))
		c.f()
		slog.Info(fmt.Sprintf("==== /%s ====", c.name))
	}
}
