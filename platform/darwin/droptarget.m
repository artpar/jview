#import <Cocoa/Cocoa.h>
#include "droptarget.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kDropDelegateKey = &kDropDelegateKey;

// JVDropDelegate wraps NSDraggingDestination behavior for any NSView.
// It's attached as an associated object and uses method swizzling isn't needed —
// instead we use a transparent overlay view that acts as the drop target.
@interface JVDropOverlay : NSView
@property (nonatomic, assign) uint64_t callbackID;
@end

@implementation JVDropOverlay

- (instancetype)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        self.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
        [self registerForDraggedTypes:@[
            NSPasteboardTypeFileURL,
            NSPasteboardTypeString,
            NSPasteboardTypeURL
        ]];
    }
    return self;
}

- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
    return NSDragOperationCopy;
}

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
    NSPasteboard *pb = [sender draggingPasteboard];

    NSMutableArray *paths = [NSMutableArray array];
    NSString *text = nil;

    // Check for file URLs
    NSArray<NSURL *> *urls = [pb readObjectsForClasses:@[[NSURL class]]
                                               options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
    for (NSURL *url in urls) {
        if (url.path) [paths addObject:url.path];
    }

    // Check for plain text
    NSString *str = [pb stringForType:NSPasteboardTypeString];
    if (str) text = str;

    // Build JSON result
    NSMutableDictionary *result = [NSMutableDictionary dictionary];
    if (paths.count > 0) result[@"paths"] = paths;
    if (text) result[@"text"] = text;

    NSError *err = nil;
    NSData *data = [NSJSONSerialization dataWithJSONObject:result options:0 error:&err];
    if (!err) {
        NSString *json = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
        GoCallbackInvoke(self.callbackID, [json UTF8String]);
    }

    return YES;
}

- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
    return YES;
}

// Pass all mouse events through to the underlying view
- (NSView *)hitTest:(NSPoint)point {
    return nil;
}

@end

void JVEnableDropTarget(void* handle, uint64_t callbackID) {
    NSView *view = (__bridge NSView*)handle;

    // Remove existing overlay if any
    JVDropOverlay *existing = objc_getAssociatedObject(view, kDropDelegateKey);
    if (existing) {
        existing.callbackID = callbackID;
        return;
    }

    JVDropOverlay *overlay = [[JVDropOverlay alloc] initWithFrame:view.bounds];
    overlay.callbackID = callbackID;
    [view addSubview:overlay];

    // Prevent overlay from being deallocated
    objc_setAssociatedObject(view, kDropDelegateKey, overlay, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

void JVUpdateDropTargetCallbackID(void* handle, uint64_t callbackID) {
    NSView *view = (__bridge NSView*)handle;
    JVDropOverlay *overlay = objc_getAssociatedObject(view, kDropDelegateKey);
    if (overlay) {
        overlay.callbackID = callbackID;
    }
}

void JVDisableDropTarget(void* handle) {
    NSView *view = (__bridge NSView*)handle;
    JVDropOverlay *overlay = objc_getAssociatedObject(view, kDropDelegateKey);
    if (overlay) {
        [overlay removeFromSuperview];
        objc_setAssociatedObject(view, kDropDelegateKey, nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}
