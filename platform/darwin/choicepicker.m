#import <Cocoa/Cocoa.h>
#include "choicepicker.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kPickerTargetKey = &kPickerTargetKey;
static const void *kPickerValuesKey = &kPickerValuesKey;

@interface JVPickerTarget : NSObject
@property (nonatomic, assign) uint64_t callbackID;
- (void)selectionChanged:(id)sender;
@end

@implementation JVPickerTarget

- (void)selectionChanged:(id)sender {
    NSPopUpButton *popup = (NSPopUpButton*)sender;
    NSArray<NSString*> *values = objc_getAssociatedObject(popup, kPickerValuesKey);
    NSInteger idx = popup.indexOfSelectedItem;
    if (idx >= 0 && idx < (NSInteger)values.count) {
        NSString *val = values[idx];
        GoCallbackInvoke(self.callbackID, [val UTF8String]);
    }
}

@end

void* JVCreateChoicePicker(const char** labels, const char** values, int count,
                            const char* selected, uint64_t callbackID) {
    NSPopUpButton *popup = [[NSPopUpButton alloc] initWithFrame:NSZeroRect pullsDown:NO];
    popup.translatesAutoresizingMaskIntoConstraints = NO;

    NSMutableArray<NSString*> *valueArray = [NSMutableArray arrayWithCapacity:count];

    for (int i = 0; i < count; i++) {
        NSString *label = [NSString stringWithUTF8String:labels[i]];
        NSString *value = [NSString stringWithUTF8String:values[i]];
        [popup addItemWithTitle:label];
        [valueArray addObject:value];
    }

    objc_setAssociatedObject(popup, kPickerValuesKey, valueArray, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Select initial value
    NSString *selectedStr = [NSString stringWithUTF8String:selected];
    if (selectedStr.length > 0) {
        for (NSUInteger i = 0; i < valueArray.count; i++) {
            if ([valueArray[i] isEqualToString:selectedStr]) {
                [popup selectItemAtIndex:i];
                break;
            }
        }
    }

    JVPickerTarget *target = [[JVPickerTarget alloc] init];
    target.callbackID = callbackID;
    popup.target = target;
    popup.action = @selector(selectionChanged:);

    objc_setAssociatedObject(popup, kPickerTargetKey, target, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)popup;
}

void JVUpdateChoicePicker(void* handle, const char** labels, const char** values, int count,
                           const char* selected) {
    if (!handle) return;
    NSPopUpButton *popup = (__bridge NSPopUpButton*)handle;
    [popup removeAllItems];

    NSMutableArray<NSString*> *valueArray = [NSMutableArray arrayWithCapacity:count];

    for (int i = 0; i < count; i++) {
        NSString *label = [NSString stringWithUTF8String:labels[i]];
        NSString *value = [NSString stringWithUTF8String:values[i]];
        [popup addItemWithTitle:label];
        [valueArray addObject:value];
    }

    objc_setAssociatedObject(popup, kPickerValuesKey, valueArray, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    NSString *selectedStr = [NSString stringWithUTF8String:selected];
    if (selectedStr.length > 0) {
        for (NSUInteger i = 0; i < valueArray.count; i++) {
            if ([valueArray[i] isEqualToString:selectedStr]) {
                [popup selectItemAtIndex:i];
                break;
            }
        }
    }
}
