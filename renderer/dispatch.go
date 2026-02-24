package renderer

// Dispatcher queues work to run on the main thread.
// On macOS, AppKit views must only be touched from the main thread.
type Dispatcher interface {
	// RunOnMain schedules fn to execute on the main (UI) thread.
	// It returns immediately; fn runs asynchronously.
	RunOnMain(fn func())
}
