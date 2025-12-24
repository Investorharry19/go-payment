package http

import (
	"fmt"

	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, store *payment.PaymentStoreDB, bank payment.Bank) {
	// Create payment
	app.Post("/payments", func(c *fiber.Ctx) error {
		var body struct {
			ID     string `json:"id"`
			Amount int64  `json:"amount"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		p, err := store.Create(body.ID, body.Amount)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(p)
	})

	// Get all payments
	app.Get("/payments", func(c *fiber.Ctx) error {
		var payments []payment.Payment

		// Preload operations for each payment
		if err := store.DB.Preload("Operations").Find(&payments).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(payments)
	})

	// Get payment by ID
	app.Get("/payments/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		p, err := store.Get(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "payment not found"})
		}
		return c.JSON(p)
	})

	// Operate payment (idempotent)
	app.Post("/payments/:id/operate", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var body struct {
			Operation   string `json:"operation"`
			OperationID string `json:"operation_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		fmt.Println("operation:", body.Operation, "operation_id:", body.OperationID)

		// Apply operation via DB-backed store
		err := store.Apply(
			bank,
			id,
			body.OperationID,
			payment.Operation(body.Operation),
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
	})
}
