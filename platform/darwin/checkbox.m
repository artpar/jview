#import <Cocoa/Cocoa.h>
#include "checkbox.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kCheckBoxCallbackKey = &kCheckBoxCallbackKey;

@interface JVCheckBoxTarget : NSObject
@property (nonatomic, assign) uint64_t callbackID;
- (void)toggled:(id)sender;
@end

@implementation JVCheckBoxTarget

- (void)toggled:(id)sender {
    NSButton *cb = (NSButton*)sender;
    const char *val = (cb.state == NSControlStateValueOn) ? "true" : "false";
    GoCallbackInvoke(self.callbackID, val);
}

@end

void* JVCreateCheckBox(const char* label, bool checked, uint64_t callbackID) {
    NSString *labelStr = [NSString stringWithUTF8String:label];

    NSButton *checkbox = [NSButton checkboxWithTitle:labelStr target:nil action:nil];
    checkbox.state = checked ? NSControlStateValueOn : NSControlStateValueOff;
    checkbox.translatesAutoresizingMaskIntoConstraints = NO;

    JVCheckBoxTarget *target = [[JVCheckBoxTarget alloc] init];
    target.callbackID = callbackID;
    checkbox.target = target;
    checkbox.action = @selector(toggled:);

    objc_setAssociatedObject(checkbox, kCheckBoxCallbackKey, target, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)checkbox;
}

void JVUpdateCheckBox(void* handle, const char* label, bool checked) {
    NSButton *checkbox = (__bridge NSButton*)handle;
    checkbox.title = [NSString stringWithUTF8String:label];
    checkbox.state = checked ? NSControlStateValueOn : NSControlStateValueOff;
}
