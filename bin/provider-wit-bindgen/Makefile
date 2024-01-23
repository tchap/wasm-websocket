.PHONY: build

build:
	go build

gen: build
	./provider-with-bindgen -wit ../../wit -out ../../gen github.com/tchap/wasmcloud-websocket
