.PHONY: all clean build run

all: build/game-of-life

build/game-of-life:
	go build -o build/game-of-life $(shell find . -name '*.go')

run: build
	./game-of-life

clean:
	go clean
	rm -rf build