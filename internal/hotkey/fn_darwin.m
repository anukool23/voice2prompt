//go:build darwin

#import <Cocoa/Cocoa.h>

// Exported from fn_darwin.go: 1 = Fn pressed, 0 = Fn released.
extern void goFnEvent(int down);

static CFMachPortRef gFnTap;
static CFRunLoopSourceRef gFnSrc;
static BOOL gFnDown = NO;

static CGEventRef fnCallback(CGEventTapProxy proxy, CGEventType type,
                             CGEventRef event, void *ctx) {
    // Re-enable the tap if the system disabled it.
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        if (gFnTap) CGEventTapEnable(gFnTap, true);
        return event;
    }
    if (type == kCGEventFlagsChanged) {
        BOOL fn = (CGEventGetFlags(event) & kCGEventFlagMaskSecondaryFn) != 0;
        if (fn != gFnDown) {
            gFnDown = fn;
            goFnEvent(fn ? 1 : 0);
        }
    }
    return event; // listen-only: pass the event through untouched
}

// startFnTap installs a listen-only tap for flag changes on the main run loop.
// Returns 1 on success, 0 if the tap couldn't be created (missing Input Monitoring).
int startFnTap(void) {
    __block int ok = 0;
    dispatch_sync(dispatch_get_main_queue(), ^{
        CGEventMask mask = CGEventMaskBit(kCGEventFlagsChanged);
        gFnTap = CGEventTapCreate(kCGSessionEventTap, kCGHeadInsertEventTap,
                                  kCGEventTapOptionListenOnly, mask, fnCallback, NULL);
        if (gFnTap) {
            gFnSrc = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, gFnTap, 0);
            CFRunLoopAddSource(CFRunLoopGetMain(), gFnSrc, kCFRunLoopCommonModes);
            CGEventTapEnable(gFnTap, true);
            ok = 1;
        }
    });
    return ok;
}

void stopFnTap(void) {
    dispatch_sync(dispatch_get_main_queue(), ^{
        if (gFnSrc) {
            CFRunLoopRemoveSource(CFRunLoopGetMain(), gFnSrc, kCFRunLoopCommonModes);
            CFRelease(gFnSrc);
            gFnSrc = NULL;
        }
        if (gFnTap) {
            CGEventTapEnable(gFnTap, false);
            CFRelease(gFnTap);
            gFnTap = NULL;
        }
        gFnDown = NO;
    });
}
