package picogo

/*
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

//export picogoCallback
func picogoCallback(ptr unsafe.Pointer, audio unsafe.Pointer, audioBytes C.int, final bool) bool {
	return getctx(ptr).callback(C.GoBytes(audio, audioBytes), final)
}

var userLock sync.Mutex
var userData = make(map[uintptr]*ctx)
var userPtr uintptr // XXX

//SpeakCallback receives PCM audio chunks as they are being produced.
type SpeakCallback func(pcm []byte, final bool) bool

type ctx struct {
	e        *engine
	callback SpeakCallback
	ptr      uintptr
}

func (c *ctx) release() {
	userLock.Lock()
	defer userLock.Unlock()
	delete(userData, c.ptr)
}

func newctx(e *engine, cb SpeakCallback) *ctx {
	userLock.Lock()
	defer userLock.Unlock()
	userPtr++
	c := &ctx{e, cb, userPtr}
	userData[c.ptr] = c
	return c
}

func getctx(ptr unsafe.Pointer) *ctx {
	userLock.Lock()
	defer userLock.Unlock()
	return userData[uintptr(ptr)]
}
