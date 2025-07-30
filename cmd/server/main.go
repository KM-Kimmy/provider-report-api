package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    config "provider-report-api/configs" 
    "provider-report-api/internal/middleware"
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
    router.Use(middleware.CORS())
    router.Use(middleware.Logger())
    router.Use(gin.Recovery())

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

    // API routes
    api := router.Group("/api/v1")
    {
        // Provider Detail Report routes
        providers := api.Group("/providers")
        {
            // CRUD Operations (ตาม Swagger spec)
            providers.POST("", providerController.CreateProvider)           // Create new provider
            providers.GET("/:id", providerController.GetProviderByID)       // Get provider by ID
            providers.PUT("/:id", providerController.UpdateProvider)        // Update provider
            providers.DELETE("/:id", providerController.DeleteProvider)     // Delete provider
            
            // Search and Filter Operations
            providers.GET("/search", providerController.SearchProviders)
            
            // Report Operations
            providers.POST("/report", providerController.GenerateReport)
            providers.POST("/export", providerController.ExportReport)
            
            // Summary and Statistics
            providers.GET("/summary", providerController.GetProviderSummary)
            providers.GET("/stats", providerController.GetProviderStats)
            
            // Reference Data
            providers.GET("/provinces", providerController.GetProvinces)
            providers.GET("/types", providerController.GetProviderTypes)
        }

        // Template Management routes
        templates := api.Group("/templates")
        {
            templates.GET("", templateController.GetTemplates)
            templates.GET("/:id", templateController.GetTemplate)
            templates.POST("", templateController.CreateTemplate)
            templates.PUT("/:id", templateController.UpdateTemplate)
            templates.DELETE("/:id", templateController.DeleteTemplate)
        }

        // Schedule Management routes
        schedules := api.Group("/schedules")
        {
            schedules.GET("", scheduleController.GetSchedules)
            schedules.GET("/:id", scheduleController.GetSchedule)
            schedules.POST("", scheduleController.CreateSchedule)
            schedules.PUT("/:id", scheduleController.UpdateSchedule)
            schedules.DELETE("/:id", scheduleController.DeleteSchedule)
            schedules.POST("/:id/run", scheduleController.RunSchedule)
        }

        // Log Report routes
        logs := api.Group("/logs")
        {
            logs.GET("/sent-reports", logController.GetSentReportLogs)
            logs.GET("/sent-reports/:id", logController.GetSentReportLog)
        }

        // Available Fields routes
        fields := api.Group("/fields")
        {
            fields.GET("", fieldController.GetAvailableFields)
            fields.GET("/by-category/:category", fieldController.GetFieldsByCategory)
        }
    }

    port := cfg.ServerPort
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}