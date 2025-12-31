package http

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/gofiber/fiber/v2"
)

func RegisterPaymentRoutes(app *fiber.App, store *payment.PaymentStoreDB, bank payment.Bank) {

	paymentRouters := app.Group("/v1/payments")
	// Create payment
	paymentRouters.Post("/", func(c *fiber.Ctx) error {
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

	// Operate payment (idempotent)

	//
	// paymentRouters.Get("/callback/verify", func(c *fiber.Ctx) error {
	// 	return VerifyPaymentInCallbackController(c, store, bank)
	// })
	//

	//

	//
	paymentRouters.Post("/:id/refund", func(c *fiber.Ctx) error {
		return RefundPaymentController(c, store, bank)
	})

	// Webhook for Paystack events
	paymentRouters.Post("/webhooks/paystack", func(c *fiber.Ctx) error {
		fmt.Println("Paystack webhook received")
		return PaystackWebhookController(c, store, bank)
	})

	// 	// Step 1: Validate signature
	// 	signature := c.Get("X-Paystack-Signature")
	// 	if !validatePaystackSignature(body, signature, os.Getenv("PAYSTACK_SECRET")) {
	// 		return c.SendStatus(fiber.StatusUnauthorized)
	// 	}

	// 	// Step 2: Parse webhook event
	// 	var event PaystackWebhookEvent
	// 	if err := json.Unmarshal(body, &event); err != nil {
	// 		return c.SendStatus(fiber.StatusBadRequest)
	// 	}

	// 	// Only handle relevant events
	// 	if event.Event != "charge.success" && event.Event != "charge.failed" {
	// 		return c.SendStatus(fiber.StatusOK)
	// 	}

	// 	// Step 3: Lookup payment by reference
	// 	payment, ok := store.Get(event.Data.Reference)
	// 	if !ok {
	// 		// Payment not found — ignore but return 200
	// 		return c.SendStatus(fiber.StatusOK)
	// 	}

	// 	// Step 4: Verify with Paystack (optional but safer)
	// 	verifyResp, err := bank.Verify(c.Context(), event.Data.Reference)
	// 	if err != nil {
	// 		// Log error and exit, webhook will retry
	// 		fmt.Println("Paystack verify failed:", err)
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}

	// 	// Step 5: Update payment state
	// 	payment.Mu.Lock()
	// 	defer payment.Mu.Unlock()

	// 	opID := fmt.Sprintf("webhook-%s", event.Data.Reference)
	// 	operation := Operation(event.Event) // map "charge.success" → CAPTURE internally

	// 	if err := store.Apply(bank, payment.Payment.ID, opID, operation); err != nil {
	// 		fmt.Println("Payment store apply failed:", err)
	// 		// Return 500 so Paystack retries
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}

	// 	// Step 6: Respond 200 OK
	// 	return c.SendStatus(fiber.StatusOK)
	// })

}

func validatePaystackSignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
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
