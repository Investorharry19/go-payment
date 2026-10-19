package payment

import "testing"

func TestValidateVerification(t *testing.T) {
	p := Payment{ID: "pay-1", Amount: 5000, Currency: "NGN"}

	tests := []struct {
		name    string
		result  VerifyResponse
		wantErr bool
	}{
		{
			name:   "matching payment",
			result: VerifyResponse{Reference: "pay-1", Amount: 5000, Currency: "ngn"},
		},
		{
			name:    "reference mismatch",
			result:  VerifyResponse{Reference: "other", Amount: 5000, Currency: "NGN"},
			wantErr: true,
		},
		{
			name:    "amount mismatch",
			result:  VerifyResponse{Reference: "pay-1", Amount: 4999, Currency: "NGN"},
			wantErr: true,
		},
		{
			name:    "currency mismatch",
			result:  VerifyResponse{Reference: "pay-1", Amount: 5000, Currency: "USD"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVerification(p, tt.result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateVerification() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
