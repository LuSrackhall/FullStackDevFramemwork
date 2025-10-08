package server

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/logger"
	"github.com/LuSrackhall/FullStackDevFramemwork/go-sdk/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ServerRun() int {

	// 启动gin
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	mainRouters(r)

	r.GET("/stream", func(c *gin.Context) {

		logger.Logger.Debug("新生成了一个线程............................")
		logger.Debug("新生成了一个线程............................")

		// clientStoresChan := make(chan *config.Store)
		// clientAudioPackageStoresChan := make(chan *audioPackageConfig.Store)
		// clientKeyEventStoresChan := make(chan *keyEvent.Store)
		clientStoresChan := make(chan *store.Store)

		// serverStoresChan := make(chan bool, 1)
		// serverAudioPackageStoresChan := make(chan bool, 1)
		// serverKeyEventStoresChan := make(chan bool, 1)
		serverStoresChan := make(chan bool, 1)

		// config.Clients_sse_stores.Store(clientStoresChan, serverStoresChan)
		// audioPackageConfig.Clients_sse_stores.Store(clientAudioPackageStoresChan, serverAudioPackageStoresChan)
		// keyEvent.Clients_sse_stores.Store(clientKeyEventStoresChan, serverKeyEventStoresChan)

		defer func() {
			// config.Clients_sse_stores.Delete(clientStoresChan)
			// audioPackageConfig.Clients_sse_stores.Delete(clientAudioPackageStoresChan)
			// keyEvent.Clients_sse_stores.Delete(clientKeyEventStoresChan)
			store.Clients_sse_stores.Delete(clientStoresChan)

			logger.Logger.Debug("一个线程退出了............................")
			logger.Debug("一个线程退出了............................")
		}()

		clientGone := c.Request.Context().Done()

		for {
			re := c.Stream(func(w io.Writer) bool {
				select {
				case <-clientGone:
					// serverStoresChan <- false
					// serverAudioPackageStoresChan <- false
					// serverKeyEventStoresChan <- false
					serverStoresChan <- false

					return false

				// case message, ok := <-clientStoresChan:
				// 	if !ok {
				// 		logger.Error("通道clientStoresChan非正常关闭")
				// 		return true
				// 	}
				// 	c.SSEvent("message", message)
				// 	return true
				// case messageAudioPackage, ok := <-clientAudioPackageStoresChan:
				// 	if !ok {
				// 		logger.Error("通道clientAudioPackageStoresChan非正常关闭")
				// 		return true
				// 	}
				// 	c.SSEvent("messageAudioPackage", messageAudioPackage)
				// 	return true
				// case messageKeyEvent, ok := <-clientKeyEventStoresChan:
				// 	if !ok {
				// 		logger.Error("通道clientKeyEventStoresChan非正常关闭")
				// 		return true
				// 	}
				// 	c.SSEvent("messageKeyEvent", messageKeyEvent)
				// 	return true
				case message, ok := <-clientStoresChan:
					if !ok {
						logger.Error("通道clientStoresChan非正常关闭")
						return true
					}
					c.SSEvent("message", message)
					return true

				}
			})

			if !re {
				return
			}
		}

	})

	{ // 保证端口被占用时也能够正常启动服务
		// 尝试在指定端口启动服务
		listener, err := net.Listen("tcp", "localhost:38888")
		if err != nil {
			// 如果38888被占用，让系统分配一个可用端口
			listener, err = net.Listen("tcp", "localhost:0")
			if err != nil {
				logger.Error("无法启动服务:", err)
				return -1
			}
		}

		// 获取实际使用的端口
		port := listener.Addr().(*net.TCPAddr).Port

		// 创建一个channel用于服务器就绪通知
		ready := make(chan bool, 1)

		// 使用listener启动服务
		go func() {
			// 启动服务器
			go func() {
				// time.Sleep(10000 * time.Millisecond)
				if err := r.RunListener(listener); err != nil {
					logger.Error("服务器启动失败:", err)
					ready <- false
				}
			}()

			for {
				// 这里我们没有设置超时限制, 所以会一直阻塞等待所请求的相关服务的返回信息, 直到返回成功或失败, 并解除阻塞。(未设置超时限制, 不会返回超时信息)(由于不做超时限制, 会节省一些损耗, 也能第够一时间作出响应)
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", port))
				// 如果请求没有出错测继续。(否则会开启一轮新的请求->这里我不确定err的可能情况, 因为err不为nil的几率小到可以忽略不计->不知道之前设置的超时限制触发的是不是这里的err, 几率小到离谱, 懒得测试了, 到此为止)。
				if err == nil {
					resp.Body.Close()
					// 如果请求成功则向通道发送true, 以向终端输出端口号信息。(如果失败则不做任何处理, 让本grouting自行结束即可)(如果失败的话, 相关服务启动失败后就会向通道发送false了->以向终端输出服务启动失败的相关信息, 故此处无需处理。)
					if resp.StatusCode == 200 {
						ready <- true
						// fmt.Println("55555555555555666666666666666666666666666666777")
						return
					}
					// fmt.Println("55555555555555")
					// 只要请求本身没有出错, 就退出循环, 不再进行重新请求。
					return
				}
				// fmt.Println("55555555555555666666666666666666666666666666")
			}
		}()

		// 等待服务器就绪信号
		isReady := <-ready
		if !isReady {
			fmt.Println("SDK的本地server模块启动失败")
			return -1
		}
		// 输出端口信息，让Electron主进程可以捕获
		fmt.Printf("KEYTONE_PORT=%d\n", port)
		return port
	}

}

func mainRouters(r *gin.Engine) {
	settingStoreRouters := r.Group("/store")

	// 给到'客户端'或'前端'使用, 供它们获取持久化的设置。
	settingStoreRouters.GET("/get", func(ctx *gin.Context) {

		// key := ctx.Query("key")
		key := ctx.DefaultQuery("key", "unknown")

		if key == "unknown" || key == "" {
			ctx.JSON(200, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:",
			})
			return
		}

		// value := config.GetValue(key)

		// fmt.Println("查询到的value= ", value)

		// ctx.JSON(200, gin.H{
		// 	"message": "ok",
		// 	"key":     key,
		// 	// 这里的value, 会自动转换为JSON字符串
		// 	"value": value,
		// })
	})

	// 给到'前端'使用, 供其ui界面实现应用的设置功能
	settingStoreRouters.POST("/set", func(ctx *gin.Context) {
		type SettingStore struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		}

		var store_setting SettingStore
		err := ctx.ShouldBind(&store_setting)
		if err != nil || store_setting.Key == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		// config.SetValue(store_setting.Key, store_setting.Value)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

}
