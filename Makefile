export GOPATH := $(CURDIR)

bot:
	@echo "Building CmxBot ..."
	go build -o bin/bot bot/ai_x

deps:
	@echo "Install Installing dependencies"
	@go get -u github.com/golang/dep/cmd/dep
	cd src/bot; ${GOPATH}/bin/dep init; ${GOPATH}/bin/dep ensure -v