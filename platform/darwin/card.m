#import <Cocoa/Cocoa.h>
#include "card.h"
#import <objc/runtime.h>

static const void *kContentStackKey = &kContentStackKey;
static const void *kSubtitleLabelKey = &kSubtitleLabelKey;

void* JVCreateCard(const char* title, const char* subtitle, int padding) {
    NSString *titleStr = [NSString stringWithUTF8String:title];

    NSBox *box = [[NSBox alloc] init];
    box.boxType = NSBoxPrimary;
    box.titlePosition = (titleStr.length > 0) ? NSAtTop : NSNoTitle;
    box.title = titleStr;
    box.translatesAutoresizingMaskIntoConstraints = NO;
    box.contentViewMargins = NSMakeSize(padding > 0 ? padding : 16, padding > 0 ? padding : 16);

    // Create a stack view inside the box's content view for children
    NSStackView *contentStack = [[NSStackView alloc] init];
    contentStack.orientation = NSUserInterfaceLayoutOrientationVertical;
    contentStack.spacing = 10;
    contentStack.translatesAutoresizingMaskIntoConstraints = NO;

    NSView *cv = box.contentView;
    [cv addSubview:contentStack];

    // Pin stack to the box's content view
    [NSLayoutConstraint activateConstraints:@[
        [contentStack.topAnchor constraintEqualToAnchor:cv.topAnchor],
        [contentStack.leadingAnchor constraintEqualToAnchor:cv.leadingAnchor],
        [contentStack.trailingAnchor constraintEqualToAnchor:cv.trailingAnchor],
        [contentStack.bottomAnchor constraintEqualToAnchor:cv.bottomAnchor],
    ]];

    // Render subtitle if present
    NSString *subtitleStr = [NSString stringWithUTF8String:subtitle];
    if (subtitleStr.length > 0) {
        NSTextField *subtitleLabel = [NSTextField labelWithString:subtitleStr];
        subtitleLabel.font = [NSFont systemFontOfSize:11];
        subtitleLabel.textColor = [NSColor secondaryLabelColor];
        subtitleLabel.lineBreakMode = NSLineBreakByWordWrapping;
        subtitleLabel.maximumNumberOfLines = 0;
        [contentStack addArrangedSubview:subtitleLabel];
        objc_setAssociatedObject(box, kSubtitleLabelKey, subtitleLabel, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }

    // Associate the content stack for later child management
    objc_setAssociatedObject(box, kContentStackKey, contentStack, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)box;
}

void JVUpdateCard(void* handle, const char* title, const char* subtitle, int padding) {
    NSBox *box = (__bridge NSBox*)handle;
    NSString *titleStr = [NSString stringWithUTF8String:title];
    NSString *subtitleStr = [NSString stringWithUTF8String:subtitle];

    box.title = titleStr;
    box.titlePosition = (titleStr.length > 0) ? NSAtTop : NSNoTitle;
    box.contentViewMargins = NSMakeSize(padding > 0 ? padding : 16, padding > 0 ? padding : 16);

    // Update or create/remove subtitle
    NSStackView *contentStack = objc_getAssociatedObject(box, kContentStackKey);
    NSTextField *subtitleLabel = objc_getAssociatedObject(box, kSubtitleLabelKey);
    if (subtitleStr.length > 0) {
        if (subtitleLabel) {
            subtitleLabel.stringValue = subtitleStr;
        } else if (contentStack) {
            subtitleLabel = [NSTextField labelWithString:subtitleStr];
            subtitleLabel.font = [NSFont systemFontOfSize:11];
            subtitleLabel.textColor = [NSColor secondaryLabelColor];
            subtitleLabel.lineBreakMode = NSLineBreakByWordWrapping;
            subtitleLabel.maximumNumberOfLines = 0;
            [contentStack insertArrangedSubview:subtitleLabel atIndex:0];
            objc_setAssociatedObject(box, kSubtitleLabelKey, subtitleLabel, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
        }
    } else if (subtitleLabel) {
        [contentStack removeArrangedSubview:subtitleLabel];
        [subtitleLabel removeFromSuperview];
        objc_setAssociatedObject(box, kSubtitleLabelKey, nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}

void JVCardSetChildren(void* handle, void** children, int count) {
    NSBox *box = (__bridge NSBox*)handle;
    NSStackView *contentStack = objc_getAssociatedObject(box, kContentStackKey);
    if (!contentStack) return;

    // Remove existing arranged subviews
    NSArray<NSView*> *existing = [contentStack.arrangedSubviews copy];
    for (NSView *v in existing) {
        [contentStack removeArrangedSubview:v];
        [v removeFromSuperview];
    }

    // Re-add subtitle if it exists
    NSTextField *subtitleLabel = objc_getAssociatedObject(box, kSubtitleLabelKey);
    if (subtitleLabel) {
        [contentStack addArrangedSubview:subtitleLabel];
    }

    // Add new children
    for (int i = 0; i < count; i++) {
        NSView *child = (__bridge NSView*)children[i];
        [contentStack addArrangedSubview:child];
    }
}
