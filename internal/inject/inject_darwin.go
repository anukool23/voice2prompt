//go:build darwin

// Package inject places transcribed text into the focused application.
//
// Primary path: the macOS Accessibility (AX) API — find the focused UI element and
// insert text at the caret. This preserves the clipboard and lets us read
// surrounding context (used by the Phase 2 LLM). AX write support is uneven
// (native Cocoa fields work well; some browser/Electron web inputs don't), so when
// AX insertion isn't possible we fall back to clipboard + ⌘V.
//
// Both paths require Accessibility permission (System Settings → Privacy & Security
// → Accessibility) for the running app/terminal.
package inject

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation -framework AVFoundation -framework IOKit
#include <ApplicationServices/ApplicationServices.h>
#include <stdlib.h>

// Implemented in mic_darwin.m (AVFoundation).
int mic_status(void);
void mic_request(void);

// Input Monitoring (needed by the Fn CGEventTap). IOKit funcs, declared here to
// avoid header-path fragility. kIOHIDRequestTypeListenEvent == 1.
extern int IOHIDCheckAccess(int requestType);
extern int IOHIDRequestAccess(int requestType);
static int input_monitoring_status(void) { return IOHIDCheckAccess(1); } // 0=granted,1=denied,2=unknown
static void input_monitoring_request(void) { IOHIDRequestAccess(1); }

// Returns 1 if the process is trusted for Accessibility, else 0.
static int ax_is_trusted(void) {
    return AXIsProcessTrusted() ? 1 : 0;
}

// Triggers the system Accessibility prompt (adds the app to the list, greyed off).
static void ax_prompt_trust(void) {
    const void* keys[]   = { kAXTrustedCheckOptionPrompt };
    const void* values[] = { kCFBooleanTrue };
    CFDictionaryRef opts = CFDictionaryCreate(NULL, keys, values, 1,
        &kCFTypeDictionaryKeyCallBacks, &kCFTypeDictionaryValueCallBacks);
    AXIsProcessTrustedWithOptions(opts);
    CFRelease(opts);
}

// Copies the currently focused UI element into *out. Caller must CFRelease.
static AXError copy_focused_element(AXUIElementRef* out) {
    AXUIElementRef sys = AXUIElementCreateSystemWide();
    AXError err = AXUIElementCopyAttributeValue(sys, kAXFocusedUIElementAttribute,
        (CFTypeRef*)out);
    CFRelease(sys);
    return err;
}

// Inserts utf8 at the caret of the focused element by setting AXSelectedText
// (this replaces any current selection, or inserts at the caret if none).
// Returns 1 on success, 0 if there's no focused element or it rejects the write.
static int ax_insert_text(const char* utf8) {
    AXUIElementRef focused = NULL;
    if (copy_focused_element(&focused) != kAXErrorSuccess || focused == NULL) {
        return 0;
    }
    CFStringRef str = CFStringCreateWithCString(NULL, utf8, kCFStringEncodingUTF8);
    AXError setErr = AXUIElementSetAttributeValue(focused, kAXSelectedTextAttribute,
        (CFTypeRef)str);
    CFRelease(str);
    CFRelease(focused);
    return (setErr == kAXErrorSuccess) ? 1 : 0;
}

// Returns a malloc'd UTF-8 copy of a CFStringRef, or NULL. Caller must free().
static char* cfstring_dup(CFStringRef s) {
    if (s == NULL || CFGetTypeID(s) != CFStringGetTypeID()) return NULL;
    CFIndex maxLen = CFStringGetMaximumSizeForEncoding(CFStringGetLength(s),
        kCFStringEncodingUTF8) + 1;
    char* buf = (char*)malloc(maxLen);
    if (buf == NULL || !CFStringGetCString(s, buf, maxLen, kCFStringEncodingUTF8)) {
        if (buf) free(buf);
        return NULL;
    }
    return buf;
}

// Returns a malloc'd UTF-8 name of the frontmost (focused) application, or NULL.
// Caller must free().
static char* ax_focused_app_name(void) {
    AXUIElementRef sys = AXUIElementCreateSystemWide();
    AXUIElementRef app = NULL;
    AXError err = AXUIElementCopyAttributeValue(sys, kAXFocusedApplicationAttribute,
        (CFTypeRef*)&app);
    CFRelease(sys);
    if (err != kAXErrorSuccess || app == NULL) return NULL;

    CFTypeRef title = NULL;
    err = AXUIElementCopyAttributeValue(app, kAXTitleAttribute, &title);
    CFRelease(app);
    if (err != kAXErrorSuccess || title == NULL) return NULL;

    char* out = cfstring_dup((CFStringRef)title);
    CFRelease(title);
    return out;
}

// Returns a malloc'd UTF-8 copy of the focused element's text value, or NULL.
// Caller must free().
static char* ax_focused_value(void) {
    AXUIElementRef focused = NULL;
    if (copy_focused_element(&focused) != kAXErrorSuccess || focused == NULL) {
        return NULL;
    }
    CFTypeRef value = NULL;
    AXError e = AXUIElementCopyAttributeValue(focused, kAXValueAttribute, &value);
    CFRelease(focused);
    if (e != kAXErrorSuccess || value == NULL) {
        return NULL;
    }
    if (CFGetTypeID(value) != CFStringGetTypeID()) {
        CFRelease(value);
        return NULL;
    }
    CFStringRef s = (CFStringRef)value;
    CFIndex maxLen = CFStringGetMaximumSizeForEncoding(CFStringGetLength(s),
        kCFStringEncodingUTF8) + 1;
    char* buf = (char*)malloc(maxLen);
    if (buf == NULL || !CFStringGetCString(s, buf, maxLen, kCFStringEncodingUTF8)) {
        if (buf) free(buf);
        CFRelease(value);
        return NULL;
    }
    CFRelease(value);
    return buf;
}

// Synthesizes a ⌘V keystroke via CGEvent. Needs only Accessibility permission
// (no AppleScript/Automation permission, unlike osascript), and is far faster.
// kVK_ANSI_V == 9.
static void cg_paste(void) {
    CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
    CGEventRef down = CGEventCreateKeyboardEvent(src, (CGKeyCode)9, true);
    CGEventRef up   = CGEventCreateKeyboardEvent(src, (CGKeyCode)9, false);
    CGEventSetFlags(down, kCGEventFlagMaskCommand);
    CGEventSetFlags(up, kCGEventFlagMaskCommand);
    CGEventPost(kCGHIDEventTap, down);
    CGEventPost(kCGHIDEventTap, up);
    CFRelease(down);
    CFRelease(up);
    if (src) CFRelease(src);
}

// Synthesizes an arbitrary keystroke with modifier flags (for voice commands).
static void cg_key(int keycode, int cmd, int opt, int shift, int ctrl) {
    CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
    CGEventFlags flags = 0;
    if (cmd)   flags |= kCGEventFlagMaskCommand;
    if (opt)   flags |= kCGEventFlagMaskAlternate;
    if (shift) flags |= kCGEventFlagMaskShift;
    if (ctrl)  flags |= kCGEventFlagMaskControl;
    CGEventRef down = CGEventCreateKeyboardEvent(src, (CGKeyCode)keycode, true);
    CGEventRef up   = CGEventCreateKeyboardEvent(src, (CGKeyCode)keycode, false);
    CGEventSetFlags(down, flags);
    CGEventSetFlags(up, flags);
    CGEventPost(kCGHIDEventTap, down);
    CGEventPost(kCGHIDEventTap, up);
    CFRelease(down);
    CFRelease(up);
    if (src) CFRelease(src);
}

// Returns a malloc'd UTF-8 copy of the focused element's selected text, or NULL.
static char* ax_selected_text(void) {
    AXUIElementRef focused = NULL;
    if (copy_focused_element(&focused) != kAXErrorSuccess || focused == NULL) {
        return NULL;
    }
    CFTypeRef value = NULL;
    AXError e = AXUIElementCopyAttributeValue(focused, kAXSelectedTextAttribute, &value);
    CFRelease(focused);
    if (e != kAXErrorSuccess || value == NULL) {
        return NULL;
    }
    char* out = cfstring_dup((CFStringRef)value);
    CFRelease(value);
    return out;
}
*/
import "C"

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
	"unsafe"
)

// Method reports which injection path was used, for logging/telemetry.
type Method string

const (
	MethodAX        Method = "accessibility"
	MethodClipboard Method = "clipboard"
)

// Trusted reports whether the process has Accessibility permission.
func Trusted() bool {
	return C.ax_is_trusted() == 1
}

// PromptTrust asks macOS to show the Accessibility prompt for this app.
func PromptTrust() {
	C.ax_prompt_trust()
}

// MicStatus reports microphone permission: "authorized", "denied", "restricted",
// or "undetermined".
func MicStatus() string {
	switch int(C.mic_status()) {
	case 3:
		return "authorized"
	case 2:
		return "denied"
	case 1:
		return "restricted"
	default:
		return "undetermined"
	}
}

// MicAuthorized reports whether microphone access is granted.
func MicAuthorized() bool { return int(C.mic_status()) == 3 }

// RequestMic triggers the microphone permission prompt.
func RequestMic() { C.mic_request() }

// InputMonitoringStatus reports "authorized", "denied", or "undetermined".
// Required for the Fn-key CGEventTap.
func InputMonitoringStatus() string {
	switch int(C.input_monitoring_status()) {
	case 0:
		return "authorized"
	case 1:
		return "denied"
	default:
		return "undetermined"
	}
}

// InputMonitoringAuthorized reports whether Input Monitoring is granted.
func InputMonitoringAuthorized() bool { return int(C.input_monitoring_status()) == 0 }

// RequestInputMonitoring triggers the Input Monitoring permission prompt.
func RequestInputMonitoring() { C.input_monitoring_request() }

// FocusedApp returns the name of the frontmost application (best effort, "" if unknown).
// Used to pick context-aware formatting (casual for chat apps, formal for mail/docs).
func FocusedApp() string {
	cstr := C.ax_focused_app_name()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

// FocusedValue returns the current text content of the focused element (best effort).
// Empty string if unavailable. Intended as context for the Phase 2 cleanup LLM.
func FocusedValue() string {
	cstr := C.ax_focused_value()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

// macOS ANSI virtual key codes used by voice commands.
const (
	keyA      = 0
	keyZ      = 6
	keyX      = 7
	keyC      = 8
	keyV      = 9
	keyReturn = 36
	keyDelete = 51 // backspace
)

func cbool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

// KeyPress synthesizes a keystroke with the given modifiers (voice commands).
func KeyPress(keyCode int, cmd, opt, shift, ctrl bool) {
	C.cg_key(C.int(keyCode), cbool(cmd), cbool(opt), cbool(shift), cbool(ctrl))
}

// SelectedText returns the focused element's selected text (best effort, "" if none).
func SelectedText() string {
	cstr := C.ax_selected_text()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

// Named editing commands (⌘-shortcuts / edits synthesized via CGEvent).
func SelectAll()    { KeyPress(keyA, true, false, false, false) }
func Undo()         { KeyPress(keyZ, true, false, false, false) }
func Redo()         { KeyPress(keyZ, true, false, true, false) }
func Copy()         { KeyPress(keyC, true, false, false, false) }
func Cut()          { KeyPress(keyX, true, false, false, false) }
func PasteKey()     { KeyPress(keyV, true, false, false, false) }
func NewLine()      { KeyPress(keyReturn, false, false, false, false) }
func NewParagraph() { NewLine(); NewLine() }
func DeleteWord()   { KeyPress(keyDelete, false, true, false, false) } // Option+Delete
func DeleteLine()   { KeyPress(keyDelete, true, false, false, false) } // Cmd+Delete

// Paste inserts text into the focused field, preferring the AX API and falling
// back to clipboard + ⌘V. It returns which method succeeded.
func Paste(text string) (Method, error) {
	if text == "" {
		return "", nil
	}

	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	if C.ax_insert_text(cstr) == 1 {
		return MethodAX, nil
	}

	// AX unsupported for this element — fall back to clipboard paste.
	if err := pasteViaClipboard(text); err != nil {
		return "", err
	}
	return MethodClipboard, nil
}

// pasteViaClipboard sets the clipboard, sends ⌘V via CGEvent, then restores the
// previous clipboard text asynchronously (so the restore delay doesn't inflate the
// user-visible latency). A non-text clipboard (e.g. an image) is not preserved.
func pasteViaClipboard(text string) error {
	prev := readClipboard() // best effort; "" if empty or non-text

	if err := writeClipboard(text); err != nil {
		return err
	}
	// Brief settle so the pasteboard write is visible to the target app.
	time.Sleep(20 * time.Millisecond)

	C.cg_paste()

	// Restore the old clipboard in the background, after the app has consumed the paste.
	if prev != "" {
		go func() {
			time.Sleep(200 * time.Millisecond)
			_ = writeClipboard(prev)
		}()
	}
	return nil
}

func readClipboard() string {
	out, err := exec.Command("pbpaste").Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func writeClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pbcopy failed: %w", err)
	}
	return nil
}
