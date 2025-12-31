package payment

// func TestValidTransitions(t *testing.T) {
// 	p := &Payment{
// 		ID:     "p1",
// 		Amount: 1000,
// 		State:  Initiated,
// 	}

// 	// test valid transitions
// 	if err := p.Authorize(); err != nil {
// 		t.Fatal("failed to authorize")
// 	}
// 	if err := p.Capture(); err != nil {
// 		t.Fatal("failed to capture")
// 	}
// 	if err := p.Refund(); err != nil {
// 		t.Fatal("Failed to refund")
// 	}

// }

// func TestInvalidTransitons(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		p    Payment
// 		fn   func(*Payment) error
// 	}{
// 		{
// 			name: "try Capture without authorize",
// 			p:    Payment{State: Initiated},
// 			fn:   (*Payment).Capture,
// 		},
// 		{
// 			name: "try refund without capture",
// 			p:    Payment{State: Authorized},
// 			fn:   (*Payment).Refund,
// 		},
// 		{
// 			name: "try void after capture",
// 			p:    Payment{State: Captured},
// 			fn:   (*Payment).Void,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.fn(&tt.p); err != nil {
// 				t.Fatal("expected error and got nil")
// 			}
// 		})
// 	}

// }
