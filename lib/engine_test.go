package picogo

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

const (
	langDirInternal = "../picopi/pico/lang/"
)

func Test_New(t *testing.T) {

	tests := []struct {
		name      string
		lang, dir string
		wantErr   error
	}{
		{
			name:    "Check no language directory",
			lang:    LangDefault,
			dir:     "./not-existed",
			wantErr: ErrBadDirectory,
		},
		{
			name:    "Check en-GB is supported",
			lang:    "en-GB",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check en-US is supported",
			lang:    "en-US",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check de-DE is supported",
			lang:    "de-DE",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check es-ES is supported",
			lang:    "es-ES",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check fr-FR is supported",
			lang:    "fr-FR",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check it-IT is supported",
			lang:    "it-IT",
			dir:     langDirInternal,
			wantErr: nil,
		},
		{
			name:    "Check zu-ZU is not supported",
			lang:    "zu-ZU",
			dir:     langDirInternal,
			wantErr: ErrBadLanguage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEngine(tt.lang, tt.dir)
			if err != nil {
				if tt.wantErr != errors.Unwrap(err) {
					t.Errorf("want %v error, but got %v", tt.wantErr, err)
				}
			} else if e == nil {
				t.Errorf("e == nil")
			} else if e.tts == nil {
				t.Errorf("e.tts == nil")
			}
		})
	}

	// wait finalizer to be called
	runtime.GC()
	time.Sleep(400 * time.Millisecond)
}

func Test_Speak(t *testing.T) {
	tests := []struct {
		lang string
		text string
		want int
	}{
		{
			lang: "en-US",
			text: "Hello",
			want: 20224,
		},
		{
			lang: "en-GB",
			text: "Hello",
			want: 19968,
		},
		{
			lang: "en-GB",
			text: ".............................................................",
			want: 1085312,
		},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			e := newTestEngine(tt.lang)
			pcm, err := e.Speak(tt.text)
			if err != nil {
				t.Fatalf("error Speak: %s", err)
			}
			if len(pcm) != tt.want {
				t.Fatalf("want %d, got %d length", tt.want, len(pcm))
			}
		})
	}
}

func newTestEngine(lang string) *Engine {
	e, err := NewEngine(lang, langDirInternal)
	if err != nil {
		panic(err)
	}
	return e
}

func Test_Rate(t *testing.T) {
	e := newTestEngine("en-US")
	if e.Rate() != RateDefault {
		t.Errorf("e.Rate() != RateDefault: %d", e.Rate())
	}
	e.SetRate(RateMin)
	if e.Rate() != RateMin {
		t.Errorf("e.Rate() != RateMin: %d", e.Rate())
	}
	e.SetRate(RateMax)
	if e.Rate() != RateMax {
		t.Errorf("e.Rate() != RateMax: %d", e.Rate())
	}
}

func Test_Pitch(t *testing.T) {
	e := newTestEngine("en-US")
	if e.Pitch() != PitchDefault {
		t.Errorf("e.Pitch() != PitchDefault: %d", e.Pitch())
	}
	e.SetPitch(PitchMin)
	if e.Pitch() != PitchMin {
		t.Errorf("e.Pitch() != PitchMin: %d", e.Pitch())
	}
	e.SetPitch(PitchMax)
	if e.Pitch() != PitchMax {
		t.Errorf("e.Pitch() != PitchMax: %d", e.Pitch())
	}
}

func Test_Volume(t *testing.T) {
	e := newTestEngine("en-US")
	if e.Volume() != VolumeDefault {
		t.Errorf("e.Volume() != VolumeDefault: %d", e.Volume())
	}
	e.SetVolume(VolumeMin)
	if e.Volume() != VolumeMin {
		t.Errorf("e.Volume() != VolumeMin: %d", e.Volume())
	}
	e.SetVolume(VolumeMax)
	if e.Volume() != VolumeMax {
		t.Errorf("e.Volume() != VolumeMax: %d", e.Volume())
	}
}

func Test_AbortSpeakCB(t *testing.T) {
	e := newTestEngine("it-IT")

	// text must be large enough to callback three times:
	// 1st call: abort synth
	// 2nd call: receive remaining buffer
	// 3d call: not happen
	text := "Hello World, abort the engine, text must be larger enough to callback three times"
	var i int

	err := e.SpeakCB(text, func(pcm []byte, final bool) bool {
		i++
		return true
	})

	if i < 3 {
		t.Fatalf("text is not large enough to continue: %d", i)
	}
	if err != nil {
		t.Fatalf("[1] e.SpeakCB: %v", err)
	}

	i = 0
	err = e.SpeakCB(text, func(pcm []byte, final bool) bool {
		i++
		return false
	})
	if i != 2 {
		t.Errorf("e.SpeakCB is not stopped by callback return value: %d", i)
	} else if errors.Unwrap(err) != ErrSpeak {
		t.Errorf("e.SpeakCB stopped by callback returns unexpected error value: %v", err)
	}

	i = 0
	err = e.SpeakCB(text, func(pcm []byte, final bool) bool {
		i++
		e.Abort()
		return true
	})
	if i != 2 {
		t.Errorf("e.SpeakCB is not stopped by Abort() call: %d", i)
	} else if errors.Unwrap(err) != ErrSpeak {
		t.Errorf("e.SpeakCB stopped by Abort() call returns unexpected error value: %v", err)
	}
}
