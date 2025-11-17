package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Nombre      string    `json:"nombre" gorm:"type:varchar(100);not null;uniqueIndex"`
	PasswordHash string   `json:"-" gorm:"column:password_hash;type:varchar(255);not null"`
	Rol         string    `json:"rol" gorm:"type:enum('administracion','taller','mostrador');not null"`
	LastSeen    *time.Time `json:"last_seen" gorm:"column:last_seen"`

	// Relations
	OrdenesCreadas    []Order            `json:"ordenes_creadas,omitempty" gorm:"foreignKey:IDUsuarioCreador"`
	HistorialMovimientos []MovementHistory `json:"historial_movimientos,omitempty" gorm:"foreignKey:IDUsuario"`
	MensajesChat      []ChatMessage      `json:"mensajes_chat,omitempty" gorm:"foreignKey:IDUsuario"`
	Notificaciones    []UserNotification `json:"notificaciones,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "usuarios"
}

// CheckPassword verifies a password against the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}
