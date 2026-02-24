#ifndef JVIEW_IMAGE_H
#define JVIEW_IMAGE_H

void* JVCreateImage(const char* src, const char* alt, int width, int height);
void JVUpdateImage(void* handle, const char* src, const char* alt, int width, int height);

#endif
