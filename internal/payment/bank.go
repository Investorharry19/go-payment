package payment

import (
	"context"
	"fmt"
	"strings"
)

type BankStatus string

const (
	BankStatusSuccess BankStatus = "success"
	BankStatusFailed  BankStatus = "failed"
	BankStatusPending BankStatus = "pending"
	BankStatusUnknown BankStatus = "unknown"
)

type Bank interface {
	Authorize(ctx context.Context, req AuthorizeRequest) (AuthorizeResponse, error)
	Verify(ctx context.Context, reference string) (VerifyResponse, error)
	Refund(ctx context.Context, req RefundRequest) (RefundResponse, error)
}

type AuthorizeRequest struct {
	PaymentID   string
	OperationID string // idempotency key
	Amount      int64  // Amount is expressed in the currency's minor unit.
	Currency    string
	Email       string
	CallbackURL string
}

type AuthorizeResponse struct {
	Reference        string
	AuthorizationURL string
}

type VerifyResponse struct {
	Reference string
	Status    BankStatus
	Amount    int64
	Currency  string
}

func ValidateVerification(p Payment, result VerifyResponse) error {
	if result.Reference != p.ID {
		return fmt.Errorf("bank reference does not match payment")
	}
	if result.Amount != p.Amount {
		return fmt.Errorf("verified amount does not match payment")
	}
	if !strings.EqualFold(result.Currency, p.Currency) {
		return fmt.Errorf("verified currency does not match payment")
	}
	return nil
}

type RefundRequest struct {
	Reference   string
	OperationID string
	Amount      int64
}

type RefundResponse struct {
	Reference string
	Status    BankStatus
}
