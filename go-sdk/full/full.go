package full

import (
	"log/slog"
	"os"

	"time"

	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/logger"
	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/memory"
	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/server"
	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

}

func FullSdkRun(LogDirPath string, DataBaseDirPath string) int {
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
			// 检查数据库目录的存储路径是否存在
			utils.CheckAndCreatePathIfNotExist(DataBaseDirPath) //TIPS: 由于utils中的内容需要依赖logger, 因此必须在初始化日志模块之后使用。
			// TODO: 使用数据库文件路径的逻辑。
			// ...
		}

	}

	port := server.ServerRun()
	logger.Debug("testDebug")
	logger.Info("testInfo")
	logger.Warn("testWarn")
	logger.Error("testError")

	return port
}
