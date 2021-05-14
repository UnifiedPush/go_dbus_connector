package main

/*
#include <stdlib.h>
#include <stdio.h>

typedef void messageCallback( char*, char*, char*);
static void MessageCallback(messageCallback* f, char *a, char *b, char* c) {
	(*f)(a,b,c);
}

typedef void endpointCallback(char*, char*);
static void EndpointCallback(endpointCallback* f, char *a, char *b) {
	(*f)(a,b);
}

typedef void unregisteredCallback( char*);
static void UnregisteredCallback(unregisteredCallback* f, char *a) {
	(*f)(a);
}

//TODO find better way of syncing up to go defintions
typedef enum {
	UP_REGISTER_STATUS_NEW_ENDPOINT = 0,
	UP_REGISTER_STATUS_REFUSED = 1,
	UP_REGISTER_STATUS_FAILED = 2
} UP_REGISTER_STATUS;

*/
import "C"
import (
	"unsafe"

	"github.com/unifiedpush/go_dbus_connector/api"
)

type Connector struct {
	message      *C.messageCallback
	newEndpoint  *C.endpointCallback
	unregistered *C.unregisteredCallback
}

func (c Connector) Message(a, b, d string) {
	go C.MessageCallback(c.message, C.CString(a), C.CString(b), C.CString(d))
}
func (c Connector) NewEndpoint(a, b string) {
	go C.EndpointCallback(c.newEndpoint, C.CString(a), C.CString(b))
}
func (c Connector) Unregistered(a string) {
	go C.UnregisteredCallback(c.unregistered, C.CString(a))
}

//export UPInitializeAndCheck
func UPInitializeAndCheck(
	name *C.char,
	msg *C.messageCallback,
	endpoint *C.endpointCallback,
	unregistered *C.unregisteredCallback,
) {
	connector := Connector{
		message:      msg,
		newEndpoint:  endpoint,
		unregistered: unregistered,
	}
	api.InitializeAndCheck(C.GoString(name), connector)
}

//export UPInitialize
func UPInitialize(
	name *C.char,
	msg *C.messageCallback,
	endpoint *C.endpointCallback,
	unregistered *C.unregisteredCallback,
) {
	connector := Connector{
		message:      msg,
		newEndpoint:  endpoint,
		unregistered: unregistered,
	}
	api.Initialize(C.GoString(name), connector)
}

//export UPGetDistributors
func UPGetDistributors() (**C.char, C.size_t) {
	ret, err := api.GetDistributors()
	if err != nil {
		ret = []string{}
	}
	return cStringArray(ret)
}

func cStringArray(arr []string) (**C.char, C.size_t) {
	cArray := C.malloc(C.size_t(len(arr)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	a := (*[1<<30 - 1]*C.char)(cArray)
	for idx, substring := range arr {
		a[idx] = C.CString(substring)
	}
	return (**C.char)(cArray), C.size_t(len(arr))
}

//export UPGetDistributor
func UPGetDistributor() *C.char {
	return C.CString(api.GetDistributor())
}

//export UPRegister
func UPRegister(instance *C.char) (status C.UP_REGISTER_STATUS, reason *C.char) {
	statusret, reasonret, errret := api.Register(C.GoString(instance))
	status = (C.UP_REGISTER_STATUS)(statusret)
	reason = C.CString(reasonret)
	if errret != nil {
		status = (C.UP_REGISTER_STATUS)(99)
		reason = C.CString(errret.Error())
	}
	return status, reason
}

//export UPSaveDistributor
func UPSaveDistributor(dist *C.char) {
	_ = api.SaveDistributor(C.GoString(dist))
	//TODO
}

//export UPTryUnregister
func UPTryUnregister(instance *C.char) (err C.uint) {
	return
}

//export UPRemoveDistributor
func UPRemoveDistributor() {

}

func main() {}
