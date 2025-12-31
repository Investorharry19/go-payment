package payment

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PaymentStoreDB is a DB-backed payment store
type PaymentStoreDB struct {
	DB *gorm.DB
}

// Constructor
func NewPaymentStoreDB(db *gorm.DB) *PaymentStoreDB {
	return &PaymentStoreDB{DB: db}
}

func (s *PaymentStoreDB) Create(id string, amount int64, userId, orderId string) (*Payment, error) {
	p := &Payment{
		ID:      id,
		Amount:  amount,
		State:   Initiated,
		UserID:  userId,
		OrderID: orderId,
	}

	if err := s.DB.Create(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PaymentStoreDB) Get(id string) (*Payment, error) {
	var p Payment
	if err := s.DB.Preload("Operations").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PaymentStoreDB) Apply(
	bank Bank,
	paymentID string,
	operationID string,
	operation Operation,
) error {

	return s.DB.Transaction(func(tx *gorm.DB) error {
		//  Lock the payment row for update to prevent concurrent modification
		var p Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&p, "id = ?", paymentID).Error; err != nil {
			return fmt.Errorf("payment not found")
		}

		//  Check if the operation was already applied
		var op PaymentOperation
		if err := tx.First(&op, "payment_id = ? AND operation_id = ?", paymentID, operationID).Error; err == nil {
			// Already processed, return success (idempotent)
			return nil
		}

		//  Apply operation (update state)
		if err := p.ApplyOperation(operationID, operation); err != nil {
			return err
		}

		//  Record operation for idempotency
		newOp := PaymentOperation{
			PaymentID:   paymentID,
			OperationID: operationID,
			Operation:   string(operation),
			Result:      "success",
		}
		if err := tx.Create(&newOp).Error; err != nil {
			return err
		}

		//  Save updated payment state
		if err := tx.Save(&p).Error; err != nil {
			return err
		}

		return nil
	})
}
