#import <Cocoa/Cocoa.h>
#include "icon.h"

void* JVCreateIcon(const char* name, int size) {
    NSString *nameStr = [NSString stringWithUTF8String:name];

    NSImageView *imageView = [[NSImageView alloc] init];
    imageView.translatesAutoresizingMaskIntoConstraints = NO;

    if (@available(macOS 11.0, *)) {
        NSImage *image = [NSImage imageWithSystemSymbolName:nameStr accessibilityDescription:nil];
        if (image) {
            NSImageSymbolConfiguration *config = [NSImageSymbolConfiguration configurationWithPointSize:size weight:NSFontWeightRegular];
            imageView.image = [image imageWithSymbolConfiguration:config];
        }
    }

    imageView.imageScaling = NSImageScaleProportionallyUpOrDown;
    [imageView.widthAnchor constraintEqualToConstant:size].active = YES;
    [imageView.heightAnchor constraintEqualToConstant:size].active = YES;

    return (__bridge_retained void*)imageView;
}

void JVUpdateIcon(void* handle, const char* name, int size) {
    if (!handle) return;
    NSImageView *imageView = (__bridge NSImageView*)handle;
    NSString *nameStr = [NSString stringWithUTF8String:name];

    if (@available(macOS 11.0, *)) {
        NSImage *image = [NSImage imageWithSystemSymbolName:nameStr accessibilityDescription:nil];
        if (image) {
            NSImageSymbolConfiguration *config = [NSImageSymbolConfiguration configurationWithPointSize:size weight:NSFontWeightRegular];
            imageView.image = [image imageWithSymbolConfiguration:config];
        }
    }
}
