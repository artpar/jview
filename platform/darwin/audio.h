#ifndef JVIEW_AUDIO_H
#define JVIEW_AUDIO_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateAudio(const char* src, bool autoplay, bool loop, uint64_t endedCbID);
void JVUpdateAudio(void* handle, const char* src, bool loop);

#endif
