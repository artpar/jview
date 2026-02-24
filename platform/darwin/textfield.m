#import <Cocoa/Cocoa.h>
#include "textfield.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kTextFieldCallbackKey = &kTextFieldCallbackKey;
static const void *kTextFieldDelegateKey = &kTextFieldDelegateKey;

@interface JVTextFieldDelegate : NSObject <NSTextFieldDelegate>
@property (nonatomic, assign) uint64_t callbackID;
@end

@implementation JVTextFieldDelegate

- (void)controlTextDidChange:(NSNotification *)notification {
    NSTextField *field = notification.object;
    const char *val = [field.stringValue UTF8String];
    GoCallbackInvoke(self.callbackID, val);
}

@end

void* JVCreateTextField(const char* placeholder, const char* value,
                         const char* inputType, bool readOnly, uint64_t callbackID) {
    NSString *placeholderStr = [NSString stringWithUTF8String:placeholder];
    NSString *valueStr = [NSString stringWithUTF8String:value];
    NSString *inputTypeStr = [NSString stringWithUTF8String:inputType];

    NSTextField *field;

    if ([inputTypeStr isEqualToString:@"obscured"]) {
        NSSecureTextField *secureField = [[NSSecureTextField alloc] init];
        field = secureField;
    } else {
        field = [[NSTextField alloc] init];
    }

    field.placeholderString = placeholderStr;
    field.stringValue = valueStr;
    field.editable = !readOnly;
    field.bezeled = YES;
    field.bezelStyle = NSTextFieldRoundedBezel;
    field.translatesAutoresizingMaskIntoConstraints = NO;
    [field.widthAnchor constraintGreaterThanOrEqualToConstant:200].active = YES;

    // Set up delegate for change notifications
    JVTextFieldDelegate *delegate = [[JVTextFieldDelegate alloc] init];
    delegate.callbackID = callbackID;
    field.delegate = delegate;

    // Prevent delegate from being deallocated
    objc_setAssociatedObject(field, kTextFieldDelegateKey, delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)field;
}

void JVUpdateTextField(void* handle, const char* placeholder, const char* value,
                        const char* inputType, bool readOnly) {
    NSTextField *field = (__bridge NSTextField*)handle;
    field.placeholderString = [NSString stringWithUTF8String:placeholder];
    field.stringValue = [NSString stringWithUTF8String:value];
    field.editable = !readOnly;
}
