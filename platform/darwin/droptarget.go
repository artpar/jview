package darwin

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include "droptarget.h"
*/
import "C"
import "unsafe"

// EnableDropTarget enables drag-and-drop on the given view.
func EnableDropTarget(handle uintptr, callbackID uint64) {
	C.JVEnableDropTarget(unsafe.Pointer(handle), C.uint64_t(callbackID))
}

// UpdateDropTargetCallbackID updates the callback ID for an existing drop target.
func UpdateDropTargetCallbackID(handle uintptr, callbackID uint64) {
	C.JVUpdateDropTargetCallbackID(unsafe.Pointer(handle), C.uint64_t(callbackID))
}

// DisableDropTarget removes the drop target from the given view.
func DisableDropTarget(handle uintptr) {
	C.JVDisableDropTarget(unsafe.Pointer(handle))
}
