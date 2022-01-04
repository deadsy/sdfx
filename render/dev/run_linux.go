package dev

import (
	"os"
	"syscall"
)

func signals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}
