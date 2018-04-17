package main

import (
	"C"
	sdk "snail007/proxy/sdk/android-ios"
)

//export Start
func Start(serviceID *C.char,serviceArgsStr *C.char) (errStr *C.char) {
	return C.CString(sdk.Start(C.GoString(serviceID),C.GoString(serviceArgsStr)))
}

//export Stop
func Stop(serviceID *C.char) {
	sdk.Stop(C.GoString(serviceID))
}

//export IsRunning
func IsRunning(serviceID *C.char) C.int {
	if sdk.IsRunning(C.GoString(serviceID)) {
		return 1
	}
	return 0
}
func main() {
}
