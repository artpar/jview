#import <Cocoa/Cocoa.h>
#import <AVKit/AVKit.h>
#import <AVFoundation/AVFoundation.h>
#include "video.h"
#import <objc/runtime.h>

extern void GoCallbackInvoke(uint64_t callbackID, const char* data);

static const void *kVideoURLKey = &kVideoURLKey;
static const void *kWidthConstraintKey = &kWidthConstraintKey;
static const void *kHeightConstraintKey = &kHeightConstraintKey;
static const void *kLoopKey = &kLoopKey;
static const void *kEndedCbIDKey = &kEndedCbIDKey;
static const void *kEndedObserverKey = &kEndedObserverKey;

static void removeEndedObserver(AVPlayerView *playerView) {
    id observer = objc_getAssociatedObject(playerView, kEndedObserverKey);
    if (observer) {
        [[NSNotificationCenter defaultCenter] removeObserver:observer];
        objc_setAssociatedObject(playerView, kEndedObserverKey, nil, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
}

static void addEndedObserver(AVPlayerView *playerView, AVPlayerItem *item) {
    NSNumber *loopNum = objc_getAssociatedObject(playerView, kLoopKey);
    BOOL loop = [loopNum boolValue];
    NSNumber *cbNum = objc_getAssociatedObject(playerView, kEndedCbIDKey);
    uint64_t endedCbID = [cbNum unsignedLongLongValue];

    id observer = [[NSNotificationCenter defaultCenter]
        addObserverForName:AVPlayerItemDidPlayToEndTimeNotification
                    object:item
                     queue:[NSOperationQueue mainQueue]
                usingBlock:^(NSNotification *note) {
                    if (loop) {
                        [item seekToTime:kCMTimeZero completionHandler:nil];
                        [playerView.player play];
                    } else if (endedCbID != 0) {
                        GoCallbackInvoke(endedCbID, "");
                    }
                }];
    objc_setAssociatedObject(playerView, kEndedObserverKey, observer, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

static void loadVideo(AVPlayerView *playerView, NSString *urlStr, BOOL autoplay) {
    NSURL *url = [NSURL URLWithString:urlStr];
    if (!url) return;

    // Skip reload if URL unchanged
    NSString *currentURL = objc_getAssociatedObject(playerView, kVideoURLKey);
    if ([currentURL isEqualToString:urlStr]) return;
    objc_setAssociatedObject(playerView, kVideoURLKey, urlStr, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Remove old observer
    removeEndedObserver(playerView);

    AVPlayerItem *item = [AVPlayerItem playerItemWithURL:url];
    if (playerView.player) {
        [playerView.player replaceCurrentItemWithPlayerItem:item];
    } else {
        playerView.player = [AVPlayer playerWithPlayerItem:item];
    }

    // Add end-of-playback observer
    addEndedObserver(playerView, item);

    if (autoplay) {
        [playerView.player play];
    }
}

void* JVCreateVideo(const char* src, int width, int height, bool autoplay, bool loop, bool controls, bool muted, uint64_t endedCbID) {
    AVPlayerView *playerView = [[AVPlayerView alloc] init];
    playerView.translatesAutoresizingMaskIntoConstraints = NO;

    // Controls style
    playerView.controlsStyle = controls ? AVPlayerViewControlsStyleDefault : AVPlayerViewControlsStyleNone;

    // Store loop and callback
    objc_setAssociatedObject(playerView, kLoopKey, @(loop), OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    objc_setAssociatedObject(playerView, kEndedCbIDKey, @(endedCbID), OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Size constraints
    if (width > 0) {
        NSLayoutConstraint *wc = [playerView.widthAnchor constraintEqualToConstant:width];
        wc.active = YES;
        objc_setAssociatedObject(playerView, kWidthConstraintKey, wc, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
    if (height > 0) {
        NSLayoutConstraint *hc = [playerView.heightAnchor constraintEqualToConstant:height];
        hc.active = YES;
        objc_setAssociatedObject(playerView, kHeightConstraintKey, hc, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }

    // Load video
    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadVideo(playerView, srcStr, autoplay);
    }

    // Mute
    if (muted && playerView.player) {
        playerView.player.muted = YES;
    }

    return (__bridge_retained void*)playerView;
}

void JVUpdateVideo(void* handle, const char* src, int width, int height, bool loop, bool controls, bool muted) {
    AVPlayerView *playerView = (__bridge AVPlayerView*)handle;

    // Update controls style
    playerView.controlsStyle = controls ? AVPlayerViewControlsStyleDefault : AVPlayerViewControlsStyleNone;

    // Update loop flag
    objc_setAssociatedObject(playerView, kLoopKey, @(loop), OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    // Update size constraints
    NSLayoutConstraint *wc = objc_getAssociatedObject(playerView, kWidthConstraintKey);
    if (wc && width > 0) {
        wc.constant = width;
    }
    NSLayoutConstraint *hc = objc_getAssociatedObject(playerView, kHeightConstraintKey);
    if (hc && height > 0) {
        hc.constant = height;
    }

    // Update muted
    if (playerView.player) {
        playerView.player.muted = muted;
    }

    // Load new video if src changed (no autoplay on update)
    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadVideo(playerView, srcStr, NO);
    }
}
