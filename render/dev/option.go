package dev

import (
	"github.com/cenkalti/backoff/v4"
	"os/exec"
	"time"
)

// Option configures a Renderer to statically change its default behaviour.
type Option = func(r *Renderer)

// NOTE: There are more options defined in impl*.go, all starting with Opt<X>

// OptMRunCommand replaces the default run command (go run -v .) with any other command generator.
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMRunCommand(runCmd func() *exec.Cmd) Option {
	return func(r *Renderer) {
		r.runCmd = runCmd
	}
}

// OptMWatchFiles replaces the default set of files to watch for changes (["."]).
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMWatchFiles(filePaths []string) Option {
	return func(r *Renderer) {
		r.watchFiles = filePaths
	}
}

// OptMBackoff changes the default backoff algorithm used when trying to connect to the new code.
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMBackoff(backOff backoff.BackOff) Option {
	return func(r *Renderer) {
		r.backOff = backOff
	}
}

// OptMPartialRenderEvery changes the default duration between partial renders (loading a partial render takes a little
// time and slows down the full render if too frequent).
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMPartialRenderEvery(duration time.Duration) Option {
	return func(r *Renderer) {
		r.partialRenderEvery = duration
	}
}

// OptMZoom changes the default scaling factor (> 1)
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMZoom(zoom float64) Option {
	return func(r *Renderer) {
		r.zoomFactor = zoom
	}
}

// OptMResInv changes the default image pixels per rendererd pixel
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMResInv(resInv int) Option {
	return func(r *Renderer) {
		r.implState.ResInv = resInv
	}
}

// OptMColorMode changes the default color mode of the renderer
// WARNING: Need to run again the main renderer to apply a change of this option.
func OptMColorMode(colorMode int) Option {
	return func(r *Renderer) {
		r.implState.ColorMode = colorMode % r.impl.ColorModes()
	}
}
