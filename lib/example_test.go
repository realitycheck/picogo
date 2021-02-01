package picogo_test

import (
	"errors"
	"fmt"

	picogo "github.com/realitycheck/picogo/lib"
)

func ExampleEngine_Speak() {
	tts, err := picogo.NewEngine(picogo.LangDefault, picogo.LangDirDefault)
	if picogo.ErrBadDirectory == errors.Unwrap(err) {
		panic(err)
	}
	pcmBytes, err := tts.Speak("This is a text to speech test sample")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", pcmBytes)
}
