#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>
#include "searchfield.h"

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kSearchFieldDelegateKey = &kSearchFieldDelegateKey;

@interface JVSearchFieldDelegate : NSObject <NSSearchFieldDelegate>
@property (nonatomic, assign) uint64_t callbackID;
@end

@implementation JVSearchFieldDelegate

- (void)controlTextDidChange:(NSNotification *)notification {
    NSSearchField *field = notification.object;
    const char *val = [field.stringValue UTF8String];
    GoCallbackInvoke(self.callbackID, val);
}

// Handle search button click (pressing Enter)
- (void)searchFieldDidStartSearching:(NSSearchField *)sender {
    const char *val = [sender.stringValue UTF8String];
    GoCallbackInvoke(self.callbackID, val);
}

// Handle cancel button click (clear)
- (void)searchFieldDidEndSearching:(NSSearchField *)sender {
    GoCallbackInvoke(self.callbackID, "");
}

@end

void* JVCreateSearchField(const char* placeholder, const char* value, uint64_t callbackID) {
    NSSearchField *field = [[NSSearchField alloc] init];
    field.translatesAutoresizingMaskIntoConstraints = NO;
    field.placeholderString = [NSString stringWithUTF8String:placeholder];
    field.stringValue = [NSString stringWithUTF8String:value];
    field.sendsSearchStringImmediately = YES;
    [field.widthAnchor constraintGreaterThanOrEqualToConstant:100].active = YES;

    JVSearchFieldDelegate *delegate = [[JVSearchFieldDelegate alloc] init];
    delegate.callbackID = callbackID;
    field.delegate = delegate;
    objc_setAssociatedObject(field, kSearchFieldDelegateKey, delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)field;
}

void JVUpdateSearchField(void* handle, const char* placeholder, const char* value) {
    NSSearchField *field = (__bridge NSSearchField*)handle;
    field.placeholderString = [NSString stringWithUTF8String:placeholder];

    // Only update value if not currently being edited (avoid cursor jump)
    if (![field.window.firstResponder isKindOfClass:[NSTextView class]] ||
        ((NSTextView*)field.window.firstResponder).delegate != (id)field) {
        field.stringValue = [NSString stringWithUTF8String:value];
    }
}
