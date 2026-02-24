package engine

import "jview/jlog"

// logWarn emits a structured warning log with component and surface context.
func logWarn(component, surfaceID, msg string) {
	jlog.Warn(component, surfaceID, msg)
}

// logError emits a structured error log with component and surface context.
func logError(component, surfaceID, msg string) {
	jlog.Error(component, surfaceID, msg)
}
