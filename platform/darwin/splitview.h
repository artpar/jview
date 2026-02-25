#ifndef JVIEW_SPLITVIEW_H
#define JVIEW_SPLITVIEW_H

#include <stdbool.h>

void* JVCreateSplitView(const char* dividerStyle, bool vertical);
void JVUpdateSplitView(void* handle, const char* dividerStyle, bool vertical, int collapsedPane);
void JVSplitViewSetChildren(void* handle, void** children, int count);

#endif
