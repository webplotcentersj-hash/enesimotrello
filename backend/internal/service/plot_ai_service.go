package service

import (
	"task-board/internal/domain"
	"task-board/internal/repository"
)

type PlotAIService interface {
	SendMessage(userID uint, message string) (*domain.ChatResponse, error)
	GetHistory(userID uint, limit int) ([]domain.ChatMessage, error)
}

type plotAIService struct {
	repo          repository.PlotAIRepository
	geminiService *GeminiService
}

func NewPlotAIService(repo repository.PlotAIRepository, geminiService *GeminiService) PlotAIService {
	return &plotAIService{
		repo:          repo,
		geminiService: geminiService,
	}
}

func (s *plotAIService) SendMessage(userID uint, message string) (*domain.ChatResponse, error) {
	// Get conversation history
	history, err := s.repo.GetUserMessages(userID, 20)
	if err != nil {
		return nil, err
	}

	// Convert to format expected by Gemini
	geminiHistory := make([]struct {
		Role    string
		Content string
	}, len(history))
	
	for i, msg := range history {
		geminiHistory[i] = struct {
			Role    string
			Content string
		}{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Save user message
	userMsg := &domain.ChatMessage{
		UserID:  userID,
		Role:    "user",
		Content: message,
	}
	if err := s.repo.CreateMessage(userMsg); err != nil {
		return nil, err
	}

	// Get response from Gemini
	response, err := s.geminiService.SendMessage(message, geminiHistory)
	if err != nil {
		return nil, err
	}

	// Save assistant message
	assistantMsg := &domain.ChatMessage{
		UserID:  userID,
		Role:    "assistant",
		Content: response,
	}
	if err := s.repo.CreateMessage(assistantMsg); err != nil {
		return nil, err
	}

	// Get updated history
	updatedHistory, err := s.repo.GetUserMessages(userID, 20)
	if err != nil {
		return nil, err
	}

	return &domain.ChatResponse{
		Message: response,
		History: updatedHistory,
	}, nil
}

func (s *plotAIService) GetHistory(userID uint, limit int) ([]domain.ChatMessage, error) {
	return s.repo.GetUserMessages(userID, limit)
}

