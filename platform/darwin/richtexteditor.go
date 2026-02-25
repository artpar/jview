package darwin

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>
#include <stdint.h>
#include "richtexteditor.h"
*/
import "C"
import (
	"jview/renderer"
	"unsafe"
)

func createRichTextEditorView(node *renderer.RenderNode, surfaceID string) renderer.ViewHandle {
	cContent := C.CString(node.Props.RichContent)
	defer C.free(unsafe.Pointer(cContent))

	var cbID uint64
	if id, ok := node.Callbacks["change"]; ok {
		cbID = uint64(id)
	}

	ptr := C.JVCreateRichTextEditor(cContent, C.bool(node.Props.Editable), C.uint64_t(cbID))
	return renderer.ViewHandle(uintptr(ptr))
}

func updateRichTextEditorView(handle renderer.ViewHandle, node *renderer.RenderNode) {
	cContent := C.CString(node.Props.RichContent)
	defer C.free(unsafe.Pointer(cContent))

	C.JVUpdateRichTextEditor(unsafe.Pointer(handle), cContent, C.bool(node.Props.Editable))
}
