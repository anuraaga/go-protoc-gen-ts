package main

import (
	"github.com/wasilibs/go-protoc-gen-ts/internal/runner"
	"github.com/wasilibs/go-protoc-gen-ts/internal/wasm"
)

func main() {
	runner.Run("protoc-gen-ts", wasm.ProtocGenTs)
}
