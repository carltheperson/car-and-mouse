all: build serve

build-wasm:
	GOOS=js GOARCH=wasm go build -o website/bin.wasm

build-server:
	go build -o ./tmp/server ./scripts/server

serve:
	go run ./scripts/server

run-dev:
	-test ! -f third_party/air && echo "\nair is not installed. Installing\n" && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && mv bin/* third_party && rm -r bin
	third_party/air