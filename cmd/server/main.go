package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/multi-agent-testing/backend/internal/api/router"
	"github.com/multi-agent-testing/backend/internal/config"
	"github.com/multi-agent-testing/backend/pkg/logger"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "configs/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Multi-Agent Testing Platform",
		zap.String("version", "0.1.0"),
		zap.String("mode", cfg.Server.Mode),
	)

	// 3. 初始化数据库 (TODO)
	// initDB(cfg)

	// 4. 初始化aggo客户端 (TODO)
	// initAggoClients(cfg)

	// 5. 创建Hertz服务器
	h := server.Default(
		server.WithHostPorts(cfg.Server.GetAddr()),
	)

	// 6. 注册路由
	router.Setup(h, cfg)

	// 7. 启动服务器
	go func() {
		logger.Info("Server starting", zap.String("addr", cfg.Server.GetAddr()))
		h.Spin()
	}()

	// 8. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := h.Shutdown(context.Background()); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}