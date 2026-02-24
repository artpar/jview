#ifndef JVIEW_MODAL_H
#define JVIEW_MODAL_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateModal(const char* title, bool visible, const char* surfaceID, int width, int height, uint64_t dismissCbID);
void JVUpdateModal(void* handle, const char* title, bool visible);
void JVModalSetChildren(void* handle, void** children, int count);
void JVCleanupModal(void* handle);

#endif
