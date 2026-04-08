#import <Cocoa/Cocoa.h>
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
