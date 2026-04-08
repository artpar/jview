#import <Cocoa/Cocoa.h>
#import <SystemConfiguration/SystemConfiguration.h>
#include "systemevents.h"

extern void GoSystemEvent(const char* event, const char* data);

// --- Appearance Observer ---

static id appearanceObserver = nil;

void JVStartAppearanceObserver(void) {
    if (appearanceObserver) return;

    appearanceObserver = [[NSDistributedNotificationCenter defaultCenter]
        addObserverForName:@"AppleInterfaceThemeChangedNotification"
                    object:nil
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification * _Nonnull note) {
        NSString *appearance = @"unknown";
        if (@available(macOS 10.14, *)) {
            NSAppearanceName name = [NSApp.effectiveAppearance bestMatchFromAppearancesWithNames:@[
                NSAppearanceNameAqua, NSAppearanceNameDarkAqua
            ]];
            if ([name isEqualToString:NSAppearanceNameDarkAqua]) {
                appearance = @"dark";
            } else {
                appearance = @"light";
            }
        }
        NSString *json = [NSString stringWithFormat:@"{\"appearance\":\"%@\"}", appearance];
        GoSystemEvent("system.appearance", [json UTF8String]);
    }];
}

void JVStopAppearanceObserver(void) {
    if (appearanceObserver) {
        [[NSDistributedNotificationCenter defaultCenter] removeObserver:appearanceObserver];
        appearanceObserver = nil;
    }
}

// --- Power Observer ---

static id powerSleepObserver = nil;
static id powerWakeObserver = nil;

void JVStartPowerObserver(void) {
    if (powerSleepObserver) return;

    NSNotificationCenter *wsnc = [[NSWorkspace sharedWorkspace] notificationCenter];

    powerSleepObserver = [wsnc addObserverForName:NSWorkspaceWillSleepNotification
                                           object:nil
                                            queue:[NSOperationQueue mainQueue]
                                       usingBlock:^(NSNotification * _Nonnull note) {
        GoSystemEvent("system.power.sleep", "{\"state\":\"sleep\"}");
    }];

    powerWakeObserver = [wsnc addObserverForName:NSWorkspaceDidWakeNotification
                                          object:nil
                                           queue:[NSOperationQueue mainQueue]
                                      usingBlock:^(NSNotification * _Nonnull note) {
        GoSystemEvent("system.power.wake", "{\"state\":\"wake\"}");
    }];
}

void JVStopPowerObserver(void) {
    NSNotificationCenter *wsnc = [[NSWorkspace sharedWorkspace] notificationCenter];
    if (powerSleepObserver) {
        [wsnc removeObserver:powerSleepObserver];
        powerSleepObserver = nil;
    }
    if (powerWakeObserver) {
        [wsnc removeObserver:powerWakeObserver];
        powerWakeObserver = nil;
    }
}

// --- Display Observer ---

static id displayObserver = nil;

void JVStartDisplayObserver(void) {
    if (displayObserver) return;

    displayObserver = [[NSNotificationCenter defaultCenter]
        addObserverForName:NSApplicationDidChangeScreenParametersNotification
                    object:nil
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification * _Nonnull note) {
        NSUInteger count = [[NSScreen screens] count];
        NSScreen *main = [NSScreen mainScreen];
        NSRect frame = main.frame;
        NSString *json = [NSString stringWithFormat:
            @"{\"screenCount\":%lu,\"mainWidth\":%.0f,\"mainHeight\":%.0f}",
            (unsigned long)count, frame.size.width, frame.size.height];
        GoSystemEvent("system.display.changed", [json UTF8String]);
    }];
}

void JVStopDisplayObserver(void) {
    if (displayObserver) {
        [[NSNotificationCenter defaultCenter] removeObserver:displayObserver];
        displayObserver = nil;
    }
}

// --- Locale Observer ---

static id localeObserver = nil;

void JVStartLocaleObserver(void) {
    if (localeObserver) return;

    localeObserver = [[NSNotificationCenter defaultCenter]
        addObserverForName:NSCurrentLocaleDidChangeNotification
                    object:nil
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification * _Nonnull note) {
        NSString *identifier = [[NSLocale currentLocale] localeIdentifier];
        NSString *json = [NSString stringWithFormat:@"{\"locale\":\"%@\"}", identifier];
        GoSystemEvent("system.locale.changed", [json UTF8String]);
    }];
}

void JVStopLocaleObserver(void) {
    if (localeObserver) {
        [[NSNotificationCenter defaultCenter] removeObserver:localeObserver];
        localeObserver = nil;
    }
}

// --- Clipboard Observer (polls change count) ---

static dispatch_source_t clipboardTimer = nil;
static long lastClipboardChangeCount = 0;

void JVStartClipboardObserver(int intervalMs) {
    if (clipboardTimer) return;
    if (intervalMs < 100) intervalMs = 100; // minimum 100ms

    lastClipboardChangeCount = [NSPasteboard generalPasteboard].changeCount;

    clipboardTimer = dispatch_source_create(DISPATCH_SOURCE_TYPE_TIMER, 0, 0, dispatch_get_main_queue());
    dispatch_source_set_timer(clipboardTimer,
                              dispatch_time(DISPATCH_TIME_NOW, intervalMs * NSEC_PER_MSEC),
                              intervalMs * NSEC_PER_MSEC,
                              (intervalMs / 10) * NSEC_PER_MSEC);
    dispatch_source_set_event_handler(clipboardTimer, ^{
        long current = [NSPasteboard generalPasteboard].changeCount;
        if (current != lastClipboardChangeCount) {
            lastClipboardChangeCount = current;
            GoSystemEvent("system.clipboard.changed", "{}");
        }
    });
    dispatch_resume(clipboardTimer);
}

void JVStopClipboardObserver(void) {
    if (clipboardTimer) {
        dispatch_source_cancel(clipboardTimer);
        clipboardTimer = nil;
    }
}

// --- Network Reachability ---

static SCNetworkReachabilityRef reachabilityRef = NULL;

static void reachabilityCallback(SCNetworkReachabilityRef target,
                                  SCNetworkReachabilityFlags flags,
                                  void *info) {
    BOOL reachable = (flags & kSCNetworkFlagsReachable) != 0;
    BOOL needsConnection = (flags & kSCNetworkFlagsConnectionRequired) != 0;
    BOOL isWWAN = NO;  // macOS doesn't have WWAN

    NSString *status = @"unreachable";
    NSString *type = @"none";
    if (reachable && !needsConnection) {
        status = @"reachable";
        type = @"wifi"; // macOS: either wifi or ethernet
    }

    NSString *json = [NSString stringWithFormat:
        @"{\"status\":\"%@\",\"type\":\"%@\",\"wwan\":%s}",
        status, type, isWWAN ? "true" : "false"];
    GoSystemEvent("system.network.reachability", [json UTF8String]);
}

void JVStartNetworkObserver(void) {
    if (reachabilityRef) return;

    struct sockaddr_in zeroAddr;
    memset(&zeroAddr, 0, sizeof(zeroAddr));
    zeroAddr.sin_len = sizeof(zeroAddr);
    zeroAddr.sin_family = AF_INET;

    reachabilityRef = SCNetworkReachabilityCreateWithAddress(NULL, (const struct sockaddr *)&zeroAddr);
    if (!reachabilityRef) return;

    SCNetworkReachabilitySetCallback(reachabilityRef, reachabilityCallback, NULL);
    SCNetworkReachabilityScheduleWithRunLoop(reachabilityRef, CFRunLoopGetMain(), kCFRunLoopDefaultMode);
}

void JVStopNetworkObserver(void) {
    if (reachabilityRef) {
        SCNetworkReachabilityUnscheduleFromRunLoop(reachabilityRef, CFRunLoopGetMain(), kCFRunLoopDefaultMode);
        CFRelease(reachabilityRef);
        reachabilityRef = NULL;
    }
}

// --- Accessibility Observer ---

static id accessibilityReduceMotionObserver = nil;
static id accessibilityReduceTransparencyObserver = nil;
static id accessibilityIncreaseContrastObserver = nil;

void JVStartAccessibilityObserver(void) {
    if (accessibilityReduceMotionObserver) return;

    NSWorkspace *ws = [NSWorkspace sharedWorkspace];
    NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];

    accessibilityReduceMotionObserver = [nc
        addObserverForName:NSWorkspaceAccessibilityDisplayOptionsDidChangeNotification
                    object:ws
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification * _Nonnull note) {
        BOOL reduceMotion = [[NSWorkspace sharedWorkspace] accessibilityDisplayShouldReduceMotion];
        BOOL reduceTransparency = [[NSWorkspace sharedWorkspace] accessibilityDisplayShouldReduceTransparency];
        BOOL increaseContrast = [[NSWorkspace sharedWorkspace] accessibilityDisplayShouldIncreaseContrast];

        NSString *json = [NSString stringWithFormat:
            @"{\"reduceMotion\":%s,\"reduceTransparency\":%s,\"increaseContrast\":%s}",
            reduceMotion ? "true" : "false",
            reduceTransparency ? "true" : "false",
            increaseContrast ? "true" : "false"];
        GoSystemEvent("system.accessibility", [json UTF8String]);
    }];
}

void JVStopAccessibilityObserver(void) {
    NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];
    if (accessibilityReduceMotionObserver) {
        [nc removeObserver:accessibilityReduceMotionObserver];
        accessibilityReduceMotionObserver = nil;
    }
}

// --- Thermal State Observer ---

static id thermalObserver = nil;

void JVStartThermalObserver(void) {
    if (thermalObserver) return;

    thermalObserver = [[NSNotificationCenter defaultCenter]
        addObserverForName:NSProcessInfoThermalStateDidChangeNotification
                    object:nil
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification * _Nonnull note) {
        NSString *state = @"nominal";
        switch ([NSProcessInfo processInfo].thermalState) {
            case NSProcessInfoThermalStateNominal: state = @"nominal"; break;
            case NSProcessInfoThermalStateFair: state = @"fair"; break;
            case NSProcessInfoThermalStateSerious: state = @"serious"; break;
            case NSProcessInfoThermalStateCritical: state = @"critical"; break;
        }
        NSString *json = [NSString stringWithFormat:@"{\"state\":\"%@\"}", state];
        GoSystemEvent("system.thermal", [json UTF8String]);
    }];
}

void JVStopThermalObserver(void) {
    if (thermalObserver) {
        [[NSNotificationCenter defaultCenter] removeObserver:thermalObserver];
        thermalObserver = nil;
    }
}
