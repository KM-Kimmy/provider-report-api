package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "provider-report-api/internal/modules/provider-detail/dtos"
    "provider-report-api/internal/modules/provider-detail/repositories"
    "provider-report-api/internal/modules/provider-detail/services"
)

// ================= PROVIDER CONTROLLER =================

type ProviderController struct {
    providerService *services.ProviderService
}

func NewProviderController(providerService *services.ProviderService) *ProviderController {
    return &ProviderController{
        providerService: providerService,
    }
}

func InitializeProviderDetailController(
	providerRoute *gin.RouterGroup,
	providerService *services.ProviderService,
	templateService *services.TemplateService,
	scheduleService *services.ScheduleService,
	logService *services.LogService,
	fieldRepo *repositories.FieldRepository,
) {
	// Initialize controllers with their dependencies
	providerController := NewProviderController(providerService)
	templateController := NewTemplateController(templateService)
	scheduleController := NewScheduleController(scheduleService)
	logController := NewLogController(logService)
	fieldController := NewFieldController(fieldRepo)

	providers := providerRoute.Group("/providers")
	{
		// Provider CRUD Operations
		providers.POST("", providerController.CreateProvider)           // Create new provider
		providers.GET("/:id", providerController.GetProvider)           // Get provider by ID
		providers.PUT("/:id", providerController.UpdateProvider)        // Update provider
		providers.DELETE("/:id", providerController.DeleteProvider)     // Delete provider
		
		// Provider Search and Filter Operations
		providers.GET("/search", providerController.SearchProviders)
		
		// Provider Report Operations
		providers.POST("/report", providerController.GenerateReport)
		providers.POST("/export", providerController.ExportReport)
		
		// Provider Summary and Statistics
		providers.GET("/summary", providerController.GetProviderSummary)
		providers.GET("/stats", providerController.GetProviderStats)
		
		// Provider Reference Data
		providers.GET("/provinces", providerController.GetProvinces)
		providers.GET("/types", providerController.GetProviderTypes)
	}

	// Template Routes
	templates := providerRoute.Group("/templates")
	{
		templates.GET("", templateController.GetTemplates)
		templates.GET("/:id", templateController.GetTemplate)
		templates.POST("", templateController.CreateTemplate)
		templates.PUT("/:id", templateController.UpdateTemplate)
		templates.DELETE("/:id", templateController.DeleteTemplate)
	}

	// Schedule Routes
	schedules := providerRoute.Group("/schedules")
	{
		schedules.GET("", scheduleController.GetSchedules)
		schedules.GET("/:id", scheduleController.GetSchedule)
		schedules.POST("", scheduleController.CreateSchedule)
		schedules.PUT("/:id", scheduleController.UpdateSchedule)
		schedules.DELETE("/:id", scheduleController.DeleteSchedule)
		schedules.POST("/:id/run", scheduleController.RunSchedule)
	}

	// Log Routes
	logs := providerRoute.Group("/logs")
	{
		logs.GET("/sent-reports", logController.GetSentReportLogs)
		logs.GET("/sent-reports/:id", logController.GetSentReportLog)
	}

	// Field Routes
	fields := providerRoute.Group("/fields")
	{
		fields.GET("", fieldController.GetAvailableFields)
		fields.GET("/by-category/:category", fieldController.GetFieldsByCategory)
	}
}

// CreateProvider godoc
// @Summary Create a new provider
// @Description Create a new provider with the provided information
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param provider body dtos.CreateProviderRequestDTO true "Provider data"
// @Success 201 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers [post]
// @Security BearerAuth
func (c *ProviderController) CreateProvider(ctx *gin.Context) {
    var req dtos.CreateProviderRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    provider, err := c.providerService.CreateProvider(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to create provider",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, dtos.APIResponse{
        Success: true,
        Message: "Provider created successfully",
        Data:    provider,
    })
}

// GetProvider godoc
// @Summary Get provider by ID
// @Description Get a single provider by its ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/{id} [get]
// @Security BearerAuth
func (c *ProviderController) GetProvider(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid provider ID",
            Details: err.Error(),
        })
        return
    }

    provider, err := c.providerService.GetProviderByID(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
            Code:    http.StatusNotFound,
            Message: "Provider not found",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider retrieved successfully",
        Data:    provider,
    })
}

// UpdateProvider godoc
// @Summary Update provider
// @Description Update provider information by ID
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Param provider body dtos.UpdateProviderRequestDTO true "Provider update data"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/{id} [put]
// @Security BearerAuth
func (c *ProviderController) UpdateProvider(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid provider ID",
            Details: err.Error(),
        })
        return
    }

    var req dtos.UpdateProviderRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    provider, err := c.providerService.UpdateProvider(id, req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update provider",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider updated successfully",
        Data:    provider,
    })
}

// DeleteProvider godoc
// @Summary Delete provider
// @Description Delete provider by ID
// @Tags providerDetail 
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/{id} [delete]
// @Security BearerAuth
func (c *ProviderController) DeleteProvider(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid provider ID",
            Details: err.Error(),
        })
        return
    }

    err = c.providerService.DeleteProvider(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete provider",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider deleted successfully",
    })
}

// SearchProviders godoc
// @Summary Search providers
// @Description Search providers with filters
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param provider_name query string false "Provider name"
// @Param is_tpa_network query bool false "Is TPA Network"
// @Param province_name query string false "Province name"
// @Param provider_type query string false "Provider type"
// @Param business_type query string false "Business type"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dtos.PaginatedResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/search [get]
// @Security BearerAuth
func (c *ProviderController) SearchProviders(ctx *gin.Context) {
    var req dtos.ProviderSearchRequestDTO
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid query parameters",
            Details: err.Error(),
        })
        return
    }

    providers, total, err := c.providerService.SearchProviders(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to search providers",
            Details: err.Error(),
        })
        return
    }

    if req.Page == 0 {
        req.Page = 1
    }
    if req.Limit == 0 {
        req.Limit = 10
    }

    totalPages := int(total) / req.Limit
    if int(total)%req.Limit > 0 {
        totalPages++
    }

    response := dtos.PaginatedResponse{
        Data:       providers,
        Total:      total,
        Page:       req.Page,
        Limit:      req.Limit,
        TotalPages: totalPages,
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Providers retrieved successfully",
        Data:    response,
    })
}

// GetProviderSummary godoc
// @Summary Get provider summary
// @Description Get summary statistics of providers with optional filters
// @Tags providerDetail
// @Produce json
// @Param provider_name query string false "Provider name"
// @Param is_tpa_network query bool false "Is TPA Network"
// @Param province_name query string false "Province name"
// @Param provider_type query string false "Provider type"
// @Param business_type query string false "Business type"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/summary [get]
// @Security BearerAuth
func (c *ProviderController) GetProviderSummary(ctx *gin.Context) {
    var req dtos.ProviderSearchRequestDTO
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid query parameters",
            Details: err.Error(),
        })
        return
    }

    summary, err := c.providerService.GetProviderSummary(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get provider summary",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider summary retrieved successfully",
        Data:    summary,
    })
}

// GetProviderStats godoc
// @Summary Get provider statistics
// @Description Get comprehensive provider statistics
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/stats [get]
// @Security BearerAuth
func (c *ProviderController) GetProviderStats(ctx *gin.Context) {
    stats, err := c.providerService.GetProviderStats()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get provider statistics",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider statistics retrieved successfully",
        Data:    stats,
    })
}

// GenerateReport godoc
// @Summary Generate provider report
// @Description Generate a provider report based on specified criteria
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param report body dtos.ProviderReportRequestDTO true "Report generation parameters"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/report [post]
// @Security BearerAuth
func (c *ProviderController) GenerateReport(ctx *gin.Context) {
    var req dtos.ProviderReportRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    reportData, err := c.providerService.GenerateReport(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to generate report",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Report generated successfully",
        Data:    reportData,
    })
}

// ExportReport godoc
// @Summary Export provider report
// @Description Export provider report in specified format (Excel, PDF, etc.)
// @Tags providerDetail
// @Accept json
// @Produce application/octet-stream
// @Param export body dtos.ProviderReportRequestDTO true "Export parameters"
// @Success 200 {file} file "Exported report file"
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/export [post]
// @Security BearerAuth
func (c *ProviderController) ExportReport(ctx *gin.Context) {
    var req dtos.ProviderReportRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    if req.FormatType == "" {
        req.FormatType = "excel"
    }

    fileData, filename, contentType, err := c.providerService.ExportReport(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to export report",
            Details: err.Error(),
        })
        return
    }

    ctx.Header("Content-Disposition", "attachment; filename="+filename)
    ctx.Data(http.StatusOK, contentType, fileData)
}

// GetProvinces godoc
// @Summary Get all provinces
// @Description Get list of all available provinces
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/provinces [get]
// @Security BearerAuth
func (c *ProviderController) GetProvinces(ctx *gin.Context) {
    provinces, err := c.providerService.GetProvinces()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get provinces",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provinces retrieved successfully",
        Data:    provinces,
    })
}

// GetProviderTypes godoc
// @Summary Get all provider types
// @Description Get list of all available provider types
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/types [get]
// @Security BearerAuth
func (c *ProviderController) GetProviderTypes(ctx *gin.Context) {
    types, err := c.providerService.GetProviderTypes()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get provider types",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Provider types retrieved successfully",
        Data:    types,
    })
}

// ================= TEMPLATE CONTROLLER =================

type TemplateController struct {
    templateService *services.TemplateService
}

func NewTemplateController(templateService *services.TemplateService) *TemplateController {
    return &TemplateController{
        templateService: templateService,
    }
}

// SetupRoutes sets up all template-related routes
func (c *TemplateController) SetupRoutes(api *gin.RouterGroup) {
    templates := api.Group("/templates")
    {
        templates.GET("", c.GetTemplates)
        templates.GET("/:id", c.GetTemplate)
        templates.POST("", c.CreateTemplate)
        templates.PUT("/:id", c.UpdateTemplate)
        templates.DELETE("/:id", c.DeleteTemplate)
    }
}

// GetTemplates godoc
// @Summary Get all templates
// @Description Get list of all report templates
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/templates [get]
// @Security BearerAuth
func (c *TemplateController) GetTemplates(ctx *gin.Context) {
    templates, err := c.templateService.GetAllTemplates()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get templates",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Templates retrieved successfully",
        Data:    templates,
    })
}

// GetTemplate godoc
// @Summary Get template by ID
// @Description Get a single template by its ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/templates/{id} [get]
// @Security BearerAuth
func (c *TemplateController) GetTemplate(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid template ID",
            Details: err.Error(),
        })
        return
    }

    template, err := c.templateService.GetTemplate(id)
    if err != nil {
        if err.Error() == "template not found" {
            ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
                Code:    http.StatusNotFound,
                Message: "Template not found",
                Details: err.Error(),
            })
            return
        }

        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get template",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Template retrieved successfully",
        Data:    template,
    })
}

// CreateTemplate godoc
// @Summary Create a new template
// @Description Create a new report template
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param template body dtos.CreateTemplateRequestDTO true "Template data"
// @Success 201 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/templates [post]
// @Security BearerAuth
func (c *TemplateController) CreateTemplate(ctx *gin.Context) {
    var req dtos.CreateTemplateRequestDTO
    
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    template, err := c.templateService.CreateTemplate(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to create template",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, dtos.APIResponse{
        Success: true,
        Message: "Template created successfully",
        Data:    template,
    })
}

// UpdateTemplate godoc
// @Summary Update template
// @Description Update template information by ID
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param template body dtos.UpdateTemplateRequestDTO true "Template update data"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/templates/{id} [put]
// @Security BearerAuth
func (c *TemplateController) UpdateTemplate(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid template ID",
            Details: err.Error(),
        })
        return
    }

    var req dtos.UpdateTemplateRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    template, err := c.templateService.UpdateTemplate(id, req)
    if err != nil {
        if err.Error() == "template not found" {
            ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
                Code:    http.StatusNotFound,
                Message: "Template not found",
                Details: err.Error(),
            })
            return
        }

        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update template",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Template updated successfully",
        Data:    template,
    })
}

// DeleteTemplate godoc
// @Summary Delete template
// @Description Delete template by ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/templates/{id} [delete]
// @Security BearerAuth
func (c *TemplateController) DeleteTemplate(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid template ID",
            Details: err.Error(),
        })
        return
    }

    err = c.templateService.DeleteTemplate(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete template",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Template deleted successfully",
    })
}

// ================= SCHEDULE CONTROLLER =================

type ScheduleController struct {
    scheduleService *services.ScheduleService
}

func NewScheduleController(scheduleService *services.ScheduleService) *ScheduleController {
    return &ScheduleController{
        scheduleService: scheduleService,
    }
}

// SetupRoutes sets up all schedule-related routes
func (c *ScheduleController) SetupRoutes(api *gin.RouterGroup) {
    schedules := api.Group("/schedules")
    {
        schedules.GET("", c.GetSchedules)
        schedules.GET("/:id", c.GetSchedule)
        schedules.POST("", c.CreateSchedule)
        schedules.PUT("/:id", c.UpdateSchedule)
        schedules.DELETE("/:id", c.DeleteSchedule)
        schedules.POST("/:id/run", c.RunSchedule)
    }
}

// GetSchedules godoc
// @Summary Get all schedules
// @Description Get list of all report schedules
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules [get]
// @Security BearerAuth
func (c *ScheduleController) GetSchedules(ctx *gin.Context) {
    schedules, err := c.scheduleService.GetAllSchedules()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get schedules",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Schedules retrieved successfully",
        Data:    schedules,
    })
}

// GetSchedule godoc
// @Summary Get schedule by ID
// @Description Get a single schedule by its ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules/{id} [get]
// @Security BearerAuth
func (c *ScheduleController) GetSchedule(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid schedule ID",
            Details: err.Error(),
        })
        return
    }

    schedule, err := c.scheduleService.GetSchedule(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
            Code:    http.StatusNotFound,
            Message: "Schedule not found",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Schedule retrieved successfully",
        Data:    schedule,
    })
}

// CreateSchedule godoc
// @Summary Create a new schedule
// @Description Create a new report schedule
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param schedule body dtos.CreateScheduleRequestDTO true "Schedule data"
// @Success 201 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules [post]
// @Security BearerAuth
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
    var req dtos.CreateScheduleRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    schedule, err := c.scheduleService.CreateSchedule(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to create schedule",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, dtos.APIResponse{
        Success: true,
        Message: "Schedule created successfully",
        Data:    schedule,
    })
}

// UpdateSchedule godoc
// @Summary Update schedule
// @Description Update schedule information by ID
// @Tags providerDetail
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param schedule body dtos.UpdateScheduleRequestDTO true "Schedule update data"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules/{id} [put]
// @Security BearerAuth
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid schedule ID",
            Details: err.Error(),
        })
        return
    }

    var req dtos.UpdateScheduleRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    schedule, err := c.scheduleService.UpdateSchedule(id, req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update schedule",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Schedule updated successfully",
        Data:    schedule,
    })
}

// DeleteSchedule godoc
// @Summary Delete schedule
// @Description Delete schedule by ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules/{id} [delete]
// @Security BearerAuth
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid schedule ID",
            Details: err.Error(),
        })
        return
    }

    err = c.scheduleService.DeleteSchedule(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete schedule",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Schedule deleted successfully",
    })
}

// RunSchedule godoc
// @Summary Run schedule manually
// @Description Execute a schedule manually by ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/schedules/{id}/run [post]
// @Security BearerAuth
func (c *ScheduleController) RunSchedule(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid schedule ID",
            Details: err.Error(),
        })
        return
    }

    result, err := c.scheduleService.RunSchedule(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to run schedule",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Schedule executed successfully",
        Data:    result,
    })
}

// ================= LOG CONTROLLER =================

type LogController struct {
    logService *services.LogService
}

func NewLogController(logService *services.LogService) *LogController {
    return &LogController{
        logService: logService,
    }
}

// SetupRoutes sets up all log-related routes
func (c *LogController) SetupRoutes(api *gin.RouterGroup) {
    logs := api.Group("/logs")
    {
        logs.GET("/sent-reports", c.GetSentReportLogs)
        logs.GET("/sent-reports/:id", c.GetSentReportLog)
    }
}

// GetSentReportLogs godoc
// @Summary Get sent report logs
// @Description Get paginated list of sent report logs with optional filters
// @Tags providerDetail
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Log status"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/logs/sent-reports [get]
// @Security BearerAuth
func (c *LogController) GetSentReportLogs(ctx *gin.Context) {
    var req dtos.LogSearchRequestDTO
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid query parameters",
            Details: err.Error(),
        })
        return
    }

    // Set default pagination
    if req.Page == 0 {
        req.Page = 1
    }
    if req.Limit == 0 {
        req.Limit = 10
    }

    result, err := c.logService.GetSentReportLogs(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get sent report logs",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Sent report logs retrieved successfully",
        Data:    result,
    })
}

// GetSentReportLog godoc
// @Summary Get sent report log by ID
// @Description Get a single sent report log by its ID
// @Tags providerDetail
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/logs/sent-reports/{id} [get]
// @Security BearerAuth
func (c *LogController) GetSentReportLog(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid log ID",
            Details: err.Error(),
        })
        return
    }

    log, err := c.logService.GetSentReportLog(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
            Code:    http.StatusNotFound,
            Message: "Log not found",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Log retrieved successfully",
        Data:    log,
    })
}

// ================= FIELD CONTROLLER =================

type FieldController struct {
    fieldRepo *repositories.FieldRepository
}

func NewFieldController(fieldRepo *repositories.FieldRepository) *FieldController {
    return &FieldController{
        fieldRepo: fieldRepo,
    }
}

// SetupRoutes sets up all field-related routes
func (c *FieldController) SetupRoutes(api *gin.RouterGroup) {
    fields := api.Group("/fields")
    {
        fields.GET("", c.GetAvailableFields)
        fields.GET("/by-category/:category", c.GetFieldsByCategory)
    }
}

// GetAvailableFields godoc
// @Summary Get available fields
// @Description Get list of all available fields for reports
// @Tags providerDetail
// @Produce json
// @Success 200 {object} dtos.APIResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/fields [get]
// @Security BearerAuth
func (c *FieldController) GetAvailableFields(ctx *gin.Context) {
    fields, err := c.fieldRepo.GetAllFields()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get available fields",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Available fields retrieved successfully",
        Data:    fields,
    })
}

// GetFieldsByCategory godoc
// @Summary Get fields by category
// @Description Get fields filtered by category (header, data, summary)
// @Tags providerDetail
// @Produce json
// @Param category path string true "Field category" Enums(header, data, summary)
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /provider-detail/providers/fields/by-category/{category} [get]
// @Security BearerAuth
func (c *FieldController) GetFieldsByCategory(ctx *gin.Context) {
    category := ctx.Param("category")
    
    // Validate category
    validCategories := map[string]bool{
        "header":  true,
        "data":    true,
        "summary": true,
    }
    
    if !validCategories[category] {
        ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "Invalid category",
            Details: "Category must be one of: header, data, summary",
        })
        return
    }

    fields, err := c.fieldRepo.GetFieldsByCategory(category)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get fields by category",
            Details: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, dtos.APIResponse{
        Success: true,
        Message: "Fields retrieved successfully",
        Data:    fields,
    })
}