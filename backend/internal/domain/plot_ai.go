package domain

import (
	"time"

	"gorm.io/gorm"
)

// PlotAIChatMessage represents a chat message specifically for Plot AI
// Note: Plot AI messages can be stored in the general chat_messages table
// with a specific room_id, or in a separate table. For now, we'll use
// a separate table for clarity, but the service can also use ChatMessage
// from chat.go with a specific room_id.
type PlotAIChatMessage struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null"` // "user" or "model"
	Content   string    `json:"content" gorm:"type:text;not null"`
	Timestamp time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for PlotAIChatMessage
func (PlotAIChatMessage) TableName() string {
	return "plot_ai_chat_messages"
}

// PlotAIConfig represents configuration for Plot AI
type PlotAIConfig struct {
	gorm.Model
	Key   string `json:"key" gorm:"type:varchar(255);uniqueIndex;not null"`
	Value string `json:"value" gorm:"type:text"`
}

// TableName specifies the table name for PlotAIConfig
func (PlotAIConfig) TableName() string {
	return "plot_ai_config"
}

// SendMessageRequest represents a request to send a message to Plot AI
type SendMessageRequest struct {
	Message string                `json:"message" binding:"required"`
	History []PlotAIChatMessage   `json:"history,omitempty"`
}

// SendMessageResponse represents a response from Plot AI
type SendMessageResponse struct {
	Reply   string                `json:"reply"`
	History []PlotAIChatMessage   `json:"history,omitempty"`
	Timestamp time.Time           `json:"timestamp"`
}

