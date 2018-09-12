export GOPATH := $(CURDIR)
export LD_LIBRARY_PATH=/Users/break/Documents/Geek/cmx_bot/lib/libmsc.so

bot:
	@echo "Building CmxBot ..."
	go build -o bin/bot bot/intelbot

fbot:
	@echo "Building FirtBot ..."
	go build -o bin/fbot bot/firebot

hbot:
	@echo "Building Hbot ..."
	go build -o bin/hbot bot/hbot

wbot:
	@echo "Building Wbot ..."
	go build -o bin/wbot bot/wikibot

theater:
	@echo "Building Theater ..."
	go build -o bin/theater bot/theater

tts:
	@echo "Building TTS ..."
	go build -o bin/tts bot/tts

mp4:
	@echo "Building Mp4Sever ..."
	go build -o bin/mp4 bot/mp4

deps:
	@echo "Install Installing dependencies"
	@go get -u github.com/golang/dep/cmd/dep
	cd src/bot; ${GOPATH}/bin/dep init; ${GOPATH}/bin/dep ensure -v

get:
	go get github.com/yanyiwu/gojieba