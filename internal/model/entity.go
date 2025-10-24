package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONField 自定义JSON字段类型
type JSONField map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONField) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONField) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// TestRecord 测试记录表
type TestRecord struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Prompts   JSONField `gorm:"type:json;not null" json:"prompts"`
	Models    JSONField `gorm:"type:json;not null" json:"models"`
	Results   JSONField `gorm:"type:json;not null" json:"results"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (TestRecord) TableName() string {
	return "test_records"
}

// PromptTemplate 提示词模板表
type PromptTemplate struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null;index:idx_name" json:"name"`
	SystemPrompt string    `gorm:"type:text" json:"system_prompt"`
	UserPrompt   string    `gorm:"type:text" json:"user_prompt"`
	AIPrompt     string    `gorm:"type:text" json:"ai_prompt"`
	Description  string    `gorm:"type:varchar(500)" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index:idx_created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (PromptTemplate) TableName() string {
	return "prompt_templates"
}

// ModelConfigEntity 模型配置表
type ModelConfigEntity struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Provider      string    `gorm:"type:varchar(50);not null;uniqueIndex:uk_provider_model" json:"provider"`
	ModelName     string    `gorm:"type:varchar(100);not null;uniqueIndex:uk_provider_model" json:"model_name"`
	ApiKey        string    `gorm:"type:varchar(500)" json:"api_key"`
	BaseURL       string    `gorm:"type:varchar(500)" json:"base_url"`
	DefaultConfig JSONField `gorm:"type:json" json:"default_config"`
	Enabled       bool      `gorm:"type:tinyint(1);default:1;index:idx_enabled" json:"enabled"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (ModelConfigEntity) TableName() string {
	return "model_configs"
}