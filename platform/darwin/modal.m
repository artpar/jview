#import <Cocoa/Cocoa.h>
#include "modal.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);
extern NSMutableDictionary<NSString*, NSWindow*> *windowMap;

static const void *kModalPanelKey = &kModalPanelKey;
static const void *kModalContentStackKey = &kModalContentStackKey;
static const void *kModalDelegateKey = &kModalDelegateKey;

@interface JVModalDelegate : NSObject <NSWindowDelegate>
@property (nonatomic, assign) uint64_t dismissCbID;
@end

@implementation JVModalDelegate

- (BOOL)windowShouldClose:(NSWindow *)sender {
    if (self.dismissCbID != 0) {
        GoCallbackInvoke(self.dismissCbID, "");
    }
    // Return NO — let the data model drive visibility via UpdateModal
    return NO;
}

@end

void* JVCreateModal(const char* title, bool visible, const char* surfaceID, int width, int height, uint64_t dismissCbID) {
    NSString *titleStr = [NSString stringWithUTF8String:title];
    NSString *sid = [NSString stringWithUTF8String:surfaceID];

    // Create hidden proxy view (zero-height, participates in component tree)
    NSView *proxy = [[NSView alloc] initWithFrame:NSZeroRect];
    proxy.translatesAutoresizingMaskIntoConstraints = NO;
    proxy.hidden = YES;
    [NSLayoutConstraint activateConstraints:@[
        [proxy.heightAnchor constraintEqualToConstant:0],
    ]];

    // Determine panel size
    int panelWidth = (width > 0) ? width : 480;
    int panelHeight = (height > 0) ? height : 320;

    // Create NSPanel
    NSPanel *panel = [[NSPanel alloc] initWithContentRect:NSMakeRect(0, 0, panelWidth, panelHeight)
                                                styleMask:(NSWindowStyleMaskTitled |
                                                           NSWindowStyleMaskClosable |
                                                           NSWindowStyleMaskResizable)
                                                  backing:NSBackingStoreBuffered
                                                    defer:YES];
    panel.title = titleStr;
    panel.floatingPanel = YES;
    panel.becomesKeyOnlyIfNeeded = NO;
    panel.releasedWhenClosed = NO;

    // Center relative to parent window
    NSWindow *parentWindow = windowMap[sid];
    if (parentWindow) {
        NSRect parentFrame = parentWindow.frame;
        CGFloat x = NSMidX(parentFrame) - panelWidth / 2.0;
        CGFloat y = NSMidY(parentFrame) - panelHeight / 2.0;
        [panel setFrameOrigin:NSMakePoint(x, y)];
    } else {
        [panel center];
    }

    // Set up content stack view inside the panel
    NSStackView *contentStack = [[NSStackView alloc] init];
    contentStack.orientation = NSUserInterfaceLayoutOrientationVertical;
    contentStack.spacing = 8;
    contentStack.translatesAutoresizingMaskIntoConstraints = NO;

    NSView *cv = panel.contentView;
    [cv addSubview:contentStack];
    [NSLayoutConstraint activateConstraints:@[
        [contentStack.topAnchor constraintEqualToAnchor:cv.topAnchor constant:16],
        [contentStack.leadingAnchor constraintEqualToAnchor:cv.leadingAnchor constant:16],
        [contentStack.trailingAnchor constraintEqualToAnchor:cv.trailingAnchor constant:-16],
        [contentStack.bottomAnchor constraintLessThanOrEqualToAnchor:cv.bottomAnchor constant:-16],
    ]];

    // Set delegate for close button
    JVModalDelegate *delegate = [[JVModalDelegate alloc] init];
    delegate.dismissCbID = dismissCbID;
    panel.delegate = delegate;

    // Associate objects with proxy view
    objc_setAssociatedObject(proxy, kModalPanelKey, panel, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(proxy, kModalContentStackKey, contentStack, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(proxy, kModalDelegateKey, delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Show or hide based on visible
    if (visible) {
        [panel makeKeyAndOrderFront:nil];
    }

    return (__bridge_retained void*)proxy;
}

void JVUpdateModal(void* handle, const char* title, bool visible) {
    if (!handle) return;
    NSView *proxy = (__bridge NSView*)handle;
    NSPanel *panel = objc_getAssociatedObject(proxy, kModalPanelKey);
    if (!panel) return;

    NSString *titleStr = [NSString stringWithUTF8String:title];
    panel.title = titleStr;

    if (visible) {
        if (![panel isVisible]) {
            [panel makeKeyAndOrderFront:nil];
        }
    } else {
        if ([panel isVisible]) {
            [panel orderOut:nil];
        }
    }
}

void JVModalSetChildren(void* handle, void** children, int count) {
    if (!handle) return;
    NSView *proxy = (__bridge NSView*)handle;
    NSStackView *contentStack = objc_getAssociatedObject(proxy, kModalContentStackKey);
    if (!contentStack) return;

    // Remove existing arranged subviews
    NSArray<NSView*> *existing = [contentStack.arrangedSubviews copy];
    for (NSView *v in existing) {
        [contentStack removeArrangedSubview:v];
        [v removeFromSuperview];
    }

    // Add new children
    for (int i = 0; i < count; i++) {
        NSView *child = (__bridge NSView*)children[i];
        [contentStack addArrangedSubview:child];
    }
}

void JVCleanupModal(void* handle) {
    if (!handle) return;
    NSView *proxy = (__bridge NSView*)handle;
    NSPanel *panel = objc_getAssociatedObject(proxy, kModalPanelKey);
    if (panel) {
        panel.delegate = nil;
        [panel orderOut:nil];
    }
}
