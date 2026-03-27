#pragma once
#include <stdint.h>

// Enable drop target on a view. Callback invoked with JSON: {"paths":["..."],"text":"..."}
void JVEnableDropTarget(void* handle, uint64_t callbackID);

// Update the callback ID for an existing drop target.
void JVUpdateDropTargetCallbackID(void* handle, uint64_t callbackID);

// Disable drop target on a view.
void JVDisableDropTarget(void* handle);
