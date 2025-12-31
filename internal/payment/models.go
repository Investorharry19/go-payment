package payment

import (
	"time"
)

// Payment represents a single payment

type Payment struct {
	ID      string `gorm:"primaryKey"`
	UserID  string `gorm:"index;not null"` // usr_xxx
	OrderID string `gorm:"index;not null"`

	Amount     int64              `gorm:"not null"`
	State      State              `gorm:"not null"`
	Operations []PaymentOperation `gorm:"foreignKey:PaymentID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PaymentOperation struct {
	ID          uint   `gorm:"primaryKey"`
	PaymentID   string `gorm:"not null;index:idx_payment_operation_id,unique"`
	OperationID string `gorm:"not null;index:idx_payment_operation_id,unique"`

	Operation string `gorm:"not null"` // AUTHORIZE, CAPTURE, VOID, REFUND
	Amount    int64  `gorm:"not null"`
	Result    string `gorm:"not null"` // success, failed

	BankReference string
	CreatedAt     time.Time
}

func (PaymentOperation) TableName() string {
	return "payment_operations"
}
