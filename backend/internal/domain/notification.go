package domain

import "time"

// Notification represents a notification for a user
type Notification struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UsuarioDestino string   `json:"usuario_destino" gorm:"type:varchar(255);not null"`
	Tipo          string    `json:"tipo" gorm:"type:varchar(50);default:'mencion'"`
	Mensaje       string    `json:"mensaje" gorm:"type:text;not null"`
	IDOrden       *uint     `json:"id_orden" gorm:"column:id_orden;index"`
	Leida         bool      `json:"leida" gorm:"default:false"`
	Timestamp     time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden *Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notificaciones"
}

// UserNotification represents a user notification (newer format)
type UserNotification struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"column:user_id;not null;index"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description *string   `json:"description" gorm:"type:text"`
	Type        string    `json:"type" gorm:"type:varchar(50);default:'info'"`
	OrdenID     *uint     `json:"orden_id" gorm:"column:orden_id;index"`
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	User  User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Orden *Order `json:"orden,omitempty" gorm:"foreignKey:OrdenID"`
}

// TableName specifies the table name for UserNotification
func (UserNotification) TableName() string {
	return "user_notifications"
}

// NotificationViewed represents viewed notifications
type NotificationViewed struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	IDUsuario  uint      `json:"id_usuario" gorm:"column:id_usuario;not null;index"`
	IDHistorial uint     `json:"id_historial" gorm:"column:id_historial;not null;index"`
	Timestamp  time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Usuario  User            `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
	Historial MovementHistory `json:"historial,omitempty" gorm:"foreignKey:IDHistorial"`
}

// TableName specifies the table name for NotificationViewed
func (NotificationViewed) TableName() string {
	return "notificaciones_vistas"
}

// AlertSent represents sent alerts
type AlertSent struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	IDOrden    uint      `json:"id_orden" gorm:"column:id_orden;not null;index"`
	TipoAlerta string    `json:"tipo_alerta" gorm:"type:enum('estancada','plazo');not null"`
	Timestamp  time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
}

// TableName specifies the table name for AlertSent
func (AlertSent) TableName() string {
	return "alertas_enviadas"
}

// SmartAlert represents intelligent alerts
type SmartAlert struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	TipoAlerta    string     `json:"tipo_alerta" gorm:"type:enum('retraso_predicho','sobrecarga_operario','cuello_botella','eficiencia_baja');not null"`
	Prioridad     string     `json:"prioridad" gorm:"type:enum('baja','media','alta','critica');not null"`
	Titulo        string     `json:"titulo" gorm:"type:varchar(255);not null"`
	Descripcion   *string    `json:"descripcion" gorm:"type:text"`
	DatosContexto *string    `json:"datos_contexto" gorm:"type:longtext"` // JSON stored as string
	FechaCreacion time.Time  `json:"fecha_creacion" gorm:"default:CURRENT_TIMESTAMP"`
	FechaResuelto *time.Time `json:"fecha_resuelto"`
	Resuelto      bool       `json:"resuelto" gorm:"default:false"`
}

// TableName specifies the table name for SmartAlert
func (SmartAlert) TableName() string {
	return "smart_alerts"
}

