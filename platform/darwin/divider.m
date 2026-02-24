#import <Cocoa/Cocoa.h>
#include "divider.h"

void* JVCreateDivider(void) {
    NSBox *separator = [[NSBox alloc] init];
    separator.boxType = NSBoxSeparator;
    separator.translatesAutoresizingMaskIntoConstraints = NO;
    return (__bridge_retained void*)separator;
}

void JVUpdateDivider(void* handle) {
    // Divider has no dynamic properties to update
}
