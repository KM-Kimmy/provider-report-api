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

// SearchProviders godoc
// @Summary Search providers
// @Description Search providers with filters
// @Tags providers
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
// @Router /providers/search [get]
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

// ================= TEMPLATE CONTROLLER =================

type TemplateController struct {
    templateService *services.TemplateService
}

func NewTemplateController(templateService *services.TemplateService) *TemplateController {
    return &TemplateController{
        templateService: templateService,
    }
}

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