package service

import (
	"errors"
	"task-board/internal/domain"
	"task-board/internal/repository"
)

type BoardService interface {
	CreateBoard(ownerID uint, title, description string) (*domain.Board, error)
	GetBoards(ownerID uint) ([]domain.Board, error)
	GetBoard(boardID, userID uint) (*domain.Board, error)
	UpdateBoard(boardID, userID uint, title, description string) (*domain.Board, error)
	DeleteBoard(boardID, userID uint) error
}

type boardService struct {
	boardRepo repository.BoardRepository
}

func NewBoardService(boardRepo repository.BoardRepository) BoardService {
	return &boardService{
		boardRepo: boardRepo,
	}
}

func (s *boardService) CreateBoard(ownerID uint, title, description string) (*domain.Board, error) {
	board := &domain.Board{
		Title:       title,
		Description: description,
		OwnerID:     ownerID,
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) GetBoards(ownerID uint) ([]domain.Board, error) {
	return s.boardRepo.GetByOwnerID(ownerID)
}

func (s *boardService) GetBoard(boardID, userID uint) (*domain.Board, error) {
	board, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the board
	if board.OwnerID != userID {
		return nil, errors.New("unauthorized access to board")
	}

	return board, nil
}

func (s *boardService) UpdateBoard(boardID, userID uint, title, description string) (*domain.Board, error) {
	board, err := s.GetBoard(boardID, userID)
	if err != nil {
		return nil, err
	}

	board.Title = title
	board.Description = description

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) DeleteBoard(boardID, userID uint) error {
	_, err := s.GetBoard(boardID, userID)
	if err != nil {
		return err
	}

	return s.boardRepo.Delete(boardID)
}
