package dev

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/hajimehoshi/ebiten"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const changeEventThrottle = 100 * time.Millisecond

func (r *Renderer) runRenderer(runCmdF func() *exec.Cmd, watchFiles []string) error {
	if len(watchFiles) > 0 {
		watcher, err := NewFsWatcher()
		if err != nil {
			log.Println("Error watching files (won't update on changes):", err)
		} else {
			defer func(watcher io.Closer) {
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

			for _, matchedFile := range watchFiles {
				err = watcher.Add(matchedFile)
				if err != nil {
					log.Println("Error watching file", matchedFile, "-", err)
				}
			}
		}
	}

	return ebiten.RunGame(r) // blocks until the window is closed
}

func (r *Renderer) runChild(requestedAddress string) error {
	// Listen for signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, signals()...)
	// Set up a remote service that the parent renderer will connect to view the new SDF
	service := newDevRendererService(r.impl, done)
	service.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	// TODO: Use service.ServeConn() on a pipe to the parent, avoiding using ports (must be as cross-platform as possible)
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
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			log.Println("[DevRenderer] srv.Serve error:", err)
		}
		done <- syscall.SIGKILL
	}()
	log.Println("[DevRenderer] Child service ready...")
	<-done // Will block until interrupt is received or the server crashes
	log.Println("[DevRenderer] Child service finished successfully...")
	return nil
}

func (r *Renderer) rendererSwapChild(runCmd *exec.Cmd, runCmdF func() *exec.Cmd) *exec.Cmd {
	r.implLock.Lock() // No more renders until we swapped the implementation
	defer r.implLock.Unlock()
	// 1. Gracefully close the previous command
	if runCmd != nil {
		log.Println("[DevRenderer] Closing previous child process")
		if rend, ok := r.impl.(*rendererClient); ok {
			err := rend.Shutdown(5 * time.Second)
			if err != nil {
				log.Println("[DevRenderer] Closing previous child process ERROR:", err, "(the child will probably keep running in background)")
			}
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
	// Note that in case of "go run ...", a new process is forked after successful compilation and the runCmd PID will die.
	startupFinished := make(chan *os.ProcessState) // true if success
	go func() {
		ps, err2 := runCmd.Process.Wait()
		if err2 != nil {
			log.Println("[DevRenderer] runCmd error:", err2)
		}
		startupFinished <- ps
		close(startupFinished)
	}()
	// 4. Connect to it as fast as possible, with exponential backoff to relax on errors.
	log.Println("[DevRenderer] Trying to connect to new code with exponential backoff...")
	r.backOff.Reset()
	err = backoff.RetryNotify(func() error {
		dialHTTP, err := rpc.DialHTTP("tcp", requestedFreeAddr)
		if err != nil {
			select {
			case ps, ok := <-startupFinished:
				if ok && !ps.Success() {
					err2 := backoff.Permanent(fmt.Errorf("new code crashed (pid " + strconv.Itoa(runCmd.Process.Pid) +
						"), fix errors: " + ps.String()))
					return err2
				}
			default: // Do not block checking if process success
			}
			return err
		}
		remoteRenderer := newDevRendererClient(dialHTTP)
		// 4.1. Swap the renderer on success
		r.impl = remoteRenderer
		r.rerender() // Render the new SDF!!!
		return nil
	}, r.backOff, func(err error, duration time.Duration) {
		log.Println("[DevRenderer] connection error:", err, "- retrying in:", duration)
	})
	if err != nil {
		log.Println("[DevRenderer] backoff.RetryNotify gave up on connecting, with error:", err)
		return runCmd
	}
	return runCmd
}
