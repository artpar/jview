#import <Cocoa/Cocoa.h>
#include "slider.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kSliderTargetKey = &kSliderTargetKey;

@interface JVSliderTarget : NSObject
@property (nonatomic, assign) uint64_t callbackID;
- (void)sliderChanged:(id)sender;
@end

@implementation JVSliderTarget

- (void)sliderChanged:(id)sender {
    NSSlider *slider = (NSSlider*)sender;
    NSString *val = [NSString stringWithFormat:@"%g", slider.doubleValue];
    GoCallbackInvoke(self.callbackID, [val UTF8String]);
}

@end

void* JVCreateSlider(double min, double max, double step, double value, uint64_t callbackID) {
    NSSlider *slider = [[NSSlider alloc] init];
    slider.minValue = min;
    slider.maxValue = max;
    slider.doubleValue = value;
    slider.translatesAutoresizingMaskIntoConstraints = NO;
    [slider.widthAnchor constraintGreaterThanOrEqualToConstant:200].active = YES;

    if (step > 0) {
        int ticks = (int)((max - min) / step) + 1;
        if (ticks > 1 && ticks <= 100) {
            slider.numberOfTickMarks = ticks;
            slider.allowsTickMarkValuesOnly = YES;
        }
    }

    JVSliderTarget *target = [[JVSliderTarget alloc] init];
    target.callbackID = callbackID;
    slider.target = target;
    slider.action = @selector(sliderChanged:);
    slider.continuous = YES;

    objc_setAssociatedObject(slider, kSliderTargetKey, target, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    return (__bridge_retained void*)slider;
}

void JVUpdateSlider(void* handle, double min, double max, double step, double value) {
    NSSlider *slider = (__bridge NSSlider*)handle;
    slider.minValue = min;
    slider.maxValue = max;
    slider.doubleValue = value;

    if (step > 0) {
        int ticks = (int)((max - min) / step) + 1;
        if (ticks > 1 && ticks <= 100) {
            slider.numberOfTickMarks = ticks;
            slider.allowsTickMarkValuesOnly = YES;
        } else {
            slider.numberOfTickMarks = 0;
            slider.allowsTickMarkValuesOnly = NO;
        }
    }
}
