package domain

import "time"

// OnlineUser represents an online user
type OnlineUser struct {
	UserID    uint      `json:"user_id" gorm:"column:user_id;primaryKey"`
	UserNombre string   `json:"user_nombre" gorm:"type:varchar(100);not null"`
	LastSeen  time.Time `json:"last_seen" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for OnlineUser
func (OnlineUser) TableName() string {
	return "online_users"
}

// StatsCache represents cached statistics
type StatsCache struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CacheKey   string    `json:"cache_key" gorm:"type:varchar(255);not null;uniqueIndex"`
	CacheValue string    `json:"cache_value" gorm:"type:longtext"` // JSON stored as string
	ExpiresAt  time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for StatsCache
func (StatsCache) TableName() string {
	return "stats_cache"
}

// TrendingMetric represents trending metrics
type TrendingMetric struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Fecha        time.Time `json:"fecha" gorm:"type:date;not null;index"`
	Metrica      string    `json:"metrica" gorm:"type:varchar(100);not null"`
	Valor        float64   `json:"valor" gorm:"type:decimal(15,4);not null"`
	Categoria    *string   `json:"categoria" gorm:"type:varchar(100)"`
	Subcategoria *string   `json:"subcategoria" gorm:"type:varchar(100)"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for TrendingMetric
func (TrendingMetric) TableName() string {
	return "trending_metrics"
}

// PredictionMetric represents prediction metrics for orders
type PredictionMetric struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	OrdenID             uint       `json:"orden_id" gorm:"column:orden_id;not null;index"`
	NumeroOP            *string    `json:"numero_op" gorm:"type:varchar(50)"`
	TiempoPredichoHoras *float64   `json:"tiempo_predicho_horas" gorm:"type:decimal(10,2)"`
	TiempoRealHoras     *float64   `json:"tiempo_real_horas" gorm:"type:decimal(10,2)"`
	ConfianzaPrediccion *int       `json:"confianza_prediccion"`
	ErrorAbsoluto       *float64   `json:"error_absoluto" gorm:"type:decimal(10,2)"`
	ErrorPorcentual     *float64   `json:"error_porcentual" gorm:"type:decimal(10,2)"`
	FactoresAplicados   *string    `json:"factores_aplicados" gorm:"type:text"`
	FechaPrediccion     time.Time  `json:"fecha_prediccion" gorm:"default:CURRENT_TIMESTAMP"`
	FechaCompletado     *time.Time `json:"fecha_completado"`

	// Relations
	Orden Order `json:"orden,omitempty" gorm:"foreignKey:OrdenID"`
}

// TableName specifies the table name for PredictionMetric
func (PredictionMetric) TableName() string {
	return "prediction_metrics"
}

