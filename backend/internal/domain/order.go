package domain

import "time"

// Order represents a work order in the system
type Order struct {
	ID                      uint      `json:"id" gorm:"primaryKey"`
	NumeroOP                string    `json:"numero_op" gorm:"column:numero_op;not null"`
	Cliente                 string    `json:"cliente" gorm:"not null"`
	Descripcion             string    `json:"descripcion" gorm:"type:text"`
	FechaEntrega            time.Time `json:"fecha_entrega" gorm:"type:date;not null"`
	Estado                  string    `json:"estado" gorm:"not null;default:'Pendiente'"`
	Prioridad               string    `json:"prioridad" gorm:"not null;default:'Normal'"`
	FechaCreacion           time.Time `json:"fecha_creacion" gorm:"default:CURRENT_TIMESTAMP"`
	FechaIngreso            time.Time `json:"fecha_ingreso" gorm:"default:CURRENT_TIMESTAMP"`
	OperarioAsignado        string    `json:"operario_asignado" gorm:"type:varchar(100)"`
	Complejidad             string    `json:"complejidad" gorm:"type:enum('Baja','Media','Alta');default:'Media'"`
	Sector                  string    `json:"sector" gorm:"type:enum('Taller Gráfico','Mostrador');default:'Taller Gráfico'"`
	HoraEstimadaEntrega     *string   `json:"hora_estimada_entrega" gorm:"type:time"`
	HoraEntregaEfectiva     *string   `json:"hora_entrega_efectiva" gorm:"type:time"`
	IDUsuarioCreador        *uint     `json:"id_usuario_creador" gorm:"column:id_usuario_creador"`
	UsuarioTrabajandoID     *uint     `json:"usuario_trabajando_id" gorm:"column:usuario_trabajando_id"`
	UsuarioTrabajandoNombre *string   `json:"usuario_trabajando_nombre" gorm:"type:varchar(100)"`
	TimestampInicioTrabajo  *time.Time `json:"timestamp_inicio_trabajo" gorm:"column:timestamp_inicio_trabajo"`

	// Relations
	UsuarioCreador *User            `json:"usuario_creador,omitempty" gorm:"foreignKey:IDUsuarioCreador"`
	Materiales     []OrderMaterial  `json:"materiales,omitempty" gorm:"foreignKey:IDOrden"`
	Sectores       []OrderSector    `json:"sectores,omitempty" gorm:"foreignKey:IDOrden"`
	Archivos       []Attachment     `json:"archivos,omitempty" gorm:"foreignKey:IDOrden"`
	Historial      []MovementHistory `json:"historial,omitempty" gorm:"foreignKey:IDOrden"`
	Tareas         []Task           `json:"tareas,omitempty" gorm:"foreignKey:IDOrden"`
	Comentarios    []OrderComment   `json:"comentarios,omitempty" gorm:"foreignKey:IDOrden"`
	Enlaces        []OrderLink      `json:"enlaces,omitempty" gorm:"foreignKey:IDOrden"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "ordenes_trabajo"
}

// OrderMaterial represents the relationship between orders and materials
type OrderMaterial struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	IDOrden    uint    `json:"id_orden" gorm:"column:id_orden;not null;index"`
	IDMaterial uint    `json:"id_material" gorm:"column:id_material;not null;index"`
	Cantidad   float64 `json:"cantidad" gorm:"type:decimal(10,3);default:1.000"`

	// Relations
	Orden    Order    `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
	Material Material `json:"material,omitempty" gorm:"foreignKey:IDMaterial"`
}

// TableName specifies the table name for OrderMaterial
func (OrderMaterial) TableName() string {
	return "orden_materiales"
}

// OrderSector represents the relationship between orders and sectors
type OrderSector struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	IDOrden        uint      `json:"id_orden" gorm:"column:id_orden;not null;index"`
	IDSector       uint      `json:"id_sector" gorm:"column:id_sector;not null;index"`
	FechaAsignacion time.Time `json:"fecha_asignacion" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden  Order  `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
	Sector Sector `json:"sector,omitempty" gorm:"foreignKey:IDSector"`
}

// TableName specifies the table name for OrderSector
func (OrderSector) TableName() string {
	return "orden_sectores"
}

// Attachment represents a file attached to an order
type Attachment struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	IDOrden        uint      `json:"id_orden" gorm:"column:id_orden;not null;index"`
	NombreArchivo  string    `json:"nombre_archivo" gorm:"type:varchar(255);not null"`
	NombreOriginal string    `json:"nombre_original" gorm:"type:varchar(255);not null"`
	FechaSubida    time.Time `json:"fecha_subida" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
}

// TableName specifies the table name for Attachment
func (Attachment) TableName() string {
	return "archivos_adjuntos"
}

// MovementHistory represents the history of state changes for an order
type MovementHistory struct {
	ID                          uint       `json:"id" gorm:"primaryKey"`
	IDOrden                     uint       `json:"id_orden" gorm:"column:id_orden;not null;index"`
	IDUsuario                   uint       `json:"id_usuario" gorm:"column:id_usuario;not null;index"`
	NombreUsuario               string     `json:"nombre_usuario" gorm:"type:varchar(100);not null"`
	EstadoAnterior              *string    `json:"estado_anterior" gorm:"type:varchar(50)"`
	EstadoNuevo                 *string    `json:"estado_nuevo" gorm:"type:varchar(50)"`
	DuracionEstadoAnteriorSeg   *int       `json:"duracion_estado_anterior_seg"`
	Timestamp                   time.Time  `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`
	Comentario                  *string    `json:"comentario" gorm:"type:text"`

	// Relations
	Orden  Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
	Usuario User `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
}

// TableName specifies the table name for MovementHistory
func (MovementHistory) TableName() string {
	return "historial_movimientos"
}

// Task represents a task within an order
type Task struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	IDOrden         uint   `json:"id_orden" gorm:"column:id_orden;not null;index"`
	DescripcionTarea string `json:"descripcion_tarea" gorm:"type:varchar(255);not null"`
	EstadoKanban    string `json:"estado_kanban" gorm:"type:enum('Pendiente','En Proceso','Finalizado');default:'Pendiente'"`

	// Relations
	Orden Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
}

// TableName specifies the table name for Task
func (Task) TableName() string {
	return "tareas"
}

// OrderComment represents a comment on an order
type OrderComment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IDOrden   uint      `json:"id_orden" gorm:"column:id_orden;not null;index"`
	IDUsuario uint      `json:"id_usuario" gorm:"column:id_usuario;not null;index"`
	Comentario string   `json:"comentario" gorm:"type:text;not null"`
	Timestamp time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden  Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
	Usuario User `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
}

// TableName specifies the table name for OrderComment
func (OrderComment) TableName() string {
	return "comentarios_orden"
}

// OrderLink represents a link attached to an order
type OrderLink struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	IDOrden     uint      `json:"id_orden" gorm:"column:id_orden;not null;index"`
	URL         string    `json:"url" gorm:"type:text;not null"`
	Descripcion *string   `json:"descripcion" gorm:"type:text"`
	Timestamp   time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Orden Order `json:"orden,omitempty" gorm:"foreignKey:IDOrden"`
}

// TableName specifies the table name for OrderLink
func (OrderLink) TableName() string {
	return "enlaces_adjuntos"
}

