//go:build darwin

package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
int startFnTap(void);
void stopFnTap(void);
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
)

// Fn interaction tuning.
const (
	doubleTapWindow = 400 * time.Millisecond // max gap between taps to count as a double-tap
	quickTapMax     = 300 * time.Millisecond // a press shorter than this is a "tap", not a hold
)

// FnController turns raw Fn key up/down events into dictation control:
//
//   - Hold Fn         → push-to-talk (record while held; stop & transcribe on release)
//   - Double-tap Fn   → lock: keep recording hands-free until Fn is tapped again
//
// onStart begins a capture; onStop ends it and transcribes. Both must be safe to
// call from the main thread.
type FnController struct {
	onStart func()
	onStop  func()

	mu              sync.Mutex
	locked          bool
	pttActive       bool
	downTime        time.Time
	lastDownTime    time.Time
	lastWasQuickTap bool

	actions chan bool // true=start, false=stop; drained by a worker goroutine
}

// fnActive is the singleton the C callback dispatches to (goFnEvent is package-level).
var fnActive *FnController

// NewFn builds a controller. onStart/onStop drive the engine's capture.
func NewFn(onStart, onStop func()) *FnController {
	return &FnController{onStart: onStart, onStop: onStop}
}

//export goFnEvent
func goFnEvent(down C.int) {
	if fnActive != nil {
		fnActive.handle(down == 1)
	}
}

// Register installs the Fn event tap. Fails if Input Monitoring isn't granted.
func (f *FnController) Register() error {
	fnActive = f
	f.actions = make(chan bool, 16)
	go func() {
		for start := range f.actions {
			if start {
				f.onStart()
			} else {
				f.onStop()
			}
		}
	}()
	if C.startFnTap() == 0 {
		fnActive = nil
		close(f.actions)
		f.actions = nil
		return fmt.Errorf("could not start Fn listener — grant Input Monitoring permission and relaunch")
	}
	return nil
}

// Unregister removes the tap and stops the worker.
func (f *FnController) Unregister() {
	C.stopFnTap()
	fnActive = nil
	if f.actions != nil {
		close(f.actions)
		f.actions = nil
	}
}

// handle runs the tap→dictation state machine. Callbacks are invoked outside the
// lock (but synchronously) so start/stop ordering is preserved.
func (f *FnController) handle(down bool) {
	now := time.Now()
	const (
		none = iota
		start
		stop
	)
	action := none

	f.mu.Lock()
	if down {
		switch {
		case f.locked:
			// A press while locked ends the hands-free session.
			f.locked = false
			f.lastWasQuickTap = false
			action = stop
		case f.lastWasQuickTap && now.Sub(f.lastDownTime) < doubleTapWindow:
			// Second quick tap → lock and start a fresh hands-free capture.
			f.locked = true
			f.pttActive = false
			f.lastWasQuickTap = false
			action = start
		default:
			f.pttActive = true
			f.downTime = now
			f.lastDownTime = now
			action = start
		}
	} else if !f.locked && f.pttActive {
		f.pttActive = false
		f.lastWasQuickTap = now.Sub(f.downTime) < quickTapMax
		action = stop // very short taps become <0.2s clips the engine ignores
	}
	f.mu.Unlock()

	if action != none && f.actions != nil {
		f.actions <- (action == start)
	}
}

// Locked reports whether a hands-free (double-tap) session is active.
func (f *FnController) Locked() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.locked
}
