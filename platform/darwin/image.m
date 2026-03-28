#import <Cocoa/Cocoa.h>
#include "image.h"
#import <objc/runtime.h>

static const void *kImageURLKey = &kImageURLKey;
static const void *kWidthConstraintKey = &kWidthConstraintKey;
static const void *kHeightConstraintKey = &kHeightConstraintKey;

static void loadImageFromURL(NSImageView *imageView, NSString *urlStr) {
    NSURL *url = [NSURL URLWithString:urlStr];
    if (!url) return;

    // Store the URL to avoid reloading the same image
    NSString *currentURL = objc_getAssociatedObject(imageView, kImageURLKey);
    if ([currentURL isEqualToString:urlStr]) return;
    objc_setAssociatedObject(imageView, kImageURLKey, urlStr, OBJC_ASSOCIATION_RETAIN_NONATOMIC);

    NSURLSession *session = [NSURLSession sharedSession];
    NSURLSessionDataTask *task = [session dataTaskWithURL:url completionHandler:^(NSData *data, NSURLResponse *response, NSError *error) {
        if (error || !data) return;
        NSImage *image = [[NSImage alloc] initWithData:data];
        if (!image) return;
        dispatch_async(dispatch_get_main_queue(), ^{
            // Verify URL hasn't changed while loading
            NSString *latestURL = objc_getAssociatedObject(imageView, kImageURLKey);
            if ([latestURL isEqualToString:urlStr]) {
                imageView.image = image;
            }
        });
    }];
    [task resume];
}

void* JVCreateImage(const char* src, const char* alt, int width, int height) {
    NSImageView *imageView = [[NSImageView alloc] init];
    imageView.translatesAutoresizingMaskIntoConstraints = NO;
    imageView.imageScaling = NSImageScaleProportionallyUpOrDown;

    if (alt) {
        imageView.toolTip = [NSString stringWithUTF8String:alt];
    }

    if (width > 0) {
        NSLayoutConstraint *wc = [imageView.widthAnchor constraintEqualToConstant:width];
        wc.active = YES;
        objc_setAssociatedObject(imageView, kWidthConstraintKey, wc, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }
    if (height > 0) {
        NSLayoutConstraint *hc = [imageView.heightAnchor constraintEqualToConstant:height];
        hc.active = YES;
        objc_setAssociatedObject(imageView, kHeightConstraintKey, hc, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    }

    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadImageFromURL(imageView, srcStr);
    }

    return (__bridge_retained void*)imageView;
}

void JVUpdateImage(void* handle, const char* src, const char* alt, int width, int height) {
    if (!handle) return;
    NSImageView *imageView = (__bridge NSImageView*)handle;

    if (alt) {
        imageView.toolTip = [NSString stringWithUTF8String:alt];
    }

    // Update size constraints
    NSLayoutConstraint *wc = objc_getAssociatedObject(imageView, kWidthConstraintKey);
    if (wc && width > 0) {
        wc.constant = width;
    }
    NSLayoutConstraint *hc = objc_getAssociatedObject(imageView, kHeightConstraintKey);
    if (hc && height > 0) {
        hc.constant = height;
    }

    NSString *srcStr = [NSString stringWithUTF8String:src];
    if (srcStr.length > 0) {
        loadImageFromURL(imageView, srcStr);
    }
}
