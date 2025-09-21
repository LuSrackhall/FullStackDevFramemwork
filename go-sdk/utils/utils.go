package utils

import (
	"os"
	"test-go-sdk/logger"
)

// 检查指定的路径是否存在(TIPS: 实际上这部分逻辑, 应该放在特定的需求函数之内才对)
func CheckAndCreatePathIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 如果路径不存在，创建路径
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger.Error("配置文件路径创建时出错。", "err", err.Error())
		} else {
			logger.Info("配置文件路径创建成功。", "你的配置文件路径为", path)
		}
	} else if err != nil {
		logger.Error("检查配置文件路径时出错。", "err", err.Error())
	} else {
		logger.Info("配置文件路径已存在且无异常。", "你的配置文件路径为", path)
	}
}
