package zhipu

import (
	"context"
	"errors"
	"time"

	"github.com/CoolBanHub/aggo/agent"
	"github.com/CoolBanHub/aggo/model"
	"github.com/cloudwego/eino/schema"
	internalModel "github.com/multi-agent-testing/backend/internal/model"
	"github.com/multi-agent-testing/backend/internal/providers/base"
	"github.com/multi-agent-testing/backend/pkg/logger"
	"go.uber.org/zap"
)

// Provider zhipu模型提供者
type Provider struct {
	config base.ProviderConfig
}

// NewProvider 创建zhipu提供者
func NewProvider(config base.ProviderConfig) *Provider {
	return &Provider{
		config: config,
	}
}

// Name 返回提供者名称
func (p *Provider) Name() string {
	return "zhipu"
}

// Call 调用zhipu模型(非流式)
func (p *Provider) Call(ctx context.Context, req *internalModel.CallProvidersRequest) (*internalModel.ModelResponse, error) {
	startTime := time.Now()

	logger.Info("Calling zhipu model with aggo",
		zap.String("provider", p.Name()),
		zap.String("base_url", p.config.BaseURL),
	)

	// 创建聊天模型
	cm, err := model.NewChatModel(
		model.WithBaseUrl(p.config.BaseURL),
		model.WithAPIKey(p.config.ApiKey),
		model.WithModel(req.Models.Name), // 使用配置中的模型
	)
	if err != nil {
		logger.Error("Failed to create chat model",
			zap.String("provider", p.Name()),
			zap.Error(err),
		)
		return &internalModel.ModelResponse{
			ModelName:    req.Models.Name,
			Provider:     p.Name(),
			Content:      "",
			Success:      false,
			Error:        err.Error(),
			ResponseTime: time.Since(startTime).Milliseconds(),
			StartTime:    startTime,
			EndTime:      time.Now(),
		}, err
	}

	// 创建代理 (不使用记忆管理)
	var systemPrompt string
	if req.Prompts.System != "" {
		systemPrompt = req.Prompts.System
	}

	ag, err := agent.NewAgent(ctx, cm,
		agent.WithSystemPrompt(systemPrompt),
	)
	if err != nil {
		logger.Error("Failed to create agent",
			zap.String("provider", p.Name()),
			zap.Error(err),
		)
		return &internalModel.ModelResponse{
			ModelName:    "gpt-5-mini",
			Provider:     p.Name(),
			Content:      "",
			Success:      false,
			Error:        err.Error(),
			ResponseTime: time.Since(startTime).Milliseconds(),
			StartTime:    startTime,
			EndTime:      time.Now(),
		}, err
	}

	// 构建消息列表
	messages := []*schema.Message{}
	for _, _v := range req.Prompts.Message {
		v := _v
		if v.Role == "user" {
			messages = append(messages, schema.UserMessage(v.Content))
		} else if v.Role == "assistant" {
			messages = append(messages, schema.AssistantMessage(v.Content, []schema.ToolCall{}))
		}
	}
	// 进行对话
	response, err := ag.Generate(ctx, messages)
	if err != nil {
		logger.Error("Model call failed",
			zap.String("provider", p.Name()),
			zap.Error(err),
		)
		return &internalModel.ModelResponse{
			ModelName:    req.Models.Name,
			Provider:     p.Name(),
			Content:      "",
			Success:      false,
			Error:        err.Error(),
			ResponseTime: time.Since(startTime).Milliseconds(),
			StartTime:    startTime,
			EndTime:      time.Now(),
		}, err
	}

	endTime := time.Now()
	responseTime := endTime.Sub(startTime).Milliseconds()

	logger.Info("zhipu model response received",
		zap.String("provider", p.Name()),
		zap.Int64("response_time_ms", responseTime),
		zap.Int("content_length", len(response.Content)),
	)

	return &internalModel.ModelResponse{
		ModelName:    req.Models.Name,
		Provider:     p.Name(),
		Content:      response.Content,
		Success:      true,
		TokensUsed:   0, // aggo暂时不返回token信息
		ResponseTime: responseTime,
		StartTime:    startTime,
		EndTime:      endTime,
	}, nil
}

// Stream 流式调用zhipu模型
func (p *Provider) Stream(ctx context.Context, req *internalModel.CallProvidersRequest) (<-chan *internalModel.StreamChunk, error) {
	ch := make(chan *internalModel.StreamChunk, 10)

	go func() {
		defer close(ch)

		// TODO: 实现真实的流式调用
		// 这里先发送模拟数据
		chunks := []string{"这是", "zhipu", "模型的", "流式", "响应", "。"}

		for _, chunk := range chunks {
			select {
			case <-ctx.Done():
				return
			case ch <- &internalModel.StreamChunk{
				Model:   "gpt-4",
				Content: chunk,
				Done:    false,
			}:
				time.Sleep(100 * time.Millisecond) // 模拟延迟
			}
		}

		// 发送结束标记
		ch <- &internalModel.StreamChunk{
			Model:   "gpt-4",
			Content: "",
			Done:    true,
		}
	}()

	return ch, nil
}

// ValidateConfig 验证配置
func (p *Provider) ValidateConfig(config map[string]interface{}) error {
	if p.config.ApiKey == "" {
		return errors.New("zhipu API key is required")
	}
	return nil
}
