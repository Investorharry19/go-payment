package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Investorharry19/go-payment/docs"
	_ "github.com/Investorharry19/go-payment/docs" // import generated docs
	"github.com/Investorharry19/go-payment/internal/http"
	"github.com/Investorharry19/go-payment/internal/payment"
	"github.com/Investorharry19/go-payment/internal/paystack"
	"github.com/Investorharry19/go-payment/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Payment Gateway API
// @version 1.0
// @description This is a mock Payment Gateway API.
// @contact.name API Support
// @contact.email amehharrison202017@gmail.com
// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter your JWT token in the format: "Bearer <your_token>"
// @BasePath /
func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
	if err := http.LoadPrivateKey(); err != nil {
		panic(err)
	}
	if err := middlewares.LoadPublicKey(); err != nil {
		panic(err)
	}

	app := fiber.New()

	env := os.Getenv("ENV")
	switch env {
	case "PROD":
		docs.SwaggerInfo.Host = "harrison-go-payment-microservice.up.railway.app"
	default:
		docs.SwaggerInfo.Host = "localhost:8080"
	}

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	db, err := payment.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	store := payment.NewPaymentStoreDB(db)
	bank := paystack.NewPaystackClient(os.Getenv("PAYSTACK_SECRET_KEY"))

	http.RegisterPaymentRoutes(app, store, bank)
	http.RegisterUserRoutes(app)

	// Run migrations at startup

	go func() {
		if err := db.AutoMigrate(&payment.Payment{}, &payment.PaymentOperation{}); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migrations completed!")
	}()
	fmt.Println("Connected to database successfully!")

	fmt.Println("Server started on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}

}
