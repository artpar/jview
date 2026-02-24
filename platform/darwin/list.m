#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>
#include "list.h"
#include "stackview.h"

static const void *kListStackKey = &kListStackKey;

// Flipped container so the scroll view's initial scroll position shows the top of the content.
@interface JVFlippedView : NSView
@end

@implementation JVFlippedView
- (BOOL)isFlipped { return YES; }
@end

void* JVCreateList(const char* justify, const char* align, int gap, int padding) {
    const char* effectiveAlign = (align && strlen(align) > 0) ? align : "stretch";
    const char* effectiveJustify = (justify && strlen(justify) > 0) ? justify : "start";

    // Create the inner stack view (transfers ownership from JVCreateStackView's __bridge_retained)
    NSStackView *stack = (__bridge_transfer NSStackView*)JVCreateStackView(false, effectiveJustify, effectiveAlign, gap, padding);

    // Flipped container wraps the stack so content starts from the top
    JVFlippedView *container = [[JVFlippedView alloc] init];
    container.translatesAutoresizingMaskIntoConstraints = NO;
    [container addSubview:stack];

    // Pin stack to container edges — container height is determined by stack content
    [NSLayoutConstraint activateConstraints:@[
        [stack.topAnchor constraintEqualToAnchor:container.topAnchor],
        [stack.leadingAnchor constraintEqualToAnchor:container.leadingAnchor],
        [stack.trailingAnchor constraintEqualToAnchor:container.trailingAnchor],
        [stack.bottomAnchor constraintEqualToAnchor:container.bottomAnchor],
    ]];

    // Create scroll view
    NSScrollView *scrollView = [[NSScrollView alloc] init];
    scrollView.translatesAutoresizingMaskIntoConstraints = NO;
    scrollView.hasVerticalScroller = YES;
    scrollView.hasHorizontalScroller = NO;
    scrollView.autohidesScrollers = YES;
    scrollView.borderType = NSNoBorder;
    scrollView.drawsBackground = NO;

    // Set flipped container as document view
    scrollView.documentView = container;

    // Pin container width to clip view (prevents horizontal scrolling)
    [container.widthAnchor constraintEqualToAnchor:scrollView.contentView.widthAnchor].active = YES;

    // Store reference to inner stack for update/setChildren
    objc_setAssociatedObject(scrollView, kListStackKey, stack, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)scrollView;
}

void JVUpdateList(void* handle, const char* justify, const char* align, int gap, int padding) {
    NSScrollView *scrollView = (__bridge NSScrollView*)handle;
    NSStackView *stack = objc_getAssociatedObject(scrollView, kListStackKey);
    if (!stack) return;

    const char* effectiveAlign = (align && strlen(align) > 0) ? align : "stretch";
    const char* effectiveJustify = (justify && strlen(justify) > 0) ? justify : "start";
    JVUpdateStackView((__bridge void*)stack, effectiveJustify, effectiveAlign, gap, padding);
}

void JVListSetChildren(void* handle, void** children, int count) {
    NSScrollView *scrollView = (__bridge NSScrollView*)handle;
    NSStackView *stack = objc_getAssociatedObject(scrollView, kListStackKey);
    if (!stack) return;
    JVStackViewSetChildren((__bridge void*)stack, children, count);

    // Force layout then scroll to top of content
    [stack.superview layoutSubtreeIfNeeded];
    [scrollView.documentView scrollPoint:NSMakePoint(0, 0)];
}
