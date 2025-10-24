package model

// TestRequest 多模型测试请求
type TestRequest struct {
	Prompts PromptSet  `json:"prompts" binding:"required"`
	Models  []ModelReq `json:"models" binding:"required,min=1"`
}

type CallProvidersRequest struct {
	Prompts PromptSet `json:"prompts" binding:"required"`
	Models  ModelReq  `json:"models"`
}

// PromptSet 提示词配置
type PromptSet struct {
	System  string    `json:"system"`                  // 系统提示词
	User    string    `json:"user" binding:"required"` // 用户提示词
	AI      string    `json:"ai"`                      // AI助手预设回复
	Message []Message `json:"message"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ModelReq 单个模型请求配置
type ModelReq struct {
	Name     string                 `json:"name" binding:"required"`     // 模型名称 如: gpt-4
	Provider string                 `json:"provider" binding:"required"` // 提供商 如: openai
	Config   map[string]interface{} `json:"config"`                      // 模型参数配置
}

// ModelConfig 模型配置参数
type ModelConfig struct {
	Temperature float64 `json:"temperature,omitempty"` // 温度参数
	MaxTokens   int     `json:"max_tokens,omitempty"`  // 最大token数
	TopP        float64 `json:"top_p,omitempty"`       // Top-P采样
	Stream      bool    `json:"stream,omitempty"`      // 是否流式响应
}

// SaveTemplateRequest 保存模板请求
type SaveTemplateRequest struct {
	Name         string `json:"name" binding:"required"`
	SystemPrompt string `json:"system_prompt"`
	UserPrompt   string `json:"user_prompt"`
	AIPrompt     string `json:"ai_prompt"`
	Description  string `json:"description"`
}

// SaveModelConfigRequest 保存模型配置请求
type SaveModelConfigRequest struct {
	Provider      string                 `json:"provider" binding:"required"`
	ModelName     string                 `json:"model_name" binding:"required"`
	ApiKey        string                 `json:"api_key"`
	BaseURL       string                 `json:"base_url"`
	DefaultConfig map[string]interface{} `json:"default_config"`
	Enabled       bool                   `json:"enabled"`
}
