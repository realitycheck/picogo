package picogo

/*
#cgo CFLAGS: -I./picopi/pico/lib -I./picopi/pico/tts
#cgo linux LDFLAGS: -lm

#include <stdlib.h>

#include <tts_engine.c>
#include <langfiles.c>

bool cgo_speak(void *user, uint32_t rate, uint32_t format, int channels, uint8_t *audio, uint32_t audio_bytes, bool final) {
	bool speak(void*, uint8_t*, uint32_t, bool); // callback.go:speak
	return speak(user, audio, audio_bytes, final);
}
*/
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

const (
	//LangDefault is default engine's language.
	LangDefault = "en-GB"

	//LangDirDefault is default directory that should contains language files.
	LangDirDefault = "/usr/share/pico/lang/"

	//RateDefault is default speech rate value.
	RateDefault = 100

	//RateMin is min speech rate value.
	RateMin = 20

	//RateMax is max speech rate value.
	RateMax = 500

	//PitchDefault is default speech pitch value.
	PitchDefault = 100

	//PitchMin is default min pitch value.
	PitchMin = 50

	//PitchMax is default max pitch value.
	PitchMax = 200

	//VolumeDefault is default speech volume value.
	VolumeDefault = 100

	//VolumeMin is min speech volume value.
	VolumeMin = 0

	//VolumeMax is max speech volume value.
	VolumeMax = 500

	//AudioRate is audio sample rate constant.
	AudioRate = 16000

	//AudioDepth is audio sample format depth constant.
	AudioDepth = 16

	//AudioChannels is audio channels number constant.
	AudioChannels = 1
)

var (
	//ErrBadLanguage represents not supported language value error.
	ErrBadLanguage = errors.New("language is not supported")

	//ErrBadDirectory represents not existed languages directory error.
	ErrBadDirectory = errors.New("directory is not exists")

	//ErrCreate represents create engine failure.
	ErrCreate = errors.New("failed to create engine")

	//ErrSpeak represents speech generation failure.
	ErrSpeak = errors.New("speak failure")

	supportedLangs map[string]struct{} = map[string]struct{}{
		"en-GB": struct{}{},
		"en-US": struct{}{},
		"de-DE": struct{}{},
		"es-ES": struct{}{},
		"fr-FR": struct{}{},
		"it-IT": struct{}{},
	}
)

//New returns PicoTTS engine.
func New(lang, dir string) (*Engine, error) {
	if _, ok := supportedLangs[lang]; !ok {
		return nil, fmt.Errorf("%s: %w", lang, ErrBadLanguage)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("%s: %w", dir, ErrBadDirectory)
	}

	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	clang := C.CString(lang)
	defer C.free(unsafe.Pointer(clang))

	e := &Engine{
		tts: C.TtsEngine_Create(cdir, clang, C.tts_callback_t(C.cgo_speak)),
	}

	if e.tts == nil {
		return nil, ErrCreate
	}

	runtime.SetFinalizer(&e.tts, func(tts **C.TTS_Engine) {
		C.TtsEngine_Destroy(*tts)
	})

	return e, nil
}

//Callback receives PCM audio chunks as they are being produced.
type Callback func(pcm []byte, final bool) bool

//Engine provides PicoTTS bindings.
type Engine struct {
	tts *C.TTS_Engine
}

//Speak produces PCM audio output of text value.
func (e *Engine) Speak(text string) ([]byte, error) {
	var b bytes.Buffer

	err := e.SpeakCB(text, func(pcm []byte, final bool) bool {
		_, err := b.Write(pcm)
		return err == nil
	})

	return b.Bytes(), err
}

//SpeakCB produces PCM audio output of text value to specified callback function.
func (e *Engine) SpeakCB(text string, cb Callback) error {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	ptr := userDataCreate(cb)
	defer userDataDestroy(ptr)

	if !C.TtsEngine_Speak(e.tts, ctext, unsafe.Pointer(ptr)) {
		return fmt.Errorf("%s: %w", text, ErrSpeak)
	}
	return nil
}

//Abort sends an abort signal to the underlying Speak/SpeakCB audio synth routine.
func (e *Engine) Abort() {
	C.TtsEngine_Stop(e.tts)
}

//Rate gets speech rate.
func (e *Engine) Rate() int {
	return int(C.TtsEngine_GetRate(e.tts))
}

//SetRate sets speech rate.
func (e *Engine) SetRate(rate int) {
	C.TtsEngine_SetRate(e.tts, C.int(rate))
}

//Volume gets speech volume.
func (e *Engine) Volume() int {
	return int(C.TtsEngine_GetVolume(e.tts))
}

//SetVolume sets speech volume.
func (e *Engine) SetVolume(volume int) {
	C.TtsEngine_SetVolume(e.tts, C.int(volume))
}

//Pitch gets speech pitch.
func (e *Engine) Pitch() int {
	return int(C.TtsEngine_GetPitch(e.tts))
}

//SetPitch sets speech pitch.
func (e *Engine) SetPitch(pitch int) {
	C.TtsEngine_SetPitch(e.tts, C.int(pitch))
}
