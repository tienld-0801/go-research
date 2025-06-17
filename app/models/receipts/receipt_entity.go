package receipt_model

import (
	"time"
)

type Receipt struct {
	ID          string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	StaffID     string    `gorm:"type:varchar(100);not null" json:"staff_id"`
	TotalAmount int       `gorm:"not null" json:"total_amount"`
	TotalPrice  float64   `gorm:"not null" json:"total_price"`
	Note        string    `gorm:"type:text" json:"note"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	DeletedFlag bool      `gorm:"default:false" json:"deleted_flag"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
