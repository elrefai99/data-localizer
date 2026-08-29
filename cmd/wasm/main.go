//go:build js && wasm

package main

import "syscall/js"

func main() {
	bridge := newBridge()
	js.Global().Set("__dataLocalizerModules", bridge.modules())
	select {}
}
