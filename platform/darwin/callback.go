package darwin

/*
#include <stdint.h>
*/
import "C"

//export GoCallbackInvoke
func GoCallbackInvoke(callbackID C.uint64_t, data *C.char) {
	id := uint64(callbackID)
	d := C.GoString(data)
	globalRegistry.Invoke(id, d)
}
