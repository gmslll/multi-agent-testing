package model

import "time"

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TestResult 多模型测试结果
type TestResult struct {
	Results   map[string]*ModelResponse `json:"results"`   // 各模型的响应结果
	StartTime time.Time                 `json:"start_time"`
	EndTime   time.Time                 `json:"end_time"`
	Duration  int64                     `json:"duration"` // 总耗时(毫秒)
}

// ModelResponse 单个模型的响应结果
type ModelResponse struct {
	ModelName    string        `json:"model_name"`
	Provider     string        `json:"provider"`
	Content      string        `json:"content"`       // 模型回复内容
	Error        string        `json:"error,omitempty"` // 错误信息
	Success      bool          `json:"success"`
	TokensUsed   int           `json:"tokens_used,omitempty"` // 使用的token数
	ResponseTime int64         `json:"response_time"` // 响应时间(毫秒)
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
}

// StreamChunk 流式响应数据块
type StreamChunk struct {
	Model   string `json:"model"`   // 模型名称
	Content string `json:"content"` // 内容片段
	Done    bool   `json:"done"`    // 是否结束
	Error   string `json:"error,omitempty"`
}

// ModelListResponse 模型列表响应
type ModelListResponse struct {
	Models []ModelInfo `json:"models"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Enabled  bool   `json:"enabled"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}