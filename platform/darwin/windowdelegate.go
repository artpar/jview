package darwin

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>
#include "windowdelegate.h"
*/
import "C"
import "unsafe"

// installWindowDelegate sets up an NSWindowDelegate on the window for the given surface.
// Window events are forwarded to Go via the GoWindowEvent export.
func installWindowDelegate(surfaceID string) {
	cSID := C.CString(surfaceID)
	defer C.free(unsafe.Pointer(cSID))
	C.JVInstallWindowDelegate(cSID)
}

// removeWindowDelegate removes the window delegate before destruction.
func removeWindowDelegate(surfaceID string) {
	cSID := C.CString(surfaceID)
	defer C.free(unsafe.Pointer(cSID))
	C.JVRemoveWindowDelegate(cSID)
}

// windowEventHandler is set by the engine to receive window events.
// Signature: func(surfaceID, event, data string)
var windowEventHandler func(string, string, string)

// SetWindowEventHandler sets the callback for window events from the native layer.
func SetWindowEventHandler(fn func(surfaceID, event, data string)) {
	windowEventHandler = fn
}

//export GoWindowEvent
func GoWindowEvent(surfaceID *C.char, event *C.char, data *C.char) {
	sid := C.GoString(surfaceID)
	evt := C.GoString(event)
	d := C.GoString(data)
	if windowEventHandler != nil {
		go windowEventHandler(sid, evt, d)
	}
}
