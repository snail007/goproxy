package main

import (
	"C"

	sdk "github.com/snail007/goproxy/sdk/android-ios"
)

//export Start
func Start(serviceID *C.char, serviceArgsStr *C.char) (errStr *C.char) {
	return C.CString(sdk.Start(C.GoString(serviceID), C.GoString(serviceArgsStr)))
}

//export Stop
func Stop(serviceID *C.char) {
	sdk.Stop(C.GoString(serviceID))
}

//export Version
func Version() (ver *C.char) {
	return C.CString(sdk.Version())
}

func main() {
}
