package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/realitycheck/picogo"
)

var (
	lang = picogo.LangDefault
	dir  = picogo.LangDirDefault

	rate   = picogo.RateDefault
	pitch  = picogo.PitchDefault
	volume = picogo.VolumeDefault

	stdin bool
)

func init() {
	flag.StringVar(&lang, "l", lang, "Language")
	flag.StringVar(&dir, "d", dir, "Languages directory")
	flag.IntVar(&rate, "r", rate, "Speech rate")
	flag.IntVar(&pitch, "p", pitch, "Speech pitch")
	flag.IntVar(&volume, "v", volume, "Speech volume")
	flag.BoolVar(&stdin, "i", stdin, "Use stdin as speech text input instead of TEXT argument")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: picogo [OPTIONS] [TEXT]\n\n")
		fmt.Fprintf(os.Stderr, "picogo command generates text to speech PCM data to stdout using pico tts engine.\n\n")

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

	log.SetOutput(os.Stderr)

	tts, err := picogo.New(lang, dir)
	checkErr(err)

	var text string
	if stdin {
		bytes, err := ioutil.ReadAll(os.Stdin)
		checkErr(err)
		text = string(bytes)
	} else {
		checkErr(fmt.Errorf("text input is not implemented"))
	}

	tts.SetRate(rate)
	tts.SetVolume(volume)
	tts.SetPitch(pitch)
	bytes, err := tts.Speak(text)
	checkErr(err)
	os.Stdout.Write(bytes)
}
