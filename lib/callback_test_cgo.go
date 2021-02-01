package picogo

import "C"

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"
)

func test_CallbackSpeak(t *testing.T) {
	audio := bytes.NewBufferString("this is audio sample").Bytes()

	var fakeCallback Callback = func(pcm []byte, final bool) bool {
		if !reflect.DeepEqual(audio, pcm) {
			t.Errorf("audio != pcm: want %v, got %v", audio, pcm)
		}
		return true
	}
	ptr := userDataCreate(fakeCallback)
	defer userDataDestroy(ptr)

	speak(unsafe.Pointer(ptr), unsafe.Pointer(&audio[0]), C.int(len(audio)), false)
}
