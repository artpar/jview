package darwin

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework SystemConfiguration

#include <stdlib.h>
#include "systemevents.h"
*/
import "C"

// systemEventHandler is set by the engine to receive system events.
var systemEventHandler func(string, string)

// SetSystemEventHandler sets the callback for system events from the native layer.
func SetSystemEventHandler(fn func(event, data string)) {
	systemEventHandler = fn
}

//export GoSystemEvent
func GoSystemEvent(event *C.char, data *C.char) {
	evt := C.GoString(event)
	d := C.GoString(data)
	if systemEventHandler != nil {
		go systemEventHandler(evt, d)
	}
}

// StartAppearanceObserver begins watching for light/dark mode changes.
func StartAppearanceObserver() { C.JVStartAppearanceObserver() }

// StopAppearanceObserver stops watching appearance changes.
func StopAppearanceObserver() { C.JVStopAppearanceObserver() }

// StartPowerObserver begins watching for sleep/wake events.
func StartPowerObserver() { C.JVStartPowerObserver() }

// StopPowerObserver stops watching power events.
func StopPowerObserver() { C.JVStopPowerObserver() }

// StartDisplayObserver begins watching for display changes.
func StartDisplayObserver() { C.JVStartDisplayObserver() }

// StopDisplayObserver stops watching display changes.
func StopDisplayObserver() { C.JVStopDisplayObserver() }

// StartLocaleObserver begins watching for locale changes.
func StartLocaleObserver() { C.JVStartLocaleObserver() }

// StopLocaleObserver stops watching locale changes.
func StopLocaleObserver() { C.JVStopLocaleObserver() }

// StartClipboardObserver begins polling the clipboard for changes.
func StartClipboardObserver(intervalMs int) { C.JVStartClipboardObserver(C.int(intervalMs)) }

// StopClipboardObserver stops clipboard polling.
func StopClipboardObserver() { C.JVStopClipboardObserver() }

// StartNetworkObserver begins watching for network reachability changes.
func StartNetworkObserver() { C.JVStartNetworkObserver() }

// StopNetworkObserver stops watching network changes.
func StopNetworkObserver() { C.JVStopNetworkObserver() }

// StartAccessibilityObserver begins watching for accessibility setting changes.
func StartAccessibilityObserver() { C.JVStartAccessibilityObserver() }

// StopAccessibilityObserver stops watching accessibility changes.
func StopAccessibilityObserver() { C.JVStopAccessibilityObserver() }

// StartThermalObserver begins watching for thermal state changes.
func StartThermalObserver() { C.JVStartThermalObserver() }

// StopThermalObserver stops watching thermal changes.
func StopThermalObserver() { C.JVStopThermalObserver() }
