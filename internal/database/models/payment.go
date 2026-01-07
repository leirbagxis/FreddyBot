package models

import "time"

type PriceConfig struct {
	Key    string `db:"key" json:"key"`       // ex: "add_to_channel_fee"
	Title  string `db:"title" json:"title"`   // ex: "Taxa para adicionar em canal"
	Amount int64  `db:"amount" json:"amount"` // inteiro (Stars, no caso de XTR)

	Active    bool      `db:"active" json:"active"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Payment struct {
	ID        string `gorm:"primaryKey"`
	UserID    int64  `gorm:"index"`
	Amount    int
	Payload   string `gorm:"uniqueIndex"`
	CouponID  *string
	Status    string `gorm:"default:pending"`
	CreatedAt time.Time
	PaidAt    *time.Time
}

//`gorm:"default:nil"`
