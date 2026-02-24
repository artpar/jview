#import <Cocoa/Cocoa.h>
#include "dispatch.h"

// Defined in dispatch.go — Go callback
extern void goDispatchCallback(uintptr_t handle);

void JVDispatchMainAsync(uintptr_t handle) {
    dispatch_async(dispatch_get_main_queue(), ^{
        goDispatchCallback(handle);
    });
}
