package payment

import (
	"fmt"
)

type State string

const (
	Initiated  State = "initiated"
	Authorized State = "authorized"
	Captured   State = "captured"
	Voided     State = "voided"
	Refunded   State = "refunded"
)

type Operation string

const (
	OPAuthorize Operation = "authorize"
	OPCapture   Operation = "capture"
	OPVoid      Operation = "void"
	OPRefund    Operation = "refund"
)

type OperationResult struct {
	Operation Operation
	State     State
}

var (
	ErrInvalidTranstion = fmt.Errorf("Invalid state transition")
)

func NewPayment(id string, amount int64) *Payment {
	return &Payment{
		ID:     id,
		Amount: amount,
		State:  Initiated,
	}
}

func (p *Payment) ApplyOperation(opID string, operation Operation) error {

	// if the operation is already execulted then return stored result
	// if res, ok := p.Operations[opID]; ok {
	// 	p.State = res.State
	// 	return nil
	// }

	// execute operation that did not exists
	var err error
	switch operation {
	case OPAuthorize:
		err = p.Authorize()
	case OPCapture:
		err = p.Capture()
	case OPVoid:
		err = p.Void()
	case OPRefund:
		err = p.Refund()

	default:
		return fmt.Errorf("unknowk operation: %s", operation)
	}

	if err != nil {
		return err
	}

	// record sucesssful operation
	// p.Operations[opID] = OperationResult{
	// 	Operation: operation,
	// 	State:     p.State,
	// }
	return nil
}

// rule you can only authorize an authorized initiated
func (p *Payment) Authorize() error {
	if p.State != Initiated {
		return fmt.Errorf("%w: cannot authorize from %s", ErrInvalidTranstion, p.State)
	}

	p.State = Authorized
	return nil
}

// rule you can only capture an authorized or initiated payment
func (p *Payment) Capture() error {
	if p.State != Authorized && p.State != Initiated {
		return fmt.Errorf("%w: cannot capture from %s", ErrInvalidTranstion, p.State)
	}

	p.State = Captured
	return nil
}

// rule you can only void an authorized payment
func (p *Payment) Void() error {
	if p.State != Authorized {
		return fmt.Errorf("%w: cannot void from %s", ErrInvalidTranstion, p.State)
	}

	p.State = Voided
	return nil
}

// rule you can only void an authorized payment
func (p *Payment) Refund() error {
	if p.State != Captured {
		return fmt.Errorf("%w: cannot refund from %s", ErrInvalidTranstion, p.State)
	}

	p.State = Refunded
	return nil
}

func PaymentFunction() {
	fmt.Println("Hello from payment")
}

/*

	payment structure for reference
	initiated
		|
	authorized
		|
	-		-
	|		|
   Void    Capture
			|
			Refund

*/
