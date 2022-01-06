//go:build freebsd || openbsd || netbsd || dragonfly || darwin || windows || linux || solaris
// +build freebsd openbsd netbsd dragonfly darwin windows linux solaris

package dev

import "github.com/fsnotify/fsnotify"

func NewFsWatcher() (*fsnotify.Watcher, error) {
	return fsnotify.NewWatcher()
}
