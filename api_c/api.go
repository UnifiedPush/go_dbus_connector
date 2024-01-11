package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

typedef void messageCallback(char* instance, uint8_t* message, size_t msglen, char* id);
static void MessageCallback(messageCallback* f, char *a, uint8_t* b, size_t len, char* c) {
	(*f)(a,b,len,c);
	free(a);
	free(b);
	free(c);
}

typedef void endpointCallback(char* instance, char* endpoint);
static void EndpointCallback(endpointCallback* f, char *a, char *b) {
	(*f)(a,b);
	free(a);
	free(b);
}

typedef void unregisteredCallback(char* instance);
static void UnregisteredCallback(unregisteredCallback* f, char *a) {
	(*f)(a);
	free(a);
}

//TODO find better way of syncing up to go defintions
typedef enum {
	UP_REGISTER_STATUS_NEW_ENDPOINT = 0,
	UP_REGISTER_STATUS_REFUSED = 1,
	UP_REGISTER_STATUS_FAILED = 2,
	UP_REGISTER_STATUS_FAILED_OTHER = 99
} UP_REGISTER_STATUS;

*/
import "C"

import (
	"unsafe"

	"unifiedpush.org/go/dbus_connector/api"
)

type Connector struct {
	message      *C.messageCallback
	newEndpoint  *C.endpointCallback
	unregistered *C.unregisteredCallback
}

func (c Connector) Message(a string, b []byte, d string) {
	go C.MessageCallback(c.message, C.CString(a), (*C.uint8_t)(C.CBytes(b)), C.size_t(len(b)), C.CString(d))
}

func (c Connector) NewEndpoint(a, b string) {
	go C.EndpointCallback(c.newEndpoint, C.CString(a), C.CString(b))
}

func (c Connector) Unregistered(a string) {
	go C.UnregisteredCallback(c.unregistered, C.CString(a))
}

/**
 * UPInitializeAndCheck takes in essentailly the same arguments as api.InitializeAndCheck but as typedef'd functions.
 */
//export UPInitializeAndCheck
func UPInitializeAndCheck(
	fullName *C.char,
	friendlyName *C.char,
	msg *C.messageCallback,
	endpoint *C.endpointCallback,
	unregistered *C.unregisteredCallback,
) (ok C.bool) {
	connector := Connector{
		message:      msg,
		newEndpoint:  endpoint,
		unregistered: unregistered,
	}
	err := api.InitializeAndCheck(C.GoString(fullName), C.GoString(friendlyName), connector)
	return err == nil
}

/**
 * UPInitialize takes in essentailly the same arguments as api.Initialize but as typedef'd functions.
 */
//export UPInitialize
func UPInitialize(
	fullName *C.char,
	friendlyName *C.char,
	msg *C.messageCallback,
	endpoint *C.endpointCallback,
	unregistered *C.unregisteredCallback,
) (ok C.bool) {
	connector := Connector{
		message:      msg,
		newEndpoint:  endpoint,
		unregistered: unregistered,
	}
	err := api.Initialize(C.GoString(fullName), C.GoString(friendlyName), connector)
	return err == nil
}

//export UPGetDistributors
func UPGetDistributors() (**C.char, C.size_t) {
	ret, err := api.GetDistributors()
	if err != nil {
		ret = []string{}
	}
	return cStringArray(ret)
}

//CHECKTHIS: TODO
/**
* UPFreeStringArray frees a string array ._. It's meant to be run on the output of UPGetDistributors.
* Primarily a convinience function for bridging to higher level languages.
 */
//export UPFreeStringArray
func UPFreeStringArray(inp **C.char, inp2 C.size_t) {
	// Slice memory layout
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(inp)), int(inp2), int(inp2)}

	// Use unsafe to turn sl into a [] slice.
	b := *(*[]*C.char)(unsafe.Pointer(&sl))
	for _, i := range b {
		UPFreeString(i)
	}

	//also free pointer to memory which contains pointers to strings
	C.free(unsafe.Pointer(inp))
}

//export UPFreeString
func UPFreeString(inp *C.char) {
	C.free(unsafe.Pointer(inp))
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
	return UPRegisterWithDescription(instance, C.CString(""))
}

//export UPRegisterWithDescription
func UPRegisterWithDescription(instance *C.char, description *C.char) (status C.UP_REGISTER_STATUS, reason *C.char) {
	statusret, reasonret, errret := api.RegisterWithDescription(C.GoString(instance), C.GoString(description))
	status = (C.UP_REGISTER_STATUS)(statusret)
	reason = C.CString(reasonret)
	if errret != nil {
		status = C.UP_REGISTER_STATUS_FAILED_OTHER
		reason = C.CString(errret.Error())
	}
	return status, reason
}

//export UPSaveDistributor
func UPSaveDistributor(dist *C.char) (ok C.bool) {
	err := api.SaveDistributor(C.GoString(dist))
	return err == nil
}

//export UPTryUnregister
func UPTryUnregister(instance *C.char) (ok C.bool) {
	err := api.TryUnregister(C.GoString(instance))
	return err == nil
}

//export UPRemoveDistributor
func UPRemoveDistributor() (ok C.bool) {
	err := api.RemoveDistributor()
	return err == nil
}

func main() {}
