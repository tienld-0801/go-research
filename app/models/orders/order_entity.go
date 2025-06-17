package order_model

import (
	"time"

	"github.com/lib/pq"

	receipt_model "go-backend/app/models/receipts"
)

type Order struct {
	StaffId        string         `gorm:"primaryKey" json:"staff_id"`
	ReceiptId      string         `gorm:"primaryKey;type:varchar(100)" json:"receipt_id"`
	SequenceNumber int            `gorm:"primaryKey" json:"sequence_number"`
	ProductId      string         `gorm:"type:varchar(100)" json:"product_id"`
	Size           string         `gorm:"type:varchar(50)" json:"size"`
	Sweet          string         `gorm:"type:varchar(50)" json:"sweet"`
	Ice            string         `gorm:"type:varchar(50)" json:"ice"`
	Toppings       pq.StringArray `gorm:"type:text[]" json:"toppings"`
	Amount         int            `json:"amount"`
	SinglePrice    float64        `json:"single_price"`
	SumPrice       float64        `json:"sum_price"`
	Profit         float64        `json:"profit"`
	Note           string         `gorm:"type:text" json:"note"`
	Date           time.Time      `json:"date"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoCreateTime" json:"updated_at"`

	// Association
	Receipt receipt_model.Receipt `gorm:"foreignKey:ReceiptId;references:ID" json:"receipt"`
}
