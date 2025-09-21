package main

import (
	"flag"
	"log/slog"
	"os"
	"test-go-sdk/logger"
	"test-go-sdk/memory"
	"test-go-sdk/server"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 定义日志目录的路径的命令行参数
var LogDirPath string

// 定义数据库目录路径的命令行参数
var DataBaseDirPath string

func init() {
	go func() {
		for {
			memory.PrintMemUsage()
			time.Sleep(5 * time.Second)
		}
	}()

	// 设置环境变量
	err := godotenv.Load()
	// 在 InitLogger 之前, 使用slog的log信息(与我们的logger模块无关), 此类信息仅在终端上看即可, 输出到日志中也无意义。
	if err != nil {
		// 没必要因为这个err退出程序, .env文件在本项目中, 主要作为开发文件使用。(后面真要上配置文件的话, 也是使用.json格式的)
		slog.Warn("无法加载.env文件", "err", err)
	} else {
		slog.Info(".env文件已被正确加载", "SDK_MODE", os.Getenv("SDK_MODE"))
	}

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

	// 检查所有配置路径是否存在
	{
		// 检查指定的路径是否存在(TIPS: 实际上这部分逻辑, 应该放在特定的需求函数之内才对)
		// utils.CheckAndCreatePathIfNotExist(LogPathAndName)//这个不需要创建路径, 因为是文件路径。TIPS: utils中的函数需要依赖logger, 因此Logger内部不能使用utils中的函数。
		// utils.CheckAndCreatePathIfNotExist(DataBasePath) //ERROR: 由于utils中的内容需要依赖logger, 而logger必须在执行InitLogger之后才能使用, 故此处不能使用utils中的函数。
	}

	// 初始化模块
	{

		// 初始化日志模块(并顺便初始化gin的MODE), 主要是为了输出到日志中, 便于在用户使用过程中记录bug数据。
		{
			logger.InitLogger(LogDirPath)

			// 设置日志级别(此处主要用于开发过程中, 自己可随时进行调整的级别设置)
			logger.ProgramLevel.Set(slog.LevelDebug)
			// logger.ProgramLevel.Set(slog.LevelInfo)
			// logger.ProgramLevel.Set(slog.LevelWarn)
			// logger.ProgramLevel.Set(slog.LevelError)

			if os.Getenv("SDK_MODE") != "debug" {
				// 设置log库, 在正式release中的默认级别
				logger.ProgramLevel.Set(slog.LevelInfo) // 设置级别为Info, 仅展示Info、Warn、Error级别的日志。不展示Debug级别的日志。

				// 设置 gin 框架, 在正式release中的 MODE 为 "release"
				gin.SetMode(gin.ReleaseMode)
			}

			logger.Info("日志模块已开始正常运行, Getenv值已获取。 ", "SDK_MODE", os.Getenv("SDK_MODE"), "GIN_MODE", os.Getenv("GIN_MODE"))
		}

		// 初始化数据库模块
		{

			// TODO: 使用数据库文件路径的逻辑。
			// ...
		}
	}

}

func main() {
	server.ServerRun()
	logger.Debug("testDebug")
	logger.Info("testInfo")
	logger.Warn("testWarn")
	logger.Error("testError")
	time.Sleep(100 * time.Hour)
}
