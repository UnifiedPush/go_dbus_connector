package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/unifiedpush/go_dbus_connector/api"
)

//temporary connector setup for build testing
type Connector struct {
}

func (c Connector) Message(a, b, d string)  {}
func (c Connector) NewEndpoint(a, b string) {}
func (c Connector) Unregistered(a string)   {}

//export DBusInitialize
func DBusInitialize(str *C.char) {
	connector := Connector{}
	api.Initialize(C.GoString(str), connector)
}

//export ListDistributors
func ListDistributors() **C.char {
	ret, err := api.GetDistributors()
	if err != nil {
		ret = []string{}
	}
	cArray := C.malloc(C.size_t(len(ret)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	a := (*[1<<30 - 1]*C.char)(cArray)
	for idx, substring := range ret {
		a[idx] = C.CString(substring)
	}
	return (**C.char)(cArray)
}

func main() {}
