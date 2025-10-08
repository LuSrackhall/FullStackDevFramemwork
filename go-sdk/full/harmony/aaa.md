# 为什么构建失败

```bash
# 关键: 强制使用 Go 版 net
go build --tags "fts5 netgo" -ldflags "-s -w" -buildmode=c-shared -v -o libfull.so .
```

# 为什么 .h 没生成

必须在 Go 代码中 导出至少一个函数 给 C 调用（使用 //export 注释指令）。(注意不能是 // export 注释指令)

```go
package main

/*
#include <stdint.h>
*/
import "C"

//export Add
func Add(a, b C.int) C.int {
  return a + b
}

//export Hello
func Hello() {
  println("Hello from Go (OHOS)")
}

func main() {} // 必须有 main
```