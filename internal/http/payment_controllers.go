package http

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/gofiber/fiber/v2"
)

// PaymentRequest represents the JSON body for creating a payment
type PaymentRequest struct {
	ID       string `json:"id" example:"pay_12345"`
	Amount   int64  `json:"amount" example:"5000"`
	Currency string `json:"currency" example:"NGN"`
	Email    string `json:"email" example:"customer@example.com"`
	UserId   string `json:"user_id" example:"user_123"`
	OrderId  string `json:"order_id" example:"order_123"`
}

// PaymentResponse represents the JSON response after creating a payment
type PaymentResponse struct {
	StripeResponse interface{} `json:"stripe_response"`
	SQLResponse    interface{} `json:"sql_response"`
}
type OperationResponse struct {
	ID        string `json:"id" example:"op_123"`
	Type      string `json:"type" example:"authorization"`
	Amount    int64  `json:"amount" example:"5000"`
	CreatedAt string `json:"created_at" example:"2025-12-31T05:23:02Z"`
}

// PaymentFullResponse represents a payment with operations
type PaymentFullResponse struct {
	ID         string              `json:"id" example:"pay_12345"`
	Amount     int64               `json:"amount" example:"5000"`
	Currency   string              `json:"currency" example:"NGN"`
	UserId     string              `json:"user_id" example:"user_123"`
	OrderId    string              `json:"order_id" example:"order_123"`
	Operations []OperationResponse `json:"operations"`
}
type PaymentResponseByID struct {
	ID         string              `json:"id" example:"pay_12345"`
	Amount     int64               `json:"amount" example:"5000"`
	Currency   string              `json:"currency" example:"NGN"`
	UserId     string              `json:"user_id" example:"user_123"`
	OrderId    string              `json:"order_id" example:"order_123"`
	Operations []OperationResponse `json:"operations"`
}

// RefundRequest represents the JSON body for a refund
type RefundRequest struct {
	OperationID string `json:"operation_id" example:"op_12345"`
}

// PaymentRefundResponse represents a payment after refund (reuse existing payment struct)
type PaymentRefundResponse = PaymentFullResponse

// controllers

// CreatePaymentController godoc
// @Summary Create a payment
// @Description Creates a new payment and stores it in the database. Requires JWT authentication.
// @Tags Payments
// @Accept json
// @Produce json
// @Param payment body PaymentRequest true "Payment details"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /v1/payments [post]
func CreatePaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	var body struct {
		ID       string `json:"id"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Email    string `json:"email"`
		UserId   string `json:"user_id"`
		OrderId  string `json:"order_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	body.Currency = strings.ToUpper(strings.TrimSpace(body.Currency))
	if len(body.Currency) != 3 {
		return c.Status(400).JSON(fiber.Map{"error": "currency must be a 3-letter code"})
	}
	for _, char := range body.Currency {
		if char < 'A' || char > 'Z' {
			return c.Status(400).JSON(fiber.Map{"error": "currency must be a 3-letter code"})
		}
	}

	env := os.Getenv("ENV")
	callbackUrl := ""
	switch env {
	case "PROD":
		callbackUrl = "harrison-go-payment-microservice.up.railway.app"
	default:
		callbackUrl = "http://localhost:8080"
	}
	fmt.Println(callbackUrl)
	req := payment.AuthorizeRequest{
		PaymentID:   body.ID,
		Amount:      body.Amount,
		Currency:    body.Currency,
		Email:       body.Email,
		CallbackURL: fmt.Sprintf("%v/v1/payments/callback/verify", callbackUrl),
		OperationID: "op-" + body.ID,
	}
	fmt.Println(req.Email)
	resp, err := bank.Authorize(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Use resp.Reference and resp.AuthorizationURL as needed
	p, err := store.Create(
		body.ID,
		body.Amount,
		body.Currency,
		body.UserId,
		body.OrderId,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(struct {
		StripeResponse interface{} `json:"stripe_response"`
		SQLResponse    interface{} `json:"sql_response"`
	}{resp, p})
}

func VerifyPaymentInCallbackController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {

	reference := c.Query("reference")
	if reference == "" {
		return c.Status(400).SendString("Invalid payment reference")
	}

	// STEP 1: Show verifying page immediately
	// (Optional UX improvement: stream HTML)
	// We'll do full flow first

	// STEP 2: Verify with Paystack
	p, err := store.Get(reference)
	if err != nil {
		return renderHTML(c, "Payment not found", false)
	}

	verifyResp, err := bank.Verify(c.Context(), reference)
	if err != nil {
		return renderHTML(c, "Payment verification failed", false)
	}
	if err := payment.ValidateVerification(*p, verifyResp); err != nil {
		return renderHTML(c, "Payment verification did not match", false)
	}
	if verifyResp.Status != payment.BankStatusSuccess {
		return renderHTML(c, "Payment is not confirmed", false)
	}

	// Idempotency key for verification
	opID := "verify-" + reference

	// STEP 4: Apply state transition
	operation := payment.OPCapture

	if err := store.Apply(
		reference,
		opID,
		operation,
	); err != nil {
		return renderHTML(c, "Failed to update payment state", false)
	}

	// STEP 5: Final HTML response
	return renderHTML(c, "Payment successful 🎉", true)
}

// GetAllPaymentsController godoc
// @Summary Get all payments
// @Description Retrieves all payments with their associated operations
// @Tags Payments
// @Produce json
// @Success 200 {array} PaymentFullResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /v1/payments [get]
func GetAllPaymentsController(c *fiber.Ctx, store *payment.PaymentStoreDB) error {
	var payments []payment.Payment

	// Preload operations for each payment
	if err := store.DB.Preload("Operations").Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payments)
}

// GetPaymentByIdController godoc
// @Summary Get a payment by ID
// @Description Retrieves a single payment by its ID, including operations
// @Tags Payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} PaymentResponseByID
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /v1/payments/{id} [get]
func GetPaymentByIdController(c *fiber.Ctx, store *payment.PaymentStoreDB) error {
	id := c.Params("id")
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}
	return c.JSON(p)
}

// RefundPaymentController godoc
// @Summary Refund a payment
// @Description Applies a refund operation to a specific payment by ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param refund body RefundRequest true "Refund operation details"
// @Success 200 {object} PaymentRefundResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /v1/payments/{id}/refund [post]
func RefundPaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	id := c.Params("id")

	var body struct {
		OperationID string `json:"operation_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if strings.TrimSpace(body.OperationID) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "operation_id is required"})
	}
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}
	if p.State != payment.Captured {
		return c.Status(400).JSON(fiber.Map{"error": "only captured payments can be refunded"})
	}
	refund, err := bank.Refund(c.Context(), payment.RefundRequest{
		Reference:   p.ID,
		OperationID: p.ID + ":" + body.OperationID,
		Amount:      p.Amount,
	})
	if err != nil {
		return c.Status(502).JSON(fiber.Map{"error": err.Error()})
	}
	if refund.Status != payment.BankStatusSuccess || refund.Reference != p.ID {
		return c.Status(409).JSON(fiber.Map{"error": "refund is not confirmed"})
	}

	// Apply operation via DB-backed store
	if err := store.Apply(
		id,
		body.OperationID,
		payment.OPRefund,
	); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch updated payment with operations
	updated, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	return c.JSON(updated)

}

func PaystackWebhookController(
	c *fiber.Ctx,
	store *payment.PaymentStoreDB,
	bank payment.Bank,
) error {
	// Get raw body for signature verification
	body := c.Body()
	signature := c.Get("x-paystack-signature")
	if signature == "" {
		return c.Status(400).SendString("Missing signature")
	}

	secret := os.Getenv("PAYSTACK_SECRET_KEY")
	if secret == "" {
		return c.Status(500).SendString("Server config error")
	}

	h := hmac.New(sha512.New, []byte(secret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return c.Status(400).SendString("Invalid signature")
	}

	// 2 Parse the webhook event
	var event struct {
		Event string `json:"event"`
		Data  struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	if err := c.BodyParser(&event); err != nil {
		return c.Status(400).SendString("Invalid JSON")
	}

	paymentID := event.Data.Reference

	// 3 Retrieve the payment from DB
	stored, err := store.Get(paymentID)
	if err != nil {
		fmt.Printf("Payment not found: %s\n", paymentID)
		return c.SendStatus(fiber.StatusOK) // acknowledge webhook
	}
	_ = stored

	// 4 Verify with Paystack API
	verifyResp, err := bank.Verify(c.Context(), paymentID)
	if err != nil {
		fmt.Printf("Paystack verify failed for %s: %v\n", paymentID, err)
		return c.SendStatus(fiber.StatusOK)
	}

	// 5 Determine operation based on verification
	if err := payment.ValidateVerification(*stored, verifyResp); err != nil {
		fmt.Printf("Paystack verification mismatch for %s: %v\n", paymentID, err)
		return c.SendStatus(fiber.StatusOK)
	}
	if verifyResp.Status != payment.BankStatusSuccess {
		return c.SendStatus(fiber.StatusOK)
	}
	operation := payment.OPCapture

	// 6 Apply operation (idempotently)
	opID := "webhook-" + paymentID
	if err := store.Apply(paymentID, opID, operation); err != nil {
		fmt.Printf("Failed to apply operation for %s: %v\n", paymentID, err)
	}

	// 7 Respond 200 OK to Paystack
	return c.SendString("OK")
}

func renderHTML(c *fiber.Ctx, message string, success bool) error {
	status := "failed"
	color := "red"

	if success {
		status = "success"
		color = "green"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<title>Payment Status</title>
</head>
<body style="font-family: sans-serif;">
	<h1 style="color:%s;">Payment %s</h1>
	<p>%s</p>
</body>
</html>
`, color, status, message)

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

/*
func AuthrizePaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	id := c.Params("id")

	var body struct {
		OperationID string `json:"operation_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	fmt.Println("operation:", "operation_id:", body.OperationID)

	// Apply operation via DB-backed store
	err := store.Apply(
		bank,
		id,
		body.OperationID,
		payment.OPAuthorize,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch updated payment with operations
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	return c.JSON(p)

}



func VoidPaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	id := c.Params("id")

	var body struct {
		OperationID string `json:"operation_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	fmt.Println("operation:", "operation_id:", body.OperationID)

	// Apply operation via DB-backed store
	err := store.Apply(
		bank,
		id,
		body.OperationID,
		payment.OPVoid,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch updated payment with operations
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	return c.JSON(p)

}

func CapturePaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	id := c.Params("id")

	var body struct {
		OperationID string `json:"operation_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	fmt.Println("operation:", "operation_id:", body.OperationID)

	// Apply operation via DB-backed store
	err := store.Apply(
		bank,
		id,
		body.OperationID,
		payment.OPCapture,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch updated payment with operations
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}

	return c.JSON(p)

}
*/
