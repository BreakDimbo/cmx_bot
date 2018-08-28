export GOPATH := $(CURDIR)

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

deps:
	@echo "Install Installing dependencies"
	@go get -u github.com/golang/dep/cmd/dep
	cd src/bot; ${GOPATH}/bin/dep init; ${GOPATH}/bin/dep ensure -v


get:
	go get github.com/yanyiwu/gojieba