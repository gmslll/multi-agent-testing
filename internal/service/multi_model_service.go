package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multi-agent-testing/backend/internal/config"
	"github.com/multi-agent-testing/backend/internal/model"
	"github.com/multi-agent-testing/backend/internal/providers/base"
	"github.com/multi-agent-testing/backend/internal/providers/deepseek"
	"github.com/multi-agent-testing/backend/internal/providers/minimax"
	"github.com/multi-agent-testing/backend/internal/providers/openai"
	"github.com/multi-agent-testing/backend/internal/providers/zhipu"
	"github.com/multi-agent-testing/backend/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// MultiModelService 多模型测试服务
type MultiModelService struct {
	providers map[string]base.ModelProvider
	config    *config.Config
}

// NewMultiModelService 创建多模型服务
func NewMultiModelService(cfg *config.Config) *MultiModelService {
	service := &MultiModelService{
		providers: make(map[string]base.ModelProvider),
		config:    cfg,
	}

	// 初始化各个模型提供者
	service.initProviders()

	return service
}

// initProviders 初始化模型提供者
func (s *MultiModelService) initProviders() {
	// 初始化OpenAI
	if modelCfg, ok := s.config.Models["openai"]; ok && modelCfg.Enabled {
		provider := openai.NewProvider(base.ProviderConfig{
			ApiKey:  modelCfg.ApiKey,
			BaseURL: modelCfg.BaseURL,
			Timeout: int(modelCfg.Timeout.Seconds()),
		})
		s.providers["openai"] = provider
		logger.Info("OpenAI provider initialized", zap.String("base_url", modelCfg.BaseURL))
	}

	if modelCfg, ok := s.config.Models["deepseek"]; ok && modelCfg.Enabled {
		provider := deepseek.NewProvider(base.ProviderConfig{
			ApiKey:  modelCfg.ApiKey,
			BaseURL: modelCfg.BaseURL,
			Timeout: int(modelCfg.Timeout.Seconds()),
		})
		s.providers["deepseek"] = provider
		logger.Info("deepseek provider initialized", zap.String("base_url", modelCfg.BaseURL))
	}

	if modelCfg, ok := s.config.Models["minimax"]; ok && modelCfg.Enabled {
		provider := minimax.NewProvider(base.ProviderConfig{
			ApiKey:  modelCfg.ApiKey,
			BaseURL: modelCfg.BaseURL,
			Timeout: int(modelCfg.Timeout.Seconds()),
		})
		s.providers["minimax"] = provider
		logger.Info("minimax provider initialized", zap.String("base_url", modelCfg.BaseURL))
	}
	if modelCfg, ok := s.config.Models["zhipu"]; ok && modelCfg.Enabled {
		provider := zhipu.NewProvider(base.ProviderConfig{
			ApiKey:  modelCfg.ApiKey,
			BaseURL: modelCfg.BaseURL,
			Timeout: int(modelCfg.Timeout.Seconds()),
		})
		s.providers["zhipu"] = provider
		logger.Info("zhipu provider initialized", zap.String("base_url", modelCfg.BaseURL))
	}
	// TODO: 初始化其他提供者(Claude, GLM, Qwen等)
	// if modelCfg, ok := s.config.Models["anthropic"]; ok && modelCfg.Enabled {
	// 	provider := anthropic.NewProvider(...)
	// 	s.providers["anthropic"] = provider
	// }
}

// ExecuteTest 执行多模型测试
func (s *MultiModelService) ExecuteTest(ctx context.Context, req *model.TestRequest) (*model.TestResult, error) {
	startTime := time.Now()

	logger.Info("Starting multi-model test",
		zap.Int("model_count", len(req.Models)),
		zap.String("user_prompt", req.Prompts.User),
	)

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// 使用errgroup并发调用多个模型
	g, ctx := errgroup.WithContext(ctx)
	results := make(map[string]*model.ModelResponse)
	var mu sync.Mutex
	for _, _modelReq := range req.Models {
		modelReq := _modelReq // 避免闭包陷阱
		callProvidersRequest := &model.CallProvidersRequest{
			Prompts: req.Prompts,
			Models:  modelReq,
		}
		g.Go(func() error {
			// 获取对应的提供者
			provider, exists := s.providers[modelReq.Provider]
			if !exists {
				return fmt.Errorf("provider %s not found or not enabled", modelReq.Provider)
			}

			// 调用模型
			resp, err := provider.Call(ctx, callProvidersRequest)
			if err != nil {
				logger.Error("Model call failed",
					zap.String("provider", modelReq.Provider),
					zap.String("model", modelReq.Name),
					zap.Error(err),
				)

				// 记录错误但继续执行其他模型
				mu.Lock()
				results[modelReq.Name] = &model.ModelResponse{
					ModelName:    modelReq.Name,
					Provider:     modelReq.Provider,
					Content:      "",
					Error:        err.Error(),
					Success:      false,
					ResponseTime: 0,
					StartTime:    time.Now(),
					EndTime:      time.Now(),
				}
				mu.Unlock()
				return nil // 返回nil以允许其他模型继续执行
			}

			// 保存成功的响应
			mu.Lock()
			resp.ModelName = modelReq.Name
			results[modelReq.Name] = resp
			mu.Unlock()

			logger.Info("Model response received",
				zap.String("provider", modelReq.Provider),
				zap.String("model", modelReq.Name),
				zap.Int64("response_time_ms", resp.ResponseTime),
			)

			return nil
		})
	}

	// 等待所有模型调用完成
	if err := g.Wait(); err != nil {
		logger.Error("Multi-model test failed", zap.Error(err))
		return nil, err
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	result := &model.TestResult{
		Results:   results,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
	}

	logger.Info("Multi-model test completed",
		zap.Int("total_models", len(req.Models)),
		zap.Int("success_count", countSuccessful(results)),
		zap.Int64("total_duration_ms", duration),
	)

	return result, nil
}

// GetAvailableModels 获取可用的模型列表
func (s *MultiModelService) GetAvailableModels() []model.ModelInfo {
	models := []model.ModelInfo{}

	// OpenAI模型
	if _, ok := s.providers["openai"]; ok {
		models = append(models,
			model.ModelInfo{Name: "gpt-4.1", Provider: "openai", Enabled: true},
			model.ModelInfo{Name: "gpt-5-mini", Provider: "openai", Enabled: true},
		)
	}

	// DeepSeek模型
	if _, ok := s.providers["deepseek"]; ok {
		models = append(models,
			model.ModelInfo{Name: "deepseek-chat", Provider: "deepseek", Enabled: true},
			model.ModelInfo{Name: "deepseek-reasoner", Provider: "deepseek", Enabled: true},
		)
	}

	// minimax模型
	if _, ok := s.providers["minimax"]; ok {
		models = append(models,
			model.ModelInfo{Name: "MiniMax-M1", Provider: "minimax", Enabled: true},
			model.ModelInfo{Name: "MiniMax-Text-01", Provider: "minimax", Enabled: true},
		)
	}
	// minimax模型
	if _, ok := s.providers["zhipu"]; ok {
		models = append(models,
			model.ModelInfo{Name: "charglm-4", Provider: "zhipu", Enabled: true},
		)
	}

	// TODO: 添加其他提供者的模型

	return models
}

// countSuccessful 统计成功的响应数量
func countSuccessful(results map[string]*model.ModelResponse) int {
	count := 0
	for _, resp := range results {
		if resp.Success {
			count++
		}
	}
	return count
}
