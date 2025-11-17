package service

import (
	"errors"
	"task-board/internal/domain"
	"task-board/internal/repository"
	"time"
)

type TaskService interface {
	CreateTask(boardID, userID uint, title, description string, priority domain.TaskPriority, assigneeID *uint, dueDate *time.Time) (*domain.Task, error)
	GetTasks(boardID, userID uint) ([]domain.Task, error)
	GetTask(taskID, userID uint) (*domain.Task, error)
	UpdateTask(taskID, userID uint, title, description string, status domain.TaskStatus, priority domain.TaskPriority, assigneeID *uint, dueDate *time.Time) (*domain.Task, error)
	DeleteTask(taskID, userID uint) error
}

type taskService struct {
	taskRepo  repository.TaskRepository
	boardRepo repository.BoardRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo:  taskRepo,
		boardRepo: nil, // Will be set by dependency injection
	}
}

func (s *taskService) SetBoardRepo(boardRepo repository.BoardRepository) {
	s.boardRepo = boardRepo
}

func (s *taskService) CreateTask(boardID, userID uint, title, description string, priority domain.TaskPriority, assigneeID *uint, dueDate *time.Time) (*domain.Task, error) {
	// Verify user has access to the board
	board, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return nil, err
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized access to board")
	}

	task := &domain.Task{
		Title:       title,
		Description: description,
		Status:      domain.StatusTodo,
		Priority:    priority,
		BoardID:     boardID,
		AssigneeID:  assigneeID,
		DueDate:     dueDate,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTasks(boardID, userID uint) ([]domain.Task, error) {
	// Verify user has access to the board
	board, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return nil, err
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized access to board")
	}

	return s.taskRepo.GetByBoardID(boardID)
}

func (s *taskService) GetTask(taskID, userID uint) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Verify user has access to the board
	board, err := s.boardRepo.GetByID(task.BoardID)
	if err != nil {
		return nil, err
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized access to task")
	}

	return task, nil
}

func (s *taskService) UpdateTask(taskID, userID uint, title, description string, status domain.TaskStatus, priority domain.TaskPriority, assigneeID *uint, dueDate *time.Time) (*domain.Task, error) {
	task, err := s.GetTask(taskID, userID)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.Status = status
	task.Priority = priority
	task.AssigneeID = assigneeID
	task.DueDate = dueDate

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(taskID, userID uint) error {
	_, err := s.GetTask(taskID, userID)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(taskID)
}
