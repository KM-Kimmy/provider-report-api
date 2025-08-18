package main

import (
	"fmt"
	"log"
	"net/http"

	docs "provider-report-api/cmd/docs"
	config "provider-report-api/configs"
	router "provider-report-api/internal/routers"
	providerRepositories "provider-report-api/internal/modules/provider-detail/repositories"
	providerServices "provider-report-api/internal/modules/provider-detail/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Provider Detail Report API
// @version 1.0
// @description API for Provider Detail Report Management System
// @termsOfService http://swagger.io/terms/

// @host localhost:8777
// @BasePath /tpa-api/report
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// corsMiddleware handles CORS settings for the API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow specific origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		// Allow specific headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		// Allow credentials
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.Initialize(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	providerRepo := providerRepositories.NewProviderRepository(db)
	templateRepo := providerRepositories.NewTemplateRepository(db)
	scheduleRepo := providerRepositories.NewScheduleRepository(db)
	logRepo := providerRepositories.NewLogRepository(db)
	fieldRepo := providerRepositories.NewFieldRepository(db)

	// Initialize services
	emailService := providerServices.NewEmailService(cfg)
	exportService := providerServices.NewExportService()
	providerService := providerServices.NewProviderService(providerRepo, exportService, fieldRepo)
	templateService := providerServices.NewTemplateService(templateRepo, fieldRepo)
	scheduleService := providerServices.NewScheduleService(scheduleRepo, templateRepo, emailService)
	logService := providerServices.NewLogService(logRepo)

	// Create dependencies struct
	deps := &router.Dependencies{
		ProviderService: providerService,
		TemplateService: templateService,
		ScheduleService: scheduleService,
		LogService:      logService,
		FieldRepo:       fieldRepo,
	}

	// Setup router
	r := gin.Default()

	// Use CORS middleware
	r.Use(corsMiddleware())

	// Custom middleware to log errors
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				fmt.Printf("Error: %v", e)
			}
		}
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := config.HealthCheck(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Validator setup
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validation functions or tags here
		_ = v // Use the validator as needed
	}

	// Define base settings for Swagger API documentation. This setup includes specifying the title,
	// description, and base path of the Swagger documentation, which enhances the discoverability
	// and usability of the API by developers.
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "TPA-API Report API Document"
	docs.SwaggerInfo.Description = "This is a list of all Report's API paths, requests, and responses."
	docs.SwaggerInfo.BasePath = "/tpa-api/report"

	// Pass dependencies to the router handler. All routes for the application modules are
	// defined within the router.InitializeRoutes function, which organizes and centralizes route
	// management, making the application easier to extend and maintain.
	router.InitializeRoutes(r, deps)

	// Get port from config
	port := cfg.ServerPort
	if port == "" {
		port = "8777"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)

	// Start the server
	r.Run(":" + port)
}