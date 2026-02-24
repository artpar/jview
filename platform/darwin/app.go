package darwin

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include "app.h"
*/
import "C"

// AppInit initializes NSApplication. Must be called from main thread.
func AppInit() {
	C.JVAppInit()
}

// AppRun starts the NSApplication run loop. Blocks forever. Must be on main thread.
func AppRun() {
	C.JVAppRun()
}
