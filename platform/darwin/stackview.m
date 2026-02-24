#import <Cocoa/Cocoa.h>
#include "stackview.h"

static NSStackViewGravity gravityForJustify(NSString *justify, bool horizontal) {
    // For now, gravity affects placement in the stack
    // NSStackView uses distribution rather than individual gravity per-item
    return NSStackViewGravityCenter;
}

static NSLayoutAttribute alignmentForAlign(NSString *align, bool horizontal) {
    if ([align isEqualToString:@"start"]) {
        return horizontal ? NSLayoutAttributeTop : NSLayoutAttributeLeading;
    }
    if ([align isEqualToString:@"end"]) {
        return horizontal ? NSLayoutAttributeBottom : NSLayoutAttributeTrailing;
    }
    if ([align isEqualToString:@"stretch"]) {
        return horizontal ? NSLayoutAttributeTop : NSLayoutAttributeLeading;
    }
    // center (default)
    return horizontal ? NSLayoutAttributeCenterY : NSLayoutAttributeCenterX;
}

static void applyDistribution(NSStackView *stack, NSString *justify) {
    if ([justify isEqualToString:@"spaceBetween"]) {
        stack.distribution = NSStackViewDistributionEqualSpacing;
    } else if ([justify isEqualToString:@"spaceAround"]) {
        stack.distribution = NSStackViewDistributionEqualSpacing;
    } else if ([justify isEqualToString:@"center"]) {
        stack.distribution = NSStackViewDistributionGravityAreas;
    } else {
        stack.distribution = NSStackViewDistributionFill;
    }
}

void* JVCreateStackView(bool horizontal, const char* justify, const char* align, int gap, int padding) {
    NSStackView *stack = [[NSStackView alloc] init];
    stack.orientation = horizontal ? NSUserInterfaceLayoutOrientationHorizontal
                                   : NSUserInterfaceLayoutOrientationVertical;
    stack.spacing = gap;
    stack.edgeInsets = NSEdgeInsetsMake(padding, padding, padding, padding);
    stack.translatesAutoresizingMaskIntoConstraints = NO;

    NSString *justifyStr = [NSString stringWithUTF8String:justify];
    NSString *alignStr = [NSString stringWithUTF8String:align];

    applyDistribution(stack, justifyStr);

    if ([alignStr isEqualToString:@"stretch"]) {
        // For stretch alignment in vertical stacks, arranged subviews should fill width
        stack.alignment = horizontal ? NSLayoutAttributeTop : NSLayoutAttributeWidth;
    } else {
        stack.alignment = alignmentForAlign(alignStr, horizontal);
    }

    return (__bridge_retained void*)stack;
}

void JVUpdateStackView(void* handle, const char* justify, const char* align, int gap, int padding) {
    NSStackView *stack = (__bridge NSStackView*)handle;
    stack.spacing = gap;
    stack.edgeInsets = NSEdgeInsetsMake(padding, padding, padding, padding);

    NSString *justifyStr = [NSString stringWithUTF8String:justify];
    NSString *alignStr = [NSString stringWithUTF8String:align];

    applyDistribution(stack, justifyStr);

    bool horizontal = (stack.orientation == NSUserInterfaceLayoutOrientationHorizontal);
    if ([alignStr isEqualToString:@"stretch"]) {
        stack.alignment = horizontal ? NSLayoutAttributeTop : NSLayoutAttributeWidth;
    } else {
        stack.alignment = alignmentForAlign(alignStr, horizontal);
    }
}

void JVStackViewSetChildren(void* handle, void** children, int count) {
    NSStackView *stack = (__bridge NSStackView*)handle;

    // Remove all existing arranged subviews
    NSArray<NSView*> *existing = [stack.arrangedSubviews copy];
    for (NSView *v in existing) {
        [stack removeArrangedSubview:v];
        [v removeFromSuperview];
    }

    // Add new children
    for (int i = 0; i < count; i++) {
        NSView *child = (__bridge NSView*)children[i];
        [stack addArrangedSubview:child];
    }
}
