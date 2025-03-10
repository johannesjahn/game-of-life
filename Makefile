.PHONY: all clean build run test

all: build/game-of-life

build/game-of-life:
	go build -o build/game-of-life cmd/game-of-life/main.go

run: build
	go run cmd/game-of-life/main.go

test:
	go test -v ./...

clean:
	go clean
	rm -rf build