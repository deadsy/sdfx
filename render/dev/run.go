package dev

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/fsnotify/fsnotify"
	"github.com/hajimehoshi/ebiten"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const changeEventThrottle = 100 * time.Millisecond

func (r *Renderer) runRenderer(runCmdF func() *exec.Cmd, watchGlobs []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Println("[DevRenderer] File watcher close error:", err)
		}
	}(watcher)

	go func() {
		var runCmd *exec.Cmd
		lastEvent := time.Now()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if time.Since(lastEvent) < changeEventThrottle {
					log.Println("[DevRenderer] Change detected (but throttled)!", event)
					continue // Events tend to be generated in bulk if using an IDE, skip them if too close together
				}
				log.Println("[DevRenderer] Change detected!", event)
				lastEvent = time.Now()
				runCmd = r.rendererSwapChild(runCmd, runCmdF)
				lastEvent = time.Now()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("[DevRenderer] File watcher error:", err)
			}
		}
	}()

	for _, watchGlob := range watchGlobs {
		glob, err := filepath.Glob(watchGlob)
		if err != nil {
			return err
		}
		for _, matchedFile := range glob {
			err = watcher.Add(matchedFile)
			if err != nil {
				return err
			}
		}
	}

	return ebiten.RunGame(r) // blocks until the window is closed
}

func (r *Renderer) runChild(requestedAddress string) error {
	// Set up a remote service that the parent renderer will connect to view the new SDF
	newDevRendererService(r.impl).HandleHTTP(rpc.DefaultRPCPath, "/debug")
	listener, err := net.Listen("tcp", requestedAddress)
	if err != nil {
		return err
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Println("[DevRenderer] listener.Close error:", err)
		}
	}(listener)
	srv := &http.Server{Addr: listener.Addr().String(), Handler: http.DefaultServeMux}
	defer func(srv *http.Server) {
		err := srv.Close()
		if err != nil {
			log.Println("[DevRenderer] srv.Close error:", err)
		}
	}(srv)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			log.Println("[DevRenderer] srv.Serve error:", err)
		}
		done <- syscall.SIGILL
	}()
	<-done // Will block until interrupt is received or the server crashes
	log.Println("[DevRenderer] Closing child service...")
	return nil
}

func (r *Renderer) rendererSwapChild(runCmd *exec.Cmd, runCmdF func() *exec.Cmd) *exec.Cmd {
	r.implLock.Lock() // No more renders until we swapped the implementation
	defer r.implLock.Unlock()
	// 1. Gracefully close the previous command
	if runCmd != nil {
		log.Println("[DevRenderer] Closing previous child process")
		err := pidTermWaitKill(runCmd.Process, 12*time.Second)
		if err != nil {
			log.Println("[DevRenderer] pidTermWaitKill error:", err)
			return nil
		}
	}
	log.Println("[DevRenderer] Compiling and running new code")
	// 2. Get a random free port to ask the child to listen on (it might not be free when the process starts, but ¯\_(ツ)_/¯)
	tmpL, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("[DevRenderer] net.Listen error:", err)
		return nil
	}
	requestedFreeAddr := tmpL.Addr().String()
	err = tmpL.Close()
	if err != nil {
		log.Println("[DevRenderer] tmpL.Close error:", err)
		return nil
	}
	// 3. Configure the process and start it in the background
	runCmd = runCmdF()
	runCmd.Env = append(os.Environ(), RequestedAddressEnvKey+"="+requestedFreeAddr)
	runCmd.Stdout = os.Stdout // Merge stdout
	runCmd.Stderr = os.Stderr // Merge stderr
	err = runCmd.Start()
	if err != nil {
		log.Println("[DevRenderer] runCmd.Start error:", err)
		return nil
	}
	// 4. Connect to it as fast as possible, with exponential backoff to relax on errors.
	log.Println("[DevRenderer] Trying to connect to new code with exponential backoff...")
	backOff := backoff.NewExponentialBackOff()
	startedAt := time.Now()
	backOff.InitialInterval = 100 * time.Millisecond
	err = backoff.RetryNotify(func() error {
		dialHTTP, err := rpc.DialHTTP("tcp", requestedFreeAddr)
		if err != nil {
			if time.Since(startedAt) > time.Second { // Make sure the process fully started
				exists, err2 := pidExists(int32(runCmd.Process.Pid))
				if !exists || err2 != nil { // unix
					err2 = backoff.Permanent(fmt.Errorf("new code crashed (pid " + strconv.Itoa(runCmd.Process.Pid) + "), fix errors"))
					return err2
				}
			}
			return err
		}
		remoteRenderer := newDevRendererClient(dialHTTP)
		// 4.1. Swap the renderer on success
		r.impl = remoteRenderer
		r.rerender() // Render the new SDF!!!
		return nil
	}, backOff, func(err error, duration time.Duration) {
		log.Println("[DevRenderer] connection error:", err, "- retrying in:", duration)
	})
	if err != nil {
		log.Println("[DevRenderer] backoff.RetryNotify gave up on connecting, with error:", err)
		return runCmd
	}
	return runCmd
}
