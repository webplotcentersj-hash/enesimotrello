package domain

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	OwnerID     uint           `json:"owner_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Owner *User  `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Tasks []Task `json:"tasks,omitempty" gorm:"foreignKey:BoardID"`
}
