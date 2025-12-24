package payment

import (
	"time"
)

// Payment represents a single payment

type Payment struct {
	ID         string             `gorm:"primaryKey"`
	Amount     int64              `gorm:"not null"`
	State      State              `gorm:"not null"`
	Operations []PaymentOperation `gorm:"foreignKey:PaymentID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PaymentOperation struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	PaymentID   string `gorm:"index;not null"` // Foreign key
	OperationID string `gorm:"not null"`       // Idempotency key
	Operation   string `gorm:"not null"`       // Operation type (e.g., "CAPTURE")
	Result      string `gorm:"not null"`       // Result of operation (e.g., "success", "failed")
	CreatedAt   time.Time
}

func (PaymentOperation) TableName() string {
	return "payment_operations"
}
