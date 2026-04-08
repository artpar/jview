#ifndef JVIEW_WINDOWDELEGATE_H
#define JVIEW_WINDOWDELEGATE_H

#include <stdint.h>

// Install an NSWindowDelegate on the window for the given surface.
// The delegate forwards window events to Go via GoWindowEvent.
void JVInstallWindowDelegate(const char* surfaceID);

// Remove the window delegate (called before window destruction).
void JVRemoveWindowDelegate(const char* surfaceID);

#endif
