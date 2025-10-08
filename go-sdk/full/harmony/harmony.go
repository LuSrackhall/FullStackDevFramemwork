package main

import "C"

import "github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/full"

// func FullSdkRun(LogDirPath, DataBaseDirPath string) int {
// 	return full.FullSdkRun(LogDirPath, DataBaseDirPath)
// }

//export FullSdkRun
func FullSdkRun(logDirPath *C.char, dataBaseDirPath *C.char) C.int {
	goLogDir := C.GoString(logDirPath)
	goDBDir := C.GoString(dataBaseDirPath)

	result := full.FullSdkRun(goLogDir, goDBDir)

	return C.int(result)
}

func main() {}
