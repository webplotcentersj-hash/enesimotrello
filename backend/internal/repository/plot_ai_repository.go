package repository

import (
	"task-board/internal/domain"

	"gorm.io/gorm"
)

type PlotAIRepository interface {
	CreateMessage(message *domain.PlotAIChatMessage) error
	GetUserMessages(userID uint, limit int) ([]domain.PlotAIChatMessage, error)
	GetConfig(key string) (*domain.PlotAIConfig, error)
	SetConfig(key, value string) error
}

type plotAIRepository struct {
	db *gorm.DB
}

func NewPlotAIRepository(db *gorm.DB) PlotAIRepository {
	return &plotAIRepository{db: db}
}

func (r *plotAIRepository) CreateMessage(message *domain.PlotAIChatMessage) error {
	return r.db.Create(message).Error
}

func (r *plotAIRepository) GetUserMessages(userID uint, limit int) ([]domain.PlotAIChatMessage, error) {
	var messages []domain.PlotAIChatMessage
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&messages).Error
	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, err
}

func (r *plotAIRepository) GetConfig(key string) (*domain.PlotAIConfig, error) {
	var config domain.PlotAIConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (r *plotAIRepository) SetConfig(key, value string) error {
	config := &domain.PlotAIConfig{
		Key:   key,
		Value: value,
	}
	return r.db.Where("key = ?", key).Assign(domain.PlotAIConfig{Value: value}).FirstOrCreate(config).Error
}

