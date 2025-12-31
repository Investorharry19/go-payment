package payment

// func TestBankFailureDoesNotCHangeState(t *testing.T) {

// 	bank := &MockBank{
// 		FailureRate: 1.0,
// 	}

// 	p := NewPayment("p1", 1000)

// 	err := p.ApplyWithBank(bank, "op1", Operation(OPAuthorize))

// 	if err != nil {
// 		t.Fatal("expected bank error")
// 	}
// 	if p.State != Initiated {
// 		t.Fatalf("state changed on failure : %v", err)
// 	}
// }
// func TestBankREtrySucceeds(t *testing.T) {

// 	bank := &MockBank{
// 		FailureRate: 0.0,
// 	}

// 	p := NewPayment("p2", 1000)

// 	err := p.ApplyWithBank(bank, "op1", Operation(OPAuthorize))

// 	if err != nil {
// 		t.Fatal("expected bank error")
// 	}
// 	if p.State != Authorized {
// 		t.Fatalf("expected authorize, got: %v", p.State)
// 	}
// }
