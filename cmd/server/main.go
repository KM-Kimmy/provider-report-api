package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    config "provider-report-api/configs" 
    // "provider-report-api/internal/middleware"
    providerControllers "provider-report-api/internal/modules/provider-detail/controllers"
    providerRepositories "provider-report-api/internal/modules/provider-detail/repositories"
    providerServices "provider-report-api/internal/modules/provider-detail/services"
    
    _ "provider-report-api/cmd/docs"
    swaggerfiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Provider Detail Report API
// @version 1.0
// @description API for Provider Detail Report Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

    // Initialize controllers
    providerController := providerControllers.NewProviderController(providerService)
    templateController := providerControllers.NewTemplateController(templateService)
    scheduleController := providerControllers.NewScheduleController(scheduleService)
    logController := providerControllers.NewLogController(logService)
    fieldController := providerControllers.NewFieldController(fieldRepo)

    // Setup router
    router := gin.Default()

    // Add middleware
    router.Use(gin.Recovery())
    
    // หากต้องการใช้ CORS และ Logger สามารถเพิ่ม middleware เหล่านี้ได้:
    // router.Use(cors.Default()) // ต้อง import "github.com/gin-contrib/cors"
    // router.Use(gin.Logger())   // ใช้ gin's built-in logger

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        if err := config.HealthCheck(db); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    // Swagger endpoint
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

    // Setup API routes - ย้าย routing logic ไปยัง controllers
    api := router.Group("/api/v1")
    {
        // Setup all routes through controllers
        providerController.SetupRoutes(api)
        templateController.SetupRoutes(api)
        scheduleController.SetupRoutes(api)
        logController.SetupRoutes(api)
        fieldController.SetupRoutes(api)
    }

    port := cfg.ServerPort
    if port == "" {
        port = "8777"
    }

    log.Printf("Server starting on port %s", port)
    log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}