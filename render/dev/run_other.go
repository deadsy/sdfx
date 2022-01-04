//go:build !linux
// +build !linux

package dev

import (
	"os"
)

func signals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
