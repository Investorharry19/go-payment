package payment

type Bank interface {
	Authorize(paymentId string, amount int64) error
	Capture(paymentId string, amount int64) error
	Void(paymentId string, amount int64) error
	Refund(paymentId string, amount int64) error
}
