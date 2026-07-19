//go:build darwin

// Package tray adds a macOS menu-bar (NSStatusItem) icon with a menu. It creates
// the status item on the main queue so it coexists with the Wails NSApplication
// run loop, and routes menu clicks back to Go callbacks.
package tray

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <stdlib.h>
void trayStart(const char* iconPath);
*/
import "C"

import "unsafe"

var (
	onOpen   func()
	onToggle func()
	onQuit   func()
)

//export goTrayOpen
func goTrayOpen() {
	if onOpen != nil {
		onOpen()
	}
}

//export goTrayToggle
func goTrayToggle() {
	if onToggle != nil {
		onToggle()
	}
}

//export goTrayQuit
func goTrayQuit() {
	if onQuit != nil {
		onQuit()
	}
}

// Start installs the menu-bar item using iconPath (a template PNG) and wires the
// menu actions. Safe to call from any goroutine; the item is created on the main queue.
func Start(iconPath string, open, toggle, quit func()) {
	onOpen, onToggle, onQuit = open, toggle, quit
	c := C.CString(iconPath)
	defer C.free(unsafe.Pointer(c))
	C.trayStart(c)
}
