// +build linux,cgo

package picogo

/*
#cgo CFLAGS: -I./picopi/pico/lib -I./picopi/pico/tts
#cgo LDFLAGS:  -lsvoxpico -lm

#include <stdlib.h>
#include <tts_engine.c>
#include <langfiles.c>

extern bool speakCallback(void *user, uint32_t rate, uint32_t format, int channels, uint8_t *audio, uint32_t audio_bytes, bool final);
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

	//LangDirDefault is default directory that contains pico languages.
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

	//AudioDepth is audio sample format depth constant (S16_LE).
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

//SpeakCallback receives PCM audio chunks as they are being produced.
type SpeakCallback func(pcm []byte, final bool) bool

//Engine interface provides pico's TTS engine bindings.
type Engine interface {
	//Rate gets speech rate.
	Rate() int

	//Volume gets speech volume.
	Volume() int

	//Pitch gets speech pitch.
	Pitch() int

	//SetRate sets speech rate.
	SetRate(int)

	//SetVolume sets speech volume.
	SetVolume(int)

	//SetPitch sets speech pitch.
	SetPitch(int)

	//Stop sends abort signal to Speak() and SpeakWithCallback() running function processes.
	Stop()

	//Speak produces PCM audio output of text value.
	Speak(string) ([]byte, error)

	//SpeakWithCallback produces PCM audio output of text value to specified callback function.
	SpeakWithCallback(string, SpeakCallback) error
}

//EngineOption represents TTS engine configuration parameter.
type EngineOption func(*engine) error

//Lang option sets engine's language.
func Lang(lang string) EngineOption {
	return func(e *engine) error {
		if _, ok := supportedLangs[lang]; !ok {
			return fmt.Errorf("%s: %w", lang, ErrBadLanguage)
		}
		e.lang = lang
		return nil
	}
}

//LangDir option sets engine's languages directory.
func LangDir(dir string) EngineOption {
	return func(e *engine) error {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			return fmt.Errorf("%s: %w", dir, ErrBadDirectory)
		}
		e.dirs = append(e.dirs, dir)
		return nil
	}
}

//New returns pico's TTS engine instance.
func New(opts ...EngineOption) (Engine, error) {
	e := &engine{
		dirs: []string{
			LangDirDefault,
		},
		lang: LangDefault,
	}
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	for _, dir := range e.dirs {
		cdir := C.CString(dir)
		clang := C.CString(e.lang)
		defer C.free(unsafe.Pointer(cdir))
		defer C.free(unsafe.Pointer(clang))

		tts := C.TtsEngine_Create(cdir, clang, C.tts_callback_t(C.speakCallback))
		if tts != nil {
			e.tts = tts
			break
		}
	}
	if e.tts == nil {
		return nil, ErrCreate
	}

	runtime.SetFinalizer(&e.tts, func(tts **C.TTS_Engine) {
		C.TtsEngine_Destroy(*tts)
	})

	return e, nil
}

type engine struct {
	tts  *C.TTS_Engine
	lang string
	dirs []string
}

func (e *engine) Speak(text string) ([]byte, error) {
	var b bytes.Buffer

	err := e.SpeakWithCallback(text, func(pcm []byte, final bool) bool {
		_, err := b.Write(pcm)
		return err == nil
	})

	return b.Bytes(), err
}

func (e *engine) SpeakWithCallback(text string, cb SpeakCallback) error {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	c := &ctx{e, cb}
	ptr := c.ptr()
	defer freectx(ptr)

	if !C.TtsEngine_Speak(e.tts, ctext, ptr) {
		return fmt.Errorf("%s: %w", text, ErrSpeak)
	}
	return nil
}

func (e *engine) Stop() {
	C.TtsEngine_Stop(e.tts)
}

func (e *engine) Rate() int {
	return int(C.TtsEngine_GetRate(e.tts))
}

func (e *engine) SetRate(rate int) {
	C.TtsEngine_SetRate(e.tts, C.int(rate))
}

func (e *engine) Volume() int {
	return int(C.TtsEngine_GetVolume(e.tts))
}

func (e *engine) SetVolume(volume int) {
	C.TtsEngine_SetVolume(e.tts, C.int(volume))
}

func (e *engine) Pitch() int {
	return int(C.TtsEngine_GetPitch(e.tts))
}

func (e *engine) SetPitch(pitch int) {
	C.TtsEngine_SetPitch(e.tts, C.int(pitch))
}
