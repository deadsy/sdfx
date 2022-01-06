//go:build plan9 || solaris || js
// +build plan9 solaris js

package dev

import (
	"errors"
)

// fsEvent represents a single file system notification.
type fsEvent struct {
	Name string // Relative path to the file or directory.
	Op   fsOp   // File operation that triggered the event.
}

// fsOp describes a set of file operations.
type fsOp uint32

var errFsNotifyNotSupported = errors.New("fsnotify is not supported for this platform")

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events chan fsEvent
	Errors chan error
}

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewFsWatcher() (*Watcher, error) {
	return nil, errFsNotifyNotSupported
}

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error {
	return nil
}

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error {
	return nil
}

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error {
	return nil
}
