package payment

// func TestConcurrentCapture(t *testing.T) {

// 	bank := &MockBank{
// 		FailureRate: 0.0,
// 	}
// 	store := NewPaymentStore()

// 	store.Create("p1", 1000)
// 	// authorize
// 	if err := store.Apply(bank, "p1", "auth", OPAuthorize); err != nil {
// 		t.Fatal(err)
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		_ = store.Apply(bank, "p1", "cap-1", OPAuthorize)
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		_ = store.Apply(bank, "p1", "cap-2", OPAuthorize)
// 	}()

// 	wg.Wait()

// 	stored, _ := store.Get("p1")
// 	if stored.Payment.State != Captured {
// 		t.Fatalf("expected captured, got : %v", stored.Payment.State)
// 	}
// }
