#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>
#include "splitview.h"

static const void *kSplitViewDelegateKey = &kSplitViewDelegateKey;

@interface JVSplitViewDelegate : NSObject <NSSplitViewDelegate>
@end

@implementation JVSplitViewDelegate

- (CGFloat)splitView:(NSSplitView *)splitView constrainMinCoordinate:(CGFloat)proposedMinimumPosition ofSubviewAt:(NSInteger)dividerIndex {
    return proposedMinimumPosition + 100;
}

- (CGFloat)splitView:(NSSplitView *)splitView constrainMaxCoordinate:(CGFloat)proposedMaximumPosition ofSubviewAt:(NSInteger)dividerIndex {
    return proposedMaximumPosition - 100;
}

@end

void* JVCreateSplitView(const char* dividerStyle, bool vertical) {
    NSSplitView *splitView = [[NSSplitView alloc] init];
    splitView.translatesAutoresizingMaskIntoConstraints = NO;
    splitView.vertical = vertical;

    NSString *styleStr = [NSString stringWithUTF8String:dividerStyle];
    if ([styleStr isEqualToString:@"thick"]) {
        splitView.dividerStyle = NSSplitViewDividerStyleThick;
    } else if ([styleStr isEqualToString:@"paneSplitter"]) {
        splitView.dividerStyle = NSSplitViewDividerStylePaneSplitter;
    } else {
        splitView.dividerStyle = NSSplitViewDividerStyleThin;
    }

    JVSplitViewDelegate *delegate = [[JVSplitViewDelegate alloc] init];
    splitView.delegate = delegate;
    objc_setAssociatedObject(splitView, kSplitViewDelegateKey, delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)splitView;
}

void JVUpdateSplitView(void* handle, const char* dividerStyle, bool vertical) {
    NSSplitView *splitView = (__bridge NSSplitView*)handle;
    splitView.vertical = vertical;

    NSString *styleStr = [NSString stringWithUTF8String:dividerStyle];
    if ([styleStr isEqualToString:@"thick"]) {
        splitView.dividerStyle = NSSplitViewDividerStyleThick;
    } else if ([styleStr isEqualToString:@"paneSplitter"]) {
        splitView.dividerStyle = NSSplitViewDividerStylePaneSplitter;
    } else {
        splitView.dividerStyle = NSSplitViewDividerStyleThin;
    }
}

void JVSplitViewSetChildren(void* handle, void** children, int count) {
    NSSplitView *splitView = (__bridge NSSplitView*)handle;

    // Remove existing arranged subviews
    NSArray<NSView*> *existing = [splitView.arrangedSubviews copy];
    for (NSView *view in existing) {
        [splitView removeArrangedSubview:view];
        [view removeFromSuperview];
    }

    // Add new children
    for (int i = 0; i < count; i++) {
        NSView *child = (__bridge NSView*)children[i];
        [splitView addArrangedSubview:child];
    }

    // Set initial proportional positions after children are added
    if (count > 1) {
        [splitView adjustSubviews];
    }
}
