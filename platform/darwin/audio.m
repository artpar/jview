#import <Cocoa/Cocoa.h>
#import <AVFoundation/AVFoundation.h>
#import <CoreMedia/CoreMedia.h>
#include "audio.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kAudioURLKey = &kAudioURLKey;
static const void *kAudioLoopKey = &kAudioLoopKey;
static const void *kAudioEndedCbIDKey = &kAudioEndedCbIDKey;
static const void *kAudioEndedObserverKey = &kAudioEndedObserverKey;
static const void *kAudioPlayerKey = &kAudioPlayerKey;
static const void *kAudioTimeObserverKey = &kAudioTimeObserverKey;
static const void *kAudioPlayButtonKey = &kAudioPlayButtonKey;
static const void *kAudioScrubberKey = &kAudioScrubberKey;
static const void *kAudioTimeLabelKey = &kAudioTimeLabelKey;
static const void *kAudioScrubActionKey = &kAudioScrubActionKey;

// --- Target-Action helper for scrubber ---
@interface JVAudioScrubAction : NSObject
@property (nonatomic, assign) NSView *container;
- (void)scrubberChanged:(NSSlider *)sender;
@end

@implementation JVAudioScrubAction
- (void)scrubberChanged:(NSSlider *)sender {
    NSView *container = self.container;
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    if (!player || !player.currentItem) return;

    CMTime duration = player.currentItem.duration;
    if (CMTIME_IS_INDEFINITE(duration)) return;

    Float64 totalSeconds = CMTimeGetSeconds(duration);
    Float64 seekTo = sender.doubleValue * totalSeconds;
    [player seekToTime:CMTimeMakeWithSeconds(seekTo, NSEC_PER_SEC) completionHandler:^(BOOL finished) {}];
}
@end

// --- Target-Action helper for play/pause button ---
@interface JVAudioPlayAction : NSObject
@property (nonatomic, assign) NSView *container;
- (void)togglePlay:(NSButton *)sender;
@end

@implementation JVAudioPlayAction
- (void)togglePlay:(NSButton *)sender {
    NSView *container = self.container;
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    if (!player) return;

    if (player.rate > 0) {
        [player pause];
        sender.image = [NSImage imageWithSystemSymbolName:@"play.fill" accessibilityDescription:@"Play"];
    } else {
        [player play];
        sender.image = [NSImage imageWithSystemSymbolName:@"pause.fill" accessibilityDescription:@"Pause"];
    }
}
@end

static NSString* formatTime(Float64 seconds) {
    if (isnan(seconds) || seconds < 0) seconds = 0;
    int mins = (int)seconds / 60;
    int secs = (int)seconds % 60;
    return [NSString stringWithFormat:@"%d:%02d", mins, secs];
}

static void removeEndedObserver(NSView *container) {
    id observer = objc_getAssociatedObject(container, kAudioEndedObserverKey);
    if (observer) {
        [[NSNotificationCenter defaultCenter] removeObserver:observer];
        objc_setAssociatedObject(container, kAudioEndedObserverKey, nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}

static void removeTimeObserver(NSView *container) {
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    id token = objc_getAssociatedObject(container, kAudioTimeObserverKey);
    if (player && token) {
        [player removeTimeObserver:token];
        objc_setAssociatedObject(container, kAudioTimeObserverKey, nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}

static void addEndedObserver(NSView *container, AVPlayerItem *item) {
    NSNumber *loopNum = objc_getAssociatedObject(container, kAudioLoopKey);
    BOOL loop = [loopNum boolValue];
    NSNumber *cbNum = objc_getAssociatedObject(container, kAudioEndedCbIDKey);
    uint64_t endedCbID = [cbNum unsignedLongLongValue];
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    NSButton *playBtn = objc_getAssociatedObject(container, kAudioPlayButtonKey);

    id observer = [[NSNotificationCenter defaultCenter]
        addObserverForName:AVPlayerItemDidPlayToEndTimeNotification
                    object:item
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification *note) {
                    if (loop) {
                        [item seekToTime:kCMTimeZero completionHandler:nil];
                        [player play];
                    } else {
                        playBtn.image = [NSImage imageWithSystemSymbolName:@"play.fill" accessibilityDescription:@"Play"];
                        if (endedCbID != 0) {
                            GoCallbackInvoke(endedCbID, "");
                        }
                    }
                }];
    objc_setAssociatedObject(container, kAudioEndedObserverKey, observer, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

static void addTimeObserver(NSView *container) {
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    NSSlider *scrubber = objc_getAssociatedObject(container, kAudioScrubberKey);
    NSTextField *timeLabel = objc_getAssociatedObject(container, kAudioTimeLabelKey);
    if (!player) return;

    __weak AVPlayer *weakPlayer = player;
    id token = [player addPeriodicTimeObserverForInterval:CMTimeMakeWithSeconds(0.25, NSEC_PER_SEC)
                                                    queue:dispatch_get_main_queue()
                                               usingBlock:^(CMTime time) {
        AVPlayer *strongPlayer = weakPlayer;
        if (!strongPlayer) return;
        Float64 current = CMTimeGetSeconds(time);
        CMTime dur = strongPlayer.currentItem.duration;
        Float64 total = CMTIME_IS_INDEFINITE(dur) ? 0 : CMTimeGetSeconds(dur);

        if (total > 0) {
            scrubber.doubleValue = current / total;
        }
        NSString *text = [NSString stringWithFormat:@"%@ / %@", formatTime(current), formatTime(total)];
        timeLabel.stringValue = text;
    }];
    objc_setAssociatedObject(container, kAudioTimeObserverKey, token, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

static void loadAudio(NSView *container, NSString *urlStr, BOOL autoplay) {
    NSURL *url = [NSURL URLWithString:urlStr];
    if (!url) return;

    // Skip reload if URL unchanged
    NSString *currentURL = objc_getAssociatedObject(container, kAudioURLKey);
    if ([currentURL isEqualToString:urlStr]) return;
    objc_setAssociatedObject(container, kAudioURLKey, urlStr, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Remove old observers
    removeEndedObserver(container);
    removeTimeObserver(container);

    AVPlayerItem *item = [AVPlayerItem playerItemWithURL:url];
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    if (player) {
        [player replaceCurrentItemWithPlayerItem:item];
    } else {
        player = [AVPlayer playerWithPlayerItem:item];
        objc_setAssociatedObject(container, kAudioPlayerKey, player, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }

    // Add observers
    addEndedObserver(container, item);
    addTimeObserver(container);

    NSButton *playBtn = objc_getAssociatedObject(container, kAudioPlayButtonKey);
    if (autoplay) {
        [player play];
        playBtn.image = [NSImage imageWithSystemSymbolName:@"pause.fill" accessibilityDescription:@"Pause"];
    } else {
        playBtn.image = [NSImage imageWithSystemSymbolName:@"play.fill" accessibilityDescription:@"Play"];
    }
}

void* JVCreateAudio(const char* src, bool autoplay, bool loop, uint64_t endedCbID) {
    NSView *container = [[NSView alloc] init];
    container.translatesAutoresizingMaskIntoConstraints = NO;

    // Play/Pause button
    NSButton *playBtn = [NSButton buttonWithImage:[NSImage imageWithSystemSymbolName:@"play.fill" accessibilityDescription:@"Play"]
                                           target:nil
                                           action:nil];
    playBtn.translatesAutoresizingMaskIntoConstraints = NO;
    playBtn.bordered = NO;
    playBtn.bezelStyle = NSBezelStyleInline;
    [playBtn setContentHuggingPriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];

    // Scrubber (progress slider)
    NSSlider *scrubber = [NSSlider sliderWithValue:0 minValue:0 maxValue:1 target:nil action:nil];
    scrubber.translatesAutoresizingMaskIntoConstraints = NO;
    [scrubber setContentHuggingPriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];

    // Time label
    NSTextField *timeLabel = [NSTextField labelWithString:@"0:00 / 0:00"];
    timeLabel.translatesAutoresizingMaskIntoConstraints = NO;
    timeLabel.font = [NSFont monospacedDigitSystemFontOfSize:11 weight:NSFontWeightRegular];
    timeLabel.textColor = [NSColor secondaryLabelColor];
    [timeLabel setContentHuggingPriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];
    [timeLabel setContentCompressionResistancePriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];

    [container addSubview:playBtn];
    [container addSubview:scrubber];
    [container addSubview:timeLabel];

    // Layout: [playBtn]-8-[scrubber]-8-[timeLabel]
    [NSLayoutConstraint activateConstraints:@[
        // Height
        [container.heightAnchor constraintEqualToConstant:40],

        // Play button
        [playBtn.leadingAnchor constraintEqualToAnchor:container.leadingAnchor constant:8],
        [playBtn.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],
        [playBtn.widthAnchor constraintEqualToConstant:24],

        // Scrubber
        [scrubber.leadingAnchor constraintEqualToAnchor:playBtn.trailingAnchor constant:8],
        [scrubber.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],

        // Time label
        [timeLabel.leadingAnchor constraintEqualToAnchor:scrubber.trailingAnchor constant:8],
        [timeLabel.trailingAnchor constraintEqualToAnchor:container.trailingAnchor constant:-8],
        [timeLabel.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],
    ]];

    // Store associated objects
    objc_setAssociatedObject(container, kAudioLoopKey, @(loop), OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(container, kAudioEndedCbIDKey, @(endedCbID), OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(container, kAudioPlayButtonKey, playBtn, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(container, kAudioScrubberKey, scrubber, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(container, kAudioTimeLabelKey, timeLabel, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Wire up play/pause action
    JVAudioPlayAction *playAction = [[JVAudioPlayAction alloc] init];
    playAction.container = container;
    playBtn.target = playAction;
    playBtn.action = @selector(togglePlay:);
    objc_setAssociatedObject(container, @selector(togglePlay:), playAction, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Wire up scrubber action
    JVAudioScrubAction *scrubAction = [[JVAudioScrubAction alloc] init];
    scrubAction.container = container;
    scrubber.target = scrubAction;
    scrubber.action = @selector(scrubberChanged:);
    objc_setAssociatedObject(container, kAudioScrubActionKey, scrubAction, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Load audio
    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadAudio(container, srcStr, autoplay);
    }

    return (__bridge_retained void*)container;
}

void JVUpdateAudio(void* handle, const char* src, bool loop) {
    NSView *container = (__bridge NSView*)handle;

    // Update loop flag
    objc_setAssociatedObject(container, kAudioLoopKey, @(loop), OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Load new audio if src changed (no autoplay on update)
    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadAudio(container, srcStr, NO);
    }
}

void JVCleanupAudio(void* handle) {
    NSView *container = (__bridge NSView*)handle;
    removeTimeObserver(container);
    removeEndedObserver(container);
    AVPlayer *player = objc_getAssociatedObject(container, kAudioPlayerKey);
    if (player) {
        [player pause];
    }
}
