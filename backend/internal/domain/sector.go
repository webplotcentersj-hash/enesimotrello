package domain

import "time"

// Sector represents a work sector/area
type Sector struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Nombre            string    `json:"nombre" gorm:"type:varchar(100);not null"`
	Color             string    `json:"color" gorm:"type:varchar(7);default:'#6B7280'"`
	Activo            bool      `json:"activo" gorm:"default:true"`
	OrdenVisualizacion int      `json:"orden_visualizacion" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Ordenes []OrderSector `json:"ordenes,omitempty" gorm:"foreignKey:IDSector"`
}

// TableName specifies the table name for Sector
func (Sector) TableName() string {
	return "sectores"
}

