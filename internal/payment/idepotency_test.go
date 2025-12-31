package payment

// func TestIdepotencCapture(t *testing.T) {
// 	p := NewPayment("p1", 1000)

// 	opAuth := "op-a1"
// 	opCap := "op-a2"

// 	if err := p.ApplyOperation(opAuth, OPAuthorize); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := p.ApplyOperation(opCap, OPCapture); err != nil {
// 		t.Fatal(err)
// 	}

// 	// retry capture
// 	if err := p.ApplyOperation(opCap, OPCapture); err != nil {
// 		t.Fatal(err)
// 	}

// 	if p.State != Captured {
// 		t.Fatalf("Expected capture, got %s", p.State)
// 	}
// 	t.Log(p)
// }

// func TestIdepotencRefund(t *testing.T) {
// 	p := NewPayment("p2", 1000)

// 	opAuth := "op-a1"
// 	opCap := "op-a2"
// 	opRef := "op-ref"
// 	if err := p.ApplyOperation(opAuth, OPAuthorize); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := p.ApplyOperation(opCap, OPCapture); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := p.ApplyOperation(opRef, OPRefund); err != nil {
// 		t.Fatal(err)
// 	}

// 	// retry refund
// 	if err := p.ApplyOperation(opRef, OPRefund); err != nil {
// 		t.Fatal(err)
// 	}

// 	if p.State != Refunded {
// 		t.Fatalf("Expected refund, got %s", p.State)
// 	}
// 	t.Log(p)
// }
