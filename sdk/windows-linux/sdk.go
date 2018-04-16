package main

import (
	"C"
	sdk "snail007/proxy/sdk/android-ios"
)

//export Start
func Start(argsStr *C.char) (errStr *C.char) {
	return C.CString(sdk.Start(C.GoString(argsStr)))
}

//export Stop
func Stop(service *C.char) {
	sdk.Stop(C.GoString(service))
}

//export IsRunning
func IsRunning(service *C.char) C.int {
	if sdk.IsRunning(C.GoString(service)) {
		return 1
	}
	return 0
}
func main() {
}
