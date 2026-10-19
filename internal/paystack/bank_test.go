package paystack

import (
	"testing"

	"github.com/Investorharry19/go-payment/internal/payment"
)

func TestNormalizeBankStatus(t *testing.T) {
	tests := map[string]payment.BankStatus{
		"success":    payment.BankStatusSuccess,
		"processed":  payment.BankStatusSuccess,
		"failed":     payment.BankStatusFailed,
		"abandoned":  payment.BankStatusFailed,
		"pending":    payment.BankStatusPending,
		"processing": payment.BankStatusPending,
		"future":     payment.BankStatusUnknown,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := normalizeBankStatus(input); got != want {
				t.Fatalf("normalizeBankStatus(%q) = %q, want %q", input, got, want)
			}
		})
	}
}
