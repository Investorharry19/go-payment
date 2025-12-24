package payment

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrBankTemporary = errors.New("temporary bank error")

type MockBank struct {
	FailureRate float64
	Delay       time.Duration
}

func (b *MockBank) simulate() error {
	time.Sleep(b.Delay)

	if rand.Float64() < b.FailureRate {
		return ErrBankTemporary
	}
	return nil
}

func (b *MockBank) Authorize(id string, amount int64) error {
	return b.simulate()
}

func (b *MockBank) Capture(id string, amount int64) error {
	return b.simulate()
}
func (b *MockBank) Void(id string, amount int64) error {
	return b.simulate()
}
func (b *MockBank) Refund(id string, amount int64) error {
	return b.simulate()
}

func (p *Payment) ApplyWithBank(bank Bank, opId string, operation Operation) error {
	// if res, ok := p.Operations[opId]; ok {
	// 	p.State = res.State
	// 	return nil
	// }

	// call bank
	var err error

	switch operation {
	case OPAuthorize:
		err = bank.Authorize(p.ID, p.Amount)
	case OPCapture:
		err = bank.Capture(p.ID, p.Amount)
	case OPVoid:
		err = bank.Void(p.ID, p.Amount)
	case OPRefund:
		err = bank.Refund(p.ID, p.Amount)
	default:
		err = fmt.Errorf("unknown operation")
	}

	if err != nil {
		return err
	}
	return p.ApplyOperation(opId, operation)

}
