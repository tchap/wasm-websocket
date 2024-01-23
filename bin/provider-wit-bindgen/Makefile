.PHONY: build

build:
	go build

gen: build
	./provider-with-bindgen -in ../../wit -out ../../gen github.com/tchap/wasmcloud-websocket/gen
