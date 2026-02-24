package engine

/*
#cgo LDFLAGS: -ldl -lffi

#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>
#include <ffi/ffi.h>

// Type IDs matching Go constants.
enum {
    TYPE_VOID    = 0,
    TYPE_INT     = 1,
    TYPE_UINT32  = 2,
    TYPE_INT64   = 3,
    TYPE_UINT64  = 4,
    TYPE_FLOAT   = 5,
    TYPE_DOUBLE  = 6,
    TYPE_POINTER = 7,
    TYPE_STRING  = 8,
    TYPE_BOOL    = 9,
};

// prepared_fn holds everything needed to invoke a function via libffi.
typedef struct {
    ffi_cif cif;
    void *fn_ptr;
    ffi_type **arg_types;
    int n_args;
    int ret_type_id;
    int *arg_type_ids;
    int is_variadic;
} prepared_fn;

static ffi_type* type_id_to_ffi(int id) {
    switch (id) {
        case TYPE_VOID:    return &ffi_type_void;
        case TYPE_INT:     return &ffi_type_sint32;
        case TYPE_UINT32:  return &ffi_type_uint32;
        case TYPE_INT64:   return &ffi_type_sint64;
        case TYPE_UINT64:  return &ffi_type_uint64;
        case TYPE_FLOAT:   return &ffi_type_float;
        case TYPE_DOUBLE:  return &ffi_type_double;
        case TYPE_POINTER: return &ffi_type_pointer;
        case TYPE_STRING:  return &ffi_type_pointer; // const char* is a pointer
        case TYPE_BOOL:    return &ffi_type_sint32;   // bool as int
        default:           return &ffi_type_void;
    }
}

// ffi_prepare creates a prepared function descriptor.
// arg_type_ids is an array of int type IDs, length n_args.
// n_fixed is for variadic functions (0 = not variadic).
static prepared_fn* ffi_prepare_fn(void *fn, int ret_id, int *arg_ids, int n_args, int n_fixed) {
    prepared_fn *f = (prepared_fn*)calloc(1, sizeof(prepared_fn));
    f->fn_ptr = fn;
    f->ret_type_id = ret_id;
    f->n_args = n_args;
    f->is_variadic = (n_fixed > 0) ? 1 : 0;

    f->arg_type_ids = NULL;
    if (n_args > 0) {
        f->arg_type_ids = (int*)calloc(n_args, sizeof(int));
        memcpy(f->arg_type_ids, arg_ids, n_args * sizeof(int));
    }

    f->arg_types = (ffi_type**)calloc(n_args, sizeof(ffi_type*));
    for (int i = 0; i < n_args; i++) {
        f->arg_types[i] = type_id_to_ffi(arg_ids[i]);
    }

    ffi_type *ret_type = type_id_to_ffi(ret_id);
    ffi_status status;

    if (n_fixed > 0) {
        status = ffi_prep_cif_var(&f->cif, FFI_DEFAULT_ABI, (unsigned int)n_fixed, (unsigned int)n_args, ret_type, f->arg_types);
    } else {
        status = ffi_prep_cif(&f->cif, FFI_DEFAULT_ABI, (unsigned int)n_args, ret_type, f->arg_types);
    }

    if (status != FFI_OK) {
        free(f->arg_types);
        free(f->arg_type_ids);
        free(f);
        return NULL;
    }
    return f;
}

// ffi_invoke_fn calls the prepared function.
static void ffi_invoke_fn(prepared_fn *f, void **arg_values, void *ret_value) {
    ffi_call(&f->cif, FFI_FN(f->fn_ptr), ret_value, arg_values);
}

static void ffi_free_fn(prepared_fn *f) {
    if (f) {
        free(f->arg_types);
        free(f->arg_type_ids);
        free(f);
    }
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// Type ID constants matching the C enum.
const (
	typeVoid    = 0
	typeInt     = 1
	typeUint32  = 2
	typeInt64   = 3
	typeUint64  = 4
	typeFloat   = 5
	typeDouble  = 6
	typePointer = 7
	typeString  = 8
	typeBool    = 9
)

// typeNameToID maps JSON type strings to type IDs.
func typeNameToID(name string) (int, error) {
	switch name {
	case "void":
		return typeVoid, nil
	case "int":
		return typeInt, nil
	case "uint32":
		return typeUint32, nil
	case "int64":
		return typeInt64, nil
	case "uint64":
		return typeUint64, nil
	case "float":
		return typeFloat, nil
	case "double":
		return typeDouble, nil
	case "pointer":
		return typePointer, nil
	case "string":
		return typeString, nil
	case "bool":
		return typeBool, nil
	default:
		return -1, fmt.Errorf("ffi: unknown type: %s", name)
	}
}

// genericFunc holds a prepared libffi call descriptor.
type genericFunc struct {
	prepared *C.prepared_fn
	lib      *nativeLib
	retType  int
	argTypes []int
}

// nativeLib holds a dlopen'd library handle.
type nativeLib struct {
	handle unsafe.Pointer
	path   string
}

// HandleTable manages opaque pointer handles for the JSON data model.
type HandleTable struct {
	mu      sync.Mutex
	handles map[uint64]unsafe.Pointer
	next    uint64
}

func NewHandleTable() *HandleTable {
	return &HandleTable{
		handles: make(map[uint64]unsafe.Pointer),
		next:    1,
	}
}

func (h *HandleTable) Register(ptr unsafe.Pointer) uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	id := h.next
	h.next++
	h.handles[id] = ptr
	return id
}

func (h *HandleTable) Resolve(id uint64) (unsafe.Pointer, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	ptr, ok := h.handles[id]
	return ptr, ok
}

func (h *HandleTable) Remove(id uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.handles, id)
}

// FFIRegistry manages loaded native libraries and their callable functions.
type FFIRegistry struct {
	mu      sync.RWMutex
	libs    []*nativeLib
	funcs   map[string]*genericFunc // "prefix.name" → callable
	Handles *HandleTable
}

// NewFFIRegistry creates an empty registry.
func NewFFIRegistry() *FFIRegistry {
	return &FFIRegistry{
		funcs:   make(map[string]*genericFunc),
		Handles: NewHandleTable(),
	}
}

// LoadFromConfig loads all libraries and functions from an FFIConfig.
func (r *FFIRegistry) LoadFromConfig(cfg *FFIConfig) error {
	for _, lib := range cfg.Libraries {
		if err := r.LoadLibrary(lib.Path, lib.Prefix, lib.Functions); err != nil {
			return err
		}
	}
	return nil
}

// LoadLibrary opens a dylib and registers each declared function with its type signature.
func (r *FFIRegistry) LoadLibrary(path, prefix string, funcs []FuncConfig) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_NOW)
	if handle == nil {
		errMsg := C.GoString(C.dlerror())
		return fmt.Errorf("ffi: dlopen %s: %s", path, errMsg)
	}

	lib := &nativeLib{handle: handle, path: path}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.libs = append(r.libs, lib)

	for _, fc := range funcs {
		cSym := C.CString(fc.Symbol)
		fnPtr := C.dlsym(handle, cSym)
		C.free(unsafe.Pointer(cSym))

		if fnPtr == nil {
			errMsg := C.GoString(C.dlerror())
			return fmt.Errorf("ffi: dlsym %s in %s: %s", fc.Symbol, path, errMsg)
		}

		// Resolve type IDs
		retTypeID := typeVoid
		if fc.ReturnType != "" {
			var err error
			retTypeID, err = typeNameToID(fc.ReturnType)
			if err != nil {
				return err
			}
		}

		argTypeIDs := make([]int, len(fc.ParamTypes))
		for i, pt := range fc.ParamTypes {
			var err error
			argTypeIDs[i], err = typeNameToID(pt)
			if err != nil {
				return err
			}
		}

		// Prepare the libffi call descriptor
		var cArgIDs *C.int
		if len(argTypeIDs) > 0 {
			cArgIDs = (*C.int)(C.malloc(C.size_t(len(argTypeIDs)) * C.size_t(unsafe.Sizeof(C.int(0)))))
			slice := unsafe.Slice((*C.int)(unsafe.Pointer(cArgIDs)), len(argTypeIDs))
			for i, id := range argTypeIDs {
				slice[i] = C.int(id)
			}
		}

		prepared := C.ffi_prepare_fn(fnPtr, C.int(retTypeID), cArgIDs, C.int(len(argTypeIDs)), C.int(fc.FixedArgs))
		if cArgIDs != nil {
			C.free(unsafe.Pointer(cArgIDs))
		}
		if prepared == nil {
			return fmt.Errorf("ffi: ffi_prep_cif failed for %s", fc.Symbol)
		}

		name := prefix + "." + fc.Name
		r.funcs[name] = &genericFunc{
			prepared: prepared,
			lib:      lib,
			retType:  retTypeID,
			argTypes: argTypeIDs,
		}
	}

	return nil
}

// Has returns true if the named function is registered.
func (r *FFIRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.funcs[name]
	return ok
}

// Call invokes a registered native function with the given Go args.
// Args are converted to the declared C types, the function is called via libffi,
// and the result is converted back to a Go value.
func (r *FFIRegistry) Call(name string, args []interface{}) (interface{}, error) {
	r.mu.RLock()
	fn, ok := r.funcs[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("ffi: unknown function: %s", name)
	}

	if len(args) != len(fn.argTypes) {
		return nil, fmt.Errorf("ffi: %s expects %d args, got %d", name, len(fn.argTypes), len(args))
	}

	// Allocate arg value pointers and the storage they point to.
	// Each element of argPtrs is a pointer to the actual value storage.
	nArgs := len(fn.argTypes)
	var argPtrs **unsafe.Pointer
	if nArgs > 0 {
		argPtrs = (**unsafe.Pointer)(C.calloc(C.size_t(nArgs), C.size_t(unsafe.Sizeof(unsafe.Pointer(nil)))))
		defer C.free(unsafe.Pointer(argPtrs))
	}
	argSlice := unsafe.Slice((*unsafe.Pointer)(unsafe.Pointer(argPtrs)), nArgs)

	// Storage for each arg value — we keep Go references to prevent GC
	cStrings := make([]*C.char, 0) // track C strings to free later
	defer func() {
		for _, cs := range cStrings {
			C.free(unsafe.Pointer(cs))
		}
	}()

	// Storage backing for scalar args
	argStorage := make([]unsafe.Pointer, nArgs)
	defer func() {
		for _, p := range argStorage {
			if p != nil {
				C.free(p)
			}
		}
	}()

	for i, typeID := range fn.argTypes {
		switch typeID {
		case typeInt, typeBool:
			p := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0)))))
			switch v := args[i].(type) {
			case float64:
				*p = C.int(int32(v))
			case bool:
				if v {
					*p = 1
				} else {
					*p = 0
				}
			default:
				return nil, fmt.Errorf("ffi: %s arg %d: expected number or bool, got %T", name, i, args[i])
			}
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeUint32:
			p := (*C.uint)(C.malloc(C.size_t(unsafe.Sizeof(C.uint(0)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected number, got %T", name, i, args[i])
			}
			*p = C.uint(uint32(v))
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeInt64:
			p := (*C.longlong)(C.malloc(C.size_t(unsafe.Sizeof(C.longlong(0)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected number, got %T", name, i, args[i])
			}
			*p = C.longlong(int64(v))
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeUint64:
			p := (*C.ulonglong)(C.malloc(C.size_t(unsafe.Sizeof(C.ulonglong(0)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected number, got %T", name, i, args[i])
			}
			*p = C.ulonglong(uint64(v))
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeFloat:
			p := (*C.float)(C.malloc(C.size_t(unsafe.Sizeof(C.float(0)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected number, got %T", name, i, args[i])
			}
			*p = C.float(float32(v))
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeDouble:
			p := (*C.double)(C.malloc(C.size_t(unsafe.Sizeof(C.double(0)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected number, got %T", name, i, args[i])
			}
			*p = C.double(v)
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typeString:
			s, ok := args[i].(string)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected string, got %T", name, i, args[i])
			}
			cs := C.CString(s)
			cStrings = append(cStrings, cs)
			// For string args, we need a pointer to the char* pointer
			p := (**C.char)(C.malloc(C.size_t(unsafe.Sizeof((*C.char)(nil)))))
			*p = cs
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)

		case typePointer:
			p := (*unsafe.Pointer)(C.malloc(C.size_t(unsafe.Sizeof(unsafe.Pointer(nil)))))
			v, ok := args[i].(float64)
			if !ok {
				return nil, fmt.Errorf("ffi: %s arg %d: expected handle ID (number), got %T", name, i, args[i])
			}
			handleID := uint64(v)
			ptr, found := r.Handles.Resolve(handleID)
			if !found {
				return nil, fmt.Errorf("ffi: %s arg %d: unknown handle %d", name, i, handleID)
			}
			*p = ptr
			argStorage[i] = unsafe.Pointer(p)
			argSlice[i] = unsafe.Pointer(p)
		}
	}

	// Invoke
	var retBuf [8]byte // large enough for any scalar return type
	retPtr := unsafe.Pointer(&retBuf[0])

	C.ffi_invoke_fn(fn.prepared, (*unsafe.Pointer)(unsafe.Pointer(argPtrs)), retPtr)

	// Convert return value
	switch fn.retType {
	case typeVoid:
		return nil, nil

	case typeInt:
		val := *(*C.int)(retPtr)
		return float64(val), nil

	case typeBool:
		val := *(*C.int)(retPtr)
		return val != 0, nil

	case typeUint32:
		val := *(*C.uint)(retPtr)
		return float64(val), nil

	case typeInt64:
		val := *(*C.longlong)(retPtr)
		return float64(val), nil

	case typeUint64:
		val := *(*C.ulonglong)(retPtr)
		return float64(val), nil

	case typeFloat:
		val := *(*C.float)(retPtr)
		return float64(val), nil

	case typeDouble:
		val := *(*C.double)(retPtr)
		return float64(val), nil

	case typePointer:
		val := *(*unsafe.Pointer)(retPtr)
		if val == nil {
			return nil, nil
		}
		handleID := r.Handles.Register(val)
		return float64(handleID), nil

	case typeString:
		val := *(**C.char)(retPtr)
		if val == nil {
			return "", nil
		}
		return C.GoString(val), nil
	}

	return nil, fmt.Errorf("ffi: unsupported return type %d", fn.retType)
}

// Close dlcloses all loaded libraries and frees prepared functions.
func (r *FFIRegistry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, fn := range r.funcs {
		C.ffi_free_fn(fn.prepared)
	}
	for _, lib := range r.libs {
		C.dlclose(lib.handle)
	}
	r.libs = nil
	r.funcs = make(map[string]*genericFunc)
}
