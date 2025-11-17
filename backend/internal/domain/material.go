package domain

// Material represents a material used in orders
type Material struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Codigo      *string `json:"codigo" gorm:"type:varchar(50)"`
	Descripcion string `json:"descripcion" gorm:"type:varchar(255);not null"`

	// Relations
	Ordenes []OrderMaterial `json:"ordenes,omitempty" gorm:"foreignKey:IDMaterial"`
}

// TableName specifies the table name for Material
func (Material) TableName() string {
	return "materiales"
}

