#ifndef JVIEW_VIDEO_H
#define JVIEW_VIDEO_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateVideo(const char* src, int width, int height, bool autoplay, bool loop, bool controls, bool muted, uint64_t endedCbID);
void JVUpdateVideo(void* handle, const char* src, int width, int height, bool loop, bool controls, bool muted);

#endif
