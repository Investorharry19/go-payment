package main

import (
	"fmt"
	"log"

	"github.com/Investorharry19/go-payment/internal/http"
	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	db, err := payment.ConnectPostgres()
	if err != nil {

		panic(err)
	}
	store := payment.NewPaymentStoreDB(db)
	bank := payment.MockBank{FailureRate: 0}
	http.RegisterRoutes(app, store, &bank)

	// Run migrations at startup
	err = db.AutoMigrate(
		&payment.Payment{},
		&payment.PaymentOperation{},
	)
	fmt.Println("Connected to database successfully!")

	app.Listen(":8080")
}
