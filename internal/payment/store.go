package payment

// type PaymentStore struct {
// 	Mu       sync.RWMutex
// 	Payments map[string]*StoredPayment
// }

// type StoredPayment struct {
// 	Mu      sync.Mutex
// 	Payment *Payment
// }

// func NewPaymentStore() *PaymentStore {
// 	return &PaymentStore{
// 		Payments: make(map[string]*StoredPayment),
// 	}
// }

// func (s *PaymentStore) Create(id string, amount int64) *Payment {
// 	s.Mu.Lock()
// 	defer s.Mu.Unlock()
// 	p := NewPayment(id, amount)
// 	s.Payments[id] = &StoredPayment{Payment: p}
// 	return p
// }

// func (s *PaymentStore) Get(id string) (*StoredPayment, bool) {
// 	s.Mu.Lock()
// 	defer s.Mu.Unlock()

// 	p, ok := s.Payments[id]
// 	return p, ok
// }

// func (s *PaymentStore) Apply(
// 	bank Bank,
// 	paymentId string,
// 	operationId string,
// 	operation Operation,
// ) error {
// 	stored, ok := s.Get(paymentId)
// 	if !ok {
// 		return fmt.Errorf("payment not found")
// 	}

// 	stored.Mu.Lock()
// 	defer stored.Mu.Unlock()

// 	return stored.Payment.ApplyOperation(operationId, operation)
// }
