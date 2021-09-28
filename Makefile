all: build serve

build:
	GOOS=js GOARCH=wasm go build -o website/bin.wasm

serve:
	go run ./scripts/server