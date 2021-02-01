package picogo

import (
	"reflect"
	"testing"
)

func Test_UserData(t *testing.T) {
	obj := struct{}{}

	if len(userData) != 0 {
		t.Errorf("len(userData) != 0")
	}

	ptr := userDataCreate(obj)
	ptrFake := ptr + 1

	if len(userData) != 1 {
		t.Errorf("userDataCreate: len(userData) != 1")
	}

	u := userDataGet(ptr)
	if !reflect.DeepEqual(obj, u) {
		t.Errorf("userDataGet(ptr): unexpected data: %d", ptr)
	}

	fake := userDataGet(ptrFake)
	if fake != nil {
		t.Errorf("userDataGet(ptrFake): unexpected data: %v", fake)
	}

	userDataDestroy(ptr)
	if len(userData) != 0 {
		t.Errorf("userDataDestroy: len(userData) != 0")
	}

	userDataDestroy(ptrFake) // no panic
}

func Test_CallbackSpeak(t *testing.T) {
	test_CallbackSpeak(t)
}
