#import <Cocoa/Cocoa.h>
#include "windowdelegate.h"

extern void GoWindowEvent(const char* surfaceID, const char* event, const char* data);
extern NSMutableDictionary<NSString*, NSWindow*> *windowMap;

@interface JVWindowDelegate : NSObject <NSWindowDelegate>
@property (nonatomic, copy) NSString *surfaceID;
@end

@implementation JVWindowDelegate

- (void)windowDidResize:(NSNotification *)notification {
    NSWindow *w = notification.object;
    NSRect frame = w.frame;
    NSString *json = [NSString stringWithFormat:@"{\"width\":%.0f,\"height\":%.0f}",
                      frame.size.width, frame.size.height];
    GoWindowEvent([self.surfaceID UTF8String], "window.resize", [json UTF8String]);
}

- (void)windowDidMove:(NSNotification *)notification {
    NSWindow *w = notification.object;
    NSRect frame = w.frame;
    NSString *json = [NSString stringWithFormat:@"{\"x\":%.0f,\"y\":%.0f}",
                      frame.origin.x, frame.origin.y];
    GoWindowEvent([self.surfaceID UTF8String], "window.move", [json UTF8String]);
}

- (BOOL)windowShouldClose:(NSWindow *)sender {
    GoWindowEvent([self.surfaceID UTF8String], "window.beforeClose", "{}");
    // Return NO — let the event handler decide whether to actually close.
    // If no handler cancels, the default behavior is to close.
    return NO;
}

- (void)windowWillClose:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.close", "{}");
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.minimize", "{}");
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.restore", "{}");
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.fullscreenEnter", "{}");
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.fullscreenExit", "{}");
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.becomeKey", "{}");
}

- (void)windowDidResignKey:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.resignKey", "{}");
}

- (void)windowDidBecomeMain:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.becomeMain", "{}");
}

- (void)windowDidResignMain:(NSNotification *)notification {
    GoWindowEvent([self.surfaceID UTF8String], "window.resignMain", "{}");
}

- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
    NSWindow *w = notification.object;
    BOOL visible = (w.occlusionState & NSWindowOcclusionStateVisible) != 0;
    NSString *json = [NSString stringWithFormat:@"{\"visible\":%s}", visible ? "true" : "false"];
    GoWindowEvent([self.surfaceID UTF8String], "window.occlude", [json UTF8String]);
}

@end

// Map from surfaceID → JVWindowDelegate (retained to prevent dealloc)
static NSMutableDictionary<NSString*, JVWindowDelegate*> *windowDelegateMap = nil;

void JVInstallWindowDelegate(const char* surfaceID) {
    NSString *sid = [NSString stringWithUTF8String:surfaceID];
    NSWindow *window = windowMap[sid];
    if (!window) return;

    if (!windowDelegateMap) {
        windowDelegateMap = [[NSMutableDictionary alloc] init];
    }

    JVWindowDelegate *delegate = [[JVWindowDelegate alloc] init];
    delegate.surfaceID = sid;
    window.delegate = delegate;
    windowDelegateMap[sid] = delegate;
}

void JVRemoveWindowDelegate(const char* surfaceID) {
    NSString *sid = [NSString stringWithUTF8String:surfaceID];
    if (windowDelegateMap) {
        NSWindow *window = windowMap[sid];
        if (window && window.delegate == windowDelegateMap[sid]) {
            window.delegate = nil;
        }
        [windowDelegateMap removeObjectForKey:sid];
    }
}
