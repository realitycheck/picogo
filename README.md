# picogo

[![GoDoc](https://pkg.go.dev/github.com/realitycheck/picogo/lib?status.svg)](https://pkg.go.dev/github.com/realitycheck/picogo/lib)
[![Go Report Card](https://goreportcard.com/badge/github.com/realitycheck/picogo/lib)](https://goreportcard.com/report/github.com/realitycheck/picogo/lib)

A Go library(CGO) and frontend app to **pico** text to speech C library.

## Usage

### `picogo` frontend app

```bash
# picogo installation
> git clone github.com/realitycheck/picogo
> cd picogo
> make
> sudo make install
> picogo -h
```

```bash
# basic usage
> echo "picogo test message" | picogo -i | aplay --rate=16000 --channels=1 --format=S16_LE
```

### `picogo` library

Prerequisites:

1. [Pico TTS languages](https://github.com/realitycheck/picopi/tree/master/pico/lang/)


```bash
> go get github.com/realitycheck/picogo/lib
```
