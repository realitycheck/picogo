# picogo

[![Go Reference](https://pkg.go.dev/badge/github.com/realitycheck/picogo?status.svg.svg)](https://pkg.go.dev/github.com/realitycheck/picogo?status.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/realitycheck/picogo)](https://goreportcard.com/report/github.com/realitycheck/picogo)

A Go library(CGO) and frontend app to **pico** text to speech C library.

## Usage

### `picogo` frontend app

```bash
# picogo installation
> git clone github.com/realitycheck/picogo
> git submodule update --init
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
> go get github.com/realitycheck/picogo
```
