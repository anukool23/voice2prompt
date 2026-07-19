//go:build darwin

package hotkey

// Carbon-based push-to-talk that coexists with an already-running Cocoa/Carbon run
// loop (e.g. the one Wails owns). Unlike golang.design/x/hotkey it does NOT start
// its own main-thread loop — it installs a handler on the application event target
// and relies on the host's loop to deliver events. Uses RegisterEventHotKey, so it
// needs no Input Monitoring permission.

/*
#cgo LDFLAGS: -framework Carbon
#include <Carbon/Carbon.h>

extern void goHotkeyEvent(int pressed);

static EventHandlerRef gHandler;
static EventHotKeyRef  gHotKey;

static OSStatus hkHandler(EventHandlerCallRef next, EventRef evt, void* userData) {
    goHotkeyEvent(GetEventKind(evt) == kEventHotKeyPressed ? 1 : 0);
    return noErr;
}

static int carbonRegister(unsigned int keyCode, unsigned int mods) {
    EventTypeSpec types[2] = {
        { kEventClassKeyboard, kEventHotKeyPressed },
        { kEventClassKeyboard, kEventHotKeyReleased },
    };
    if (gHandler == NULL) {
        InstallEventHandler(GetApplicationEventTarget(),
            NewEventHandlerUPP(hkHandler), 2, types, NULL, &gHandler);
    }
    EventHotKeyID hkID = { 'PVhk', 1 };
    OSStatus s = RegisterEventHotKey(keyCode, mods, hkID,
        GetApplicationEventTarget(), 0, &gHotKey);
    return (s == noErr) ? 1 : 0;
}

static void carbonUnregister(void) {
    if (gHotKey != NULL) {
        UnregisterEventHotKey(gHotKey);
        gHotKey = NULL;
    }
}
*/
import "C"

import "sync"

// Carbon modifier masks (from Carbon/Events.h).
const (
	modCmd     = 1 << 8
	modShift   = 1 << 9
	modOption  = 1 << 11
	modControl = 1 << 12
)

var (
	cbMu     sync.Mutex
	onDownCB func()
	onUpCB   func()
)

//export goHotkeyEvent
func goHotkeyEvent(pressed C.int) {
	cbMu.Lock()
	down, up := onDownCB, onUpCB
	cbMu.Unlock()
	if pressed == 1 {
		if down != nil {
			down()
		}
	} else if up != nil {
		up()
	}
}

// Carbon is a push-to-talk hotkey manager for use inside a host run loop (Wails).
type Carbon struct{}

// NewCarbon returns a Carbon hotkey manager.
func NewCarbon() *Carbon { return &Carbon{} }

// Register binds the chord (e.g. "Ctrl+Option+Space") and wires down/up callbacks.
// Returns an error if the chord can't be parsed or registration fails.
func (c *Carbon) Register(chord string, onDown, onUp func()) error {
	keyCode, mods, err := ParseChord(chord)
	if err != nil {
		return err
	}
	cbMu.Lock()
	onDownCB, onUpCB = onDown, onUp
	cbMu.Unlock()
	if C.carbonRegister(C.uint(keyCode), C.uint(mods)) != 1 {
		return errRegister
	}
	return nil
}

// Unregister removes the hotkey.
func (c *Carbon) Unregister() {
	C.carbonUnregister()
	cbMu.Lock()
	onDownCB, onUpCB = nil, nil
	cbMu.Unlock()
}
