package memory

import (
	"fmt"
	"runtime"
)

// https://github.com/gopxl/beep/issues/179 此测试代码块对于简单的内存泄漏检测很有帮助, 之前曾借助其定位过beep的内存泄漏问题。
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))             // 当前堆上分配的内存（MiB）。表示程序当前正在使用的堆内存量。
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)) // 累计分配的内存总量（MiB）。这个值只增不减，但如果程序持续运行且没有内存泄漏，它的增长速度应该会变慢（稳定）。
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))               // 从系统获得的总内存（MiB）。这个值包括堆、栈和其他系统分配的内存。
	fmt.Printf("\tNumGC = %v\n", m.NumGC)                    // 垃圾回收的次数。通过这个可以看GC是否在运行。
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
