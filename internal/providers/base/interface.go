package base

import (
	"context"

	"github.com/multi-agent-testing/backend/internal/model"
)

// ModelProvider 模型提供者接口
type ModelProvider interface {
	// Name 返回提供者名称
	Name() string

	// Call 调用模型(非流式)
	Call(ctx context.Context, req *model.CallProvidersRequest) (*model.ModelResponse, error)

	// Stream 流式调用模型
	Stream(ctx context.Context, req *model.CallProvidersRequest) (<-chan *model.StreamChunk, error)

	// ValidateConfig 验证配置
	ValidateConfig(config map[string]interface{}) error
}

// ProviderConfig 提供者配置
type ProviderConfig struct {
	ApiKey  string
	BaseURL string
	Timeout int // 超时时间(秒)
}
