#import <Cocoa/Cocoa.h>
#include "datetimeinput.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kDatePickerTargetKey = &kDatePickerTargetKey;

@interface JVDatePickerTarget : NSObject
@property (nonatomic, assign) uint64_t callbackID;
- (void)dateChanged:(id)sender;
@end

@implementation JVDatePickerTarget

- (void)dateChanged:(id)sender {
    NSDatePicker *picker = (NSDatePicker*)sender;
    NSISO8601DateFormatter *formatter = [[NSISO8601DateFormatter alloc] init];
    formatter.formatOptions = NSISO8601DateFormatWithInternetDateTime;
    NSString *val = [formatter stringFromDate:picker.dateValue];
    GoCallbackInvoke(self.callbackID, [val UTF8String]);
}

@end

static NSDatePickerElementFlags elementsForFlags(bool enableDate, bool enableTime) {
    NSDatePickerElementFlags flags = 0;
    if (enableDate) {
        flags |= NSDatePickerElementFlagYearMonthDay;
    }
    if (enableTime) {
        flags |= NSDatePickerElementFlagHourMinuteSecond;
    }
    if (flags == 0) {
        flags = NSDatePickerElementFlagYearMonthDay;
    }
    return flags;
}

void* JVCreateDateTimeInput(bool enableDate, bool enableTime, const char* value, uint64_t callbackID) {
    NSDatePicker *picker = [[NSDatePicker alloc] init];
    picker.datePickerStyle = NSDatePickerStyleTextFieldAndStepper;
    picker.datePickerElements = elementsForFlags(enableDate, enableTime);
    picker.translatesAutoresizingMaskIntoConstraints = NO;

    // Parse ISO 8601 value if provided
    NSString *valueStr = [NSString stringWithUTF8String:value];
    if (valueStr.length > 0) {
        NSISO8601DateFormatter *formatter = [[NSISO8601DateFormatter alloc] init];
        formatter.formatOptions = NSISO8601DateFormatWithInternetDateTime;
        NSDate *date = [formatter dateFromString:valueStr];
        if (date) {
            picker.dateValue = date;
        }
    }

    JVDatePickerTarget *target = [[JVDatePickerTarget alloc] init];
    target.callbackID = callbackID;
    picker.target = target;
    picker.action = @selector(dateChanged:);

    objc_setAssociatedObject(picker, kDatePickerTargetKey, target, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)picker;
}

void JVUpdateDateTimeInput(void* handle, bool enableDate, bool enableTime, const char* value) {
    NSDatePicker *picker = (__bridge NSDatePicker*)handle;
    picker.datePickerElements = elementsForFlags(enableDate, enableTime);

    NSString *valueStr = [NSString stringWithUTF8String:value];
    if (valueStr.length > 0) {
        NSISO8601DateFormatter *formatter = [[NSISO8601DateFormatter alloc] init];
        formatter.formatOptions = NSISO8601DateFormatWithInternetDateTime;
        NSDate *date = [formatter dateFromString:valueStr];
        if (date) {
            picker.dateValue = date;
        }
    }
}
