package http

import (
	"fmt"

	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/Investorharry19/go-payment/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterPaymentRoutes(app *fiber.App, store *payment.PaymentStoreDB, bank payment.Bank) {

	paymentRouters := app.Group("/v1/payments")
	// Create payment

	paymentRouters.Post("/", middlewares.JWTMiddleware(), func(c *fiber.Ctx) error {
		return CreatePaymentController(c, store, bank)
	})

	// Get all payments
	paymentRouters.Get("/", func(c *fiber.Ctx) error {
		return GetAllPaymentsController(c, store)
	})

	// Get payment by ID
	paymentRouters.Get("/:id", func(c *fiber.Ctx) error {
		return GetPaymentByIdController(c, store)
	})

	// varify route
	paymentRouters.Get("/callback/verify", func(c *fiber.Ctx) error {
		return VerifyPaymentInCallbackController(c, store, bank)
	})

	//
	paymentRouters.Post("/:id/refund", func(c *fiber.Ctx) error {
		return RefundPaymentController(c, store, bank)
	})

	// Webhook for Paystack events
	paymentRouters.Post("/webhooks/paystack", func(c *fiber.Ctx) error {
		fmt.Println("Paystack webhook received")
		return PaystackWebhookController(c, store, bank)
	})

}

type PaystackWebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
	} `json:"data"`
}
