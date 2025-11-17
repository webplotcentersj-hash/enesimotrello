package repository

import (
	"task-board/internal/domain"

	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(board *domain.Board) error
	GetByID(id uint) (*domain.Board, error)
	GetByOwnerID(ownerID uint) ([]domain.Board, error)
	Update(board *domain.Board) error
	Delete(id uint) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *domain.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) GetByID(id uint) (*domain.Board, error) {
	var board domain.Board
	err := r.db.Preload("Owner").Preload("Tasks").Preload("Tasks.Assignee").First(&board, id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) GetByOwnerID(ownerID uint) ([]domain.Board, error) {
	var boards []domain.Board
	err := r.db.Where("owner_id = ?", ownerID).Preload("Tasks").Find(&boards).Error
	return boards, err
}

func (r *boardRepository) Update(board *domain.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Board{}, id).Error
}
