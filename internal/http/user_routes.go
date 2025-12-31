package http

// func RegisterUserRoutes(app *fiber.App, store *payment.PaymentStoreDB, bank payment.Bank) {

// 	paymentRouters := app.Group("/v1/payments")
// 	// Create payment
// 	paymentRouters.Post("/", func(c *fiber.Ctx) error {
// 		return CreatePaymentController(c, store, bank)
// 	})

// 	// Get all payments
// 	paymentRouters.Get("/", func(c *fiber.Ctx) error {
// 		return GetAllPaymentsController(c, store)
// 	})

// 	// Get payment by ID
// 	paymentRouters.Get("/:id", func(c *fiber.Ctx) error {
// 		return GetPaymentByIdController(c, store)
// 	})

// 	// Operate payment (idempotent)

// 	//
// 	paymentRouters.Post("/:id/authorizations", func(c *fiber.Ctx) error {
// 		return AuthrizePaymentController(c, store, bank)
// 	})
// 	//
// 	paymentRouters.Post("/:id/capture", func(c *fiber.Ctx) error {
// 		return CapturePaymentController(c, store, bank)
// 	})
// 	//
// 	paymentRouters.Post("/:id/void", func(c *fiber.Ctx) error {
// 		return VoidPaymentController(c, store, bank)
// 	})
// 	//
// 	paymentRouters.Post("/:id/refund", func(c *fiber.Ctx) error {
// 		return RefundPaymentController(c, store, bank)
// 	})
// }
