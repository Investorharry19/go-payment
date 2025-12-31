package payment

import (
	"context"
)

type Bank interface {
	Authorize(ctx context.Context, req AuthorizeRequest) (AuthorizeResponse, error)
	Verify(ctx context.Context, reference string) (VerifyResponse, error)
	Refund(ctx context.Context, req RefundRequest) (RefundResponse, error)
}

type AuthorizeRequest struct {
	PaymentID   string
	OperationID string // idempotency key
	Amount      int64
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
	Status    string // success, failed
	Amount    int64
	Currency  string
}

type RefundRequest struct {
	Reference string
	Amount    int64
}

type RefundResponse struct {
	Reference string
	Status    string
}
