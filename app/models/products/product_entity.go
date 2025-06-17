package products_model

import (
	"time"
)

type Product struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Category     string    `gorm:"type:varchar(50)" json:"category"`
	Description  string    `gorm:"type:text" json:"description"`
	Photo        string    `gorm:"type:text" json:"photo"`
	Recipe       string    `gorm:"type:text" json:"recipe"`
	PriceGroupID string    `gorm:"type:uuid" json:"price_group_id"`
	ActiveFlag   bool      `gorm:"default:true" json:"active_flag"`
	DeletedFlag  bool      `gorm:"default:false" json:"deleted_flag"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
