export GODEBUG = cgocheck=2

TARGET = picogo

# PICO_SRCS is taken from picopi/pico/lib/Makefile:SRCS
PICO_SRCS := picoacph.c \
	picoapi.c \
	picobase.c \
	picocep.c \
	picoctrl.c \
	picodata.c \
	picodbg.c \
	picoextapi.c \
	picofftsg.c \
	picokdbg.c \
	picokdt.c \
	picokfst.c \
	picoklex.c \
	picoknow.c \
	picokpdf.c \
	picokpr.c \
	picoktab.c \
	picoos.c \
	picopal.c \
	picopam.c \
	picopr.c \
	picorsrc.c \
	picosa.c \
	picosig.c \
	picosig2.c \
	picospho.c \
	picotok.c \
	picotrns.c \
	picowa.c


build:
	@CGO_ENABLED=1 go build -o ${TARGET} cmd/picogo/main.go

c:
	@for s in $(PICO_SRCS) ; do \
		echo "#include <$$s>" > cgo_$$s ; \
	done

test: LANG=en-GB
test: VOLUME=100
test: PITCH=100
test: RATE=100
test: LANG_DIR=./picopi/pico/lang
test: TEST=echo "this is a test message for text to speech test"
test: TARGET=picogo
test: build
	@ $(TEST) | \
		./picogo -i -d ${LANG_DIR} -r ${RATE} -v ${VOLUME} -p ${PITCH} -l ${LANG}| \
			aplay --rate=16000 --channels=1 --format=S16_LE


CC_RASPI=/opt/cross-pi-gcc/bin/arm-linux-gnueabihf-gcc

raspi-build: ${CC_RASPI}
raspi-build: export CC=${CC_RASPI}
raspi-build: export GOOS=linux
raspi-build: export GOARCH=arm
raspi-build: TARGET=raspi-picogo
raspi-build: build