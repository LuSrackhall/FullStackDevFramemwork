package main

import (
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/full"
)

// 定义日志目录的路径的命令行参数
var LogDirPath string

// 定义数据库目录路径的命令行参数
var DataBaseDirPath string

func init() {
	// 获取命令行参数中的传入值
	{

		// 如果路径不存在, 则使用当前目录作为路径
		// * 第一个参数是指向一个字符串变量的指针，用于存储解析后的值。
		// * 第二个参数是命令行参数的名称（在命令行中使用）。  用户在使用时 go run main.go -configPath=./path
		// * 第三个参数是默认值（如果用户没有提供这个参数，则使用默认值）。
		// * 第四个参数是这个参数的描述（帮助信息）。
		flag.StringVar(&LogDirPath, "LogDirPath", "./", "用于定义日志存放目录的路径")
		flag.StringVar(&DataBaseDirPath, "DataBaseDirPath", "./", "用于定义数据库存放目录的路径")

		// 解析命令行参数
		flag.Parse()

		// 使用命令行参数
		// ...
		slog.Info("命行参数已正确解析", "LogDirPath", LogDirPath, "DataBaseDirPath", DataBaseDirPath)
	}
}

func main() {
	port := full.FullSdkRun(LogDirPath, DataBaseDirPath)

	fmt.Println("[main]: SDK服务已启动, 监听端口:", port)
	time.Sleep(100 * time.Second)
}
