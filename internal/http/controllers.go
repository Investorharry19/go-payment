package http

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/gofiber/fiber/v2"
)

func CreatePaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
	var body struct {
		ID      string `json:"id"`
		Amount  int64  `json:"amount"`
		Email   string `json:"email"`
		UserId  string `json:"user_id"`
		OrderId string `json:"order_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	req := payment.AuthorizeRequest{
		PaymentID:   body.ID,
		Amount:      body.Amount,
		Email:       body.Email,
		CallbackURL: "http://localhost:8080/v1/payments/callback/verify",
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
	verifyResp, err := bank.Verify(c.Context(), reference)
	if err != nil {
		return renderHTML(c, "Payment verification failed", false)
	}

	// STEP 3: Load payment from DB
	_, err = store.Get(reference)
	if err != nil {
		return renderHTML(c, "Payment not found", false)
	}

	// Idempotency key for verification
	opID := "verify-" + reference

	// STEP 4: Apply state transition
	var operation payment.Operation

	if verifyResp.Status == "success" {
		operation = payment.OPCapture
	} else {
		operation = payment.OPVoid
	}

	if err := store.Apply(
		bank,
		reference,
		opID,
		operation,
	); err != nil {
		return renderHTML(c, "Failed to update payment state", false)
	}

	// STEP 5: Final HTML response
	return renderHTML(c, "Payment successful 🎉", true)
}

func GetAllPaymentsController(c *fiber.Ctx, store *payment.PaymentStoreDB) error {
	var payments []payment.Payment

	// Preload operations for each payment
	if err := store.DB.Preload("Operations").Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payments)
}

func GetPaymentByIdController(c *fiber.Ctx, store *payment.PaymentStoreDB) error {
	id := c.Params("id")
	p, err := store.Get(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
	}
	return c.JSON(p)
}

func RefundPaymentController(c *fiber.Ctx, store *payment.PaymentStoreDB, bank payment.Bank) error {
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
		payment.OPRefund,
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
	var operation payment.Operation
	if verifyResp.Status == "success" {
		operation = payment.OPCapture
	} else {
		operation = payment.OPVoid
	}

	// 6 Apply operation (idempotently)
	opID := "webhook-" + paymentID
	if err := store.Apply(bank, paymentID, opID, operation); err != nil {
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
