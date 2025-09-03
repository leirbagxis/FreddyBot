package models

import "time"

// Planos de assinatura
type SubscriptionPlan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`       // Ex: "Free", "Premium"
	PricePix   float64   `json:"pricePix"`   // Valor em reais (PIX)
	PriceStars int       `json:"priceStars"` // Valor em Telegram Stars
	Duration   int       `json:"duration"`   // Duração em dias (0 = vitalício)
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Assinatura de um canal em um plano
type Subscription struct {
	ID              string           `gorm:"type:text;primaryKey" json:"id"`
	ChannelID       int64            `json:"channelId"`
	PlanID          string           `json:"planId"`
	Status          string           `json:"status"` // active, expired, cancelled, pending
	StartDate       time.Time        `json:"startDate"`
	EndDate         time.Time        `json:"endDate"`
	PaymentMethod   string           `json:"paymentMethod"` // pix, telegram_star, gift_card
	AppliedCouponID *string          `json:"appliedCouponId,omitempty"`
	RedeemedGiftID  *string          `json:"redeemedGiftId,omitempty"`
	Plan            SubscriptionPlan `gorm:"foreignKey:PlanID" json:"plan"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

// Cupons de desconto
type Coupon struct {
	ID          string    `gorm:"type:text;primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex" json:"code"` // Ex: PROMO50
	DiscountPct int       `json:"discountPct"`             // Desconto em %
	DiscountAmt int64     `json:"discountAmt"`             // Desconto fixo em centavos
	ValidFrom   time.Time `json:"validFrom"`
	ValidUntil  time.Time `json:"validUntil"`
	MaxUses     int       `json:"maxUses"`
	UsedCount   int       `json:"usedCount"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Gift Cards
type GiftCard struct {
	ID        string     `gorm:"type:text;primaryKey" json:"id"`
	Code      string     `gorm:"uniqueIndex" json:"code"` // Ex: GIFT-XXXX-YYYY
	PlanID    *string    `json:"planId,omitempty"`        // Se atrelado a um plano específico
	Credit    int64      `json:"credit"`                  // Crédito em centavos
	IsUsed    bool       `gorm:"default:false" json:"isUsed"`
	UsedBy    *int64     `json:"usedBy,omitempty"` // ChannelID ou UserID
	UsedAt    *time.Time `json:"usedAt,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
