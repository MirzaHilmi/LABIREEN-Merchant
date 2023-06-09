package main

import (
	"labireen-merchant/config"
	"labireen-merchant/handlers"
	"labireen-merchant/pkg/mail"
	"labireen-merchant/repositories"
	"labireen-merchant/routes"
	"labireen-merchant/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(".env file loading failed")
	}

	// Initialize database connection
	db, err := config.GetDB()
	if err != nil {
		log.Fatalln("Database initialization failed")
	}

	// Auto migrate entities
	if err := config.Migrate(db); err != nil {
		log.Fatalln("Auto Migration failed")
	}

	emailService := mail.NewGmailSender(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))
	authService := services.NewAuthService(repositories.NewMerchantRepository(db))
	merchantService := services.NewMerchantService(repositories.NewMerchantRepository(db))

	authHandler := handlers.NewAuthHandler(authService, emailService)
	merchantHandler := handlers.NewMerchantHandler(merchantService)

	app := gin.Default()

	// Register auth routes
	authRoutes := routes.AuthRoutes{
		Router:      app,
		AuthHandler: authHandler,
	}
	authRoutes.Register()

	// Register merchant routes
	merchantRoutes := routes.MerchantRoutes{
		Router:          app,
		MerchantHandler: merchantHandler,
	}
	merchantRoutes.Register()

	app.Run(os.Getenv("APP_PORT"))
}
