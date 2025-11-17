package domain

import "time"

// ChatRoom represents a chat room
type ChatRoom struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nombre    string    `json:"nombre" gorm:"type:varchar(255);not null"`
	Tipo      string    `json:"tipo" gorm:"type:enum('publico','privado');default:'publico'"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Messages []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:RoomID"`
}

// TableName specifies the table name for ChatRoom
func (ChatRoom) TableName() string {
	return "chat_rooms"
}

// ChatMessage represents a message in a chat room
type ChatMessage struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	RoomID        uint      `json:"room_id" gorm:"column:room_id;not null;default:1;index"`
	IDUsuario     uint      `json:"id_usuario" gorm:"column:id_usuario;not null;index"`
	NombreUsuario string    `json:"nombre_usuario" gorm:"type:varchar(100);not null"`
	Mensaje       string    `json:"mensaje" gorm:"type:text;not null"`
	MessageType   *string   `json:"message_type" gorm:"type:varchar(50)"`
	Timestamp     time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Room   ChatRoom `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Usuario User    `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
}

// TableName specifies the table name for ChatMessage
func (ChatMessage) TableName() string {
	return "chat_messages"
}

