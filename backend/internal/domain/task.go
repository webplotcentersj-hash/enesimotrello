package domain

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Status      TaskStatus     `json:"status" gorm:"default:'todo'"`
	Priority    TaskPriority   `json:"priority" gorm:"default:'medium'"`
	BoardID     uint           `json:"board_id" gorm:"not null"`
	AssigneeID  *uint          `json:"assignee_id"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Board    *Board `json:"board,omitempty" gorm:"foreignKey:BoardID"`
	Assignee *User  `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
}
