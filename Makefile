APP_NAME = magpie
VERSION = $(shell git describe --tags --always)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_ARGS = -v -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: install build

install: migrations
	CGO_ENABLED=0 go install $(BUILD_ARGS) .

build: migrations
	CGO_ENABLED=0 go build -o $(APP_NAME) $(BUILD_ARGS) .

migrations:
	go-bindata -pkg steps -o ./migrate/steps/bindata.go migrate/steps 

test:
	go test ./...

