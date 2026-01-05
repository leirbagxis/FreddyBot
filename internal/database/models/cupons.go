package models

import "time"

type Coupon struct {
	ID            string `gorm:"primaryKey"`
	Code          string `gorm:"uniqueIndex;size:32"`
	DiscountType  string // "percent" | "fixed"
	DiscountValue int    // 10 = 10% | 100 = 1.00 stars
	MaxUses       int    // 0 = ilimitado
	UsedCount     int
	ExpiresAt     *time.Time

	IsPrivate  bool
	IsAuto     bool // primeira compra
	OnePerUser bool

	CreatedAt time.Time
}

type CouponUsage struct {
	ID       string `gorm:"primaryKey"`
	CouponID string `gorm:"index"`
	UserID   int64  `gorm:"index"`
	UsedAt   time.Time

	// proteção forte
	UniqueKey string `gorm:"uniqueIndex"`
}
