package picogo

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

//export speakCallback
func speakCallback(ptr unsafe.Pointer, rate uint32, format uint32, channels int, audio unsafe.Pointer, audioBytes uint32, final bool) bool {
	return getctx(ptr).speak(C.GoBytes(audio, C.int(audioBytes)), final)
}

var userLock sync.Mutex
var userData = make(map[unsafe.Pointer]*ctx)

type ctx struct {
	e     *engine
	speak SpeakCallback
}

func (c *ctx) ptr() unsafe.Pointer {
	userLock.Lock()
	defer userLock.Unlock()
	ptr := unsafe.Pointer(C.CString(fmt.Sprintf("%p", c)))
	userData[ptr] = c
	return ptr
}

func freectx(ptr unsafe.Pointer) {
	userLock.Lock()
	defer userLock.Unlock()
	defer C.free(ptr)
	delete(userData, ptr)
}

func getctx(ptr unsafe.Pointer) *ctx {
	userLock.Lock()
	defer userLock.Unlock()
	return userData[ptr]
}
