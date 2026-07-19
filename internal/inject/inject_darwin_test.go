//go:build darwin

package inject

import "testing"

// Smoke test: the cgo AX bridge should execute without crashing. In an untrusted
// or headless context Trusted() is false and FocusedValue() is empty — that's fine;
// we're only verifying the CoreFoundation memory handling doesn't panic/leak-crash.
func TestAXBridgeSmoke(t *testing.T) {
	t.Logf("Trusted() = %v", Trusted())
	t.Logf("FocusedValue() = %q", FocusedValue())
}
