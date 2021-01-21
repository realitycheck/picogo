package picogo

/*
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

//export speak
func speak(ptr unsafe.Pointer, audio unsafe.Pointer, audioBytes C.int, final bool) bool {
	cb := userDataGet(uintptr(ptr)).(Callback)
	return cb(C.GoBytes(audio, audioBytes), final)
}

var userLock sync.Mutex
var userData = make(map[uintptr]interface{})
var userPtr uintptr

func userDataCreate(obj interface{}) uintptr {
	userLock.Lock()
	defer userLock.Unlock()
	userPtr++
	userData[userPtr] = obj
	return userPtr
}

func userDataDestroy(ptr uintptr) {
	userLock.Lock()
	defer userLock.Unlock()
	delete(userData, ptr)
}

func userDataGet(ptr uintptr) interface{} {
	userLock.Lock()
	defer userLock.Unlock()
	return userData[ptr]
}
