package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	picogo "github.com/realitycheck/picogo/lib"
)

var (
	lang = picogo.LangDefault
	dir  = picogo.LangDirDefault

	rate   = picogo.RateDefault
	pitch  = picogo.PitchDefault
	volume = picogo.VolumeDefault

	stdin bool

	wavenc bool
	input  string
	output string
)

func init() {
	flag.StringVar(&lang, "l", lang, "Language")
	flag.StringVar(&dir, "d", dir, "Languages directory")
	flag.IntVar(&rate, "R", rate, "Speech rate")
	flag.IntVar(&pitch, "P", pitch, "Speech pitch")
	flag.IntVar(&volume, "V", volume, "Speech volume")
	flag.BoolVar(&stdin, "i", stdin, "Use stdin as speech text input instead of TEXT argument")
	flag.BoolVar(&wavenc, "w", wavenc, "Encode PCM output data into Waveform Audio(wav) data")
	flag.StringVar(&output, "o", output, "Use specified file as audio data output instead of STDOUT")
	flag.StringVar(&input, "f", input, "Use specified file as speech text input instead of STDIN or TEXT argument")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: picogo [OPTIONS] [TEXT]\n\n")
		fmt.Fprintf(os.Stderr, "picogo command generates text to speech audio data using pico tts engine.\n\n")

		flag.PrintDefaults()
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("picogo: %v", err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	tts, err := picogo.NewEngine(lang, dir)
	if picogo.ErrBadDirectory == errors.Unwrap(err) {
		err = fmt.Errorf("language files are not found at %s, use -d option", dir)
	}
	checkErr(err)

	var text string
	if input != "" || stdin {
		var bytes []byte
		if input != "" {
			bytes, err = ioutil.ReadFile(input)
		} else {
			bytes, err = ioutil.ReadAll(os.Stdin)
		}
		checkErr(err)
		text = string(bytes)
	} else {
		text = strings.Join(flag.Args(), " ")
	}

	var seeker io.WriteSeeker
	var writer io.Writer
	if output != "" {
		f, err := os.Create(output)
		checkErr(err)
		defer f.Close()
		seeker = f
	} else {
		seeker = os.Stdout
	}

	if wavenc {
		e := wav.NewEncoder(seeker, picogo.AudioRate, picogo.AudioDepth, picogo.AudioChannels, 1 /*PCM*/)
		defer e.Close()
		writer = &wavWriter{e}
	} else {
		writer = seeker
	}

	tts.SetRate(rate)
	tts.SetVolume(volume)
	tts.SetPitch(pitch)
	err = tts.SpeakCB(text, func(pcm []byte, final bool) bool {
		_, err := writer.Write(pcm)
		checkErr(err)
		return true
	})
	checkErr(err)
}

type wavWriter struct {
	e *wav.Encoder
}

func (w *wavWriter) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(b)
	data := make([]int, len(b)/2)
	for i := 0; i < len(data); i++ {
		data[i] = int(int16(binary.LittleEndian.Uint16(buf.Next(2))))
	}
	pcm := audio.IntBuffer{
		Format:         &audio.Format{picogo.AudioChannels, picogo.AudioRate},
		SourceBitDepth: picogo.AudioDepth,
		Data:           data,
	}
	err := w.e.Write(&pcm)
	return w.e.WrittenBytes, err
}
