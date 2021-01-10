# picogo

Go cgo bindings and frontend application for Pico TTS.

## CGO build requirements

* [Pico TTS library](https://github.com/realitycheck/picopi/tree/master/pico/lib/) headers.
* [Pico TTS tts engine](https://github.com/realitycheck/picopi/tree/master/pico/tts/) sources.

## CGO runtime requirements

* [Pico TTS library](https://github.com/realitycheck/picopi/tree/master/pico/lib/) installed.
* [Pico TTS languages](https://github.com/realitycheck/picopi/tree/master/pico/lang/) installed.


## HOW-TO

### Use `picogo` frontend app

```
> go install github.com/realitycheck/picogo/cmd/picogo
> picogo -h
```