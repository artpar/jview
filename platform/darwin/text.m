#import <Cocoa/Cocoa.h>
#include "text.h"

static NSFont* fontForVariant(NSString *variant) {
    if ([variant isEqualToString:@"h1"]) return [NSFont boldSystemFontOfSize:28];
    if ([variant isEqualToString:@"h2"]) return [NSFont boldSystemFontOfSize:22];
    if ([variant isEqualToString:@"h3"]) return [NSFont boldSystemFontOfSize:18];
    if ([variant isEqualToString:@"h4"]) return [NSFont boldSystemFontOfSize:16];
    if ([variant isEqualToString:@"h5"]) return [NSFont boldSystemFontOfSize:14];
    if ([variant isEqualToString:@"caption"]) return [NSFont systemFontOfSize:11];
    // body (default)
    return [NSFont systemFontOfSize:13];
}

void* JVCreateText(const char* content, const char* variant) {
    NSString *text = [NSString stringWithUTF8String:content];
    NSString *var_ = [NSString stringWithUTF8String:variant];

    NSTextField *label = [NSTextField labelWithString:text];
    label.font = fontForVariant(var_);
    label.lineBreakMode = NSLineBreakByWordWrapping;
    label.maximumNumberOfLines = 0; // unlimited lines
    label.translatesAutoresizingMaskIntoConstraints = NO;
    [label setContentHuggingPriority:NSLayoutPriorityDefaultHigh forOrientation:NSLayoutConstraintOrientationVertical];

    return (__bridge_retained void*)label;
}

void JVUpdateText(void* handle, const char* content, const char* variant) {
    NSTextField *label = (__bridge NSTextField*)handle;
    label.stringValue = [NSString stringWithUTF8String:content];
    label.font = fontForVariant([NSString stringWithUTF8String:variant]);
}
