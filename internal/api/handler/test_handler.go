package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/multi-agent-testing/backend/internal/model"
	"github.com/multi-agent-testing/backend/internal/service"
	"github.com/multi-agent-testing/backend/pkg/logger"
	"go.uber.org/zap"
)

// TestHandler 测试处理器
type TestHandler struct {
	service *service.MultiModelService
}

// NewTestHandler 创建测试处理器
func NewTestHandler(service *service.MultiModelService) *TestHandler {
	return &TestHandler{
		service: service,
	}
}

// ExecuteTest 执行多模型测试
func (h *TestHandler) ExecuteTest(ctx context.Context, c *app.RequestContext) {
	var req model.TestRequest

	// 绑定请求参数
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "Invalid request body"))
		return
	}

	// 验证请求
	if req.Prompts.User == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "User prompt is required"))
		return
	}

	if len(req.Models) == 0 {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "At least one model is required"))
		return
	}

	logger.Info("Received test request",
		zap.Int("model_count", len(req.Models)),
		//zap.String("user_prompt", req.Prompts.User),
	)

	// 执行测试
	result, err := h.service.ExecuteTest(ctx, &req)
	if err != nil {
		logger.Error("Test execution failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, err.Error()))
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

// GetModelList 获取可用模型列表
func (h *TestHandler) GetModelList(ctx context.Context, c *app.RequestContext) {
	models := h.service.GetAvailableModels()

	response := model.ModelListResponse{
		Models: models,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}
