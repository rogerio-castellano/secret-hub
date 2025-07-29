APP_NAME = secret-hub
OS := $(shell go env GOOS)
EXT :=

ifeq ($(OS),windows)
    EXT := .exe
endif

build:
	go build -o $(APP_NAME)$(EXT) ./main.go

test:
	go test ./...


sync:
	git pull
	git push
