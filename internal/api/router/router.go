package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/multi-agent-testing/backend/internal/api/handler"
	"github.com/multi-agent-testing/backend/internal/config"
	"github.com/multi-agent-testing/backend/internal/service"
	"github.com/multi-agent-testing/backend/pkg/logger"
	"go.uber.org/zap"
)

// CORS 中间件
func CORS() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}

// Setup 设置路由
func Setup(h *server.Hertz, cfg *config.Config) {
	// 添加全局CORS中间件
	h.Use(CORS())

	// 初始化服务
	multiModelService := service.NewMultiModelService(cfg)

	// 初始化处理器
	testHandler := handler.NewTestHandler(multiModelService)

	// API分组
	api := h.Group("/api/v1")

	// 健康检查
	api.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// 测试相关路由
	testGroup := api.Group("/test")
	{
		testGroup.POST("/execute", testHandler.ExecuteTest)
		// TODO: 实现流式响应
		// testGroup.GET("/stream", testHandler.StreamTest)
	}

	// 模型相关路由
	modelGroup := api.Group("/models")
	{
		modelGroup.GET("/list", testHandler.GetModelList)
		// TODO: 实现模型配置相关接口
		// modelGroup.POST("/config", modelHandler.SaveConfig)
	}

	// TODO: 提示词模板相关路由
	// templateGroup := api.Group("/prompt")
	// {
	// 	templateGroup.GET("/templates", promptHandler.GetTemplates)
	// 	templateGroup.POST("/templates", promptHandler.SaveTemplate)
	// }

	// TODO: 历史记录相关路由
	// historyGroup := api.Group("/history")
	// {
	// 	historyGroup.GET("", historyHandler.GetHistory)
	// }

	logger.Info("Routes registered successfully",
		zap.Int("route_count", 3),
	)
}