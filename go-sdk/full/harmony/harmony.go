// package main

// import "test-go-sdk/full"

// func FullSdkRun(LogDirPath string, DataBaseDirPath string) int {
// 	return full.FullSdkRun(LogDirPath, DataBaseDirPath)
// }
// func main() {}

package main

import "C"
import (
	"fmt"
	"log/slog"
)

// import "fmt"

//export Add
func Add(x, y C.double) C.double {
	return x + y
}

//export Hello
func Hello() *C.char {
	fmt.Println("LuSrackhall Hello")
	slog.Info("LuSrackhall Hello slog")
	return C.CString("Hello, HarmonyOS")
}
func main() {}
