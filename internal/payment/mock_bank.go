package payment

import (
	"errors"
	"math/rand/v2"
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
