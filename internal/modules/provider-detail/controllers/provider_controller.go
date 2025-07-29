package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "provider-report-api/internal/modules/provider-detail/dtos"
    "provider-report-api/internal/modules/provider-detail/services"
)

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

// GetProviderSummary godoc
// @Summary Get provider summary
// @Description Get provider summary statistics
// @Tags providers
// @Accept json
// @Produce json
// @Param provider_name query string false "Provider name"
// @Param province_name query string false "Province name"
// @Param provider_type query string false "Provider type"
// @Param business_type query string false "Business type"
// @Success 200 {object} dtos.APIResponse{data=dtos.ProviderSummaryDTO}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/summary [get]
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

// GenerateReport godoc
// @Summary Generate provider report
// @Description Generate provider report with specified criteria
// @Tags providers
// @Accept json
// @Produce json
// @Param request body dtos.ProviderReportRequestDTO true "Report request"
// @Success 200 {object} dtos.APIResponse{data=dtos.ProviderReportDataDTO}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/report [post]
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
// @Description Export provider report in specified format
// @Tags providers
// @Accept json
// @Produce application/octet-stream
// @Param request body dtos.ProviderReportRequestDTO true "Export request"
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/export [post]
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
// @Description Get list of all provinces
// @Tags providers
// @Accept json
// @Produce json
// @Success 200 {object} dtos.APIResponse{data=[]string}
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/provinces [get]
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
// @Description Get list of all provider types
// @Tags providers
// @Accept json
// @Produce json
// @Success 200 {object} dtos.APIResponse{data=[]string}
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/types [get]
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

// GetProviderStats godoc
// @Summary Get provider statistics
// @Description Get provider statistics
// @Tags providers
// @Accept json
// @Produce json
// @Success 200 {object} dtos.APIResponse{data=dtos.ProviderStatsDTO}
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/stats [get]
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

// CreateProvider godoc
// @Summary Create a new provider
// @Description Create a new provider
// @Tags providers
// @Accept json
// @Produce json
// @Param request body dtos.CreateProviderRequestDTO true "Create provider request"
// @Success 201 {object} dtos.APIResponse{data=dtos.ProviderDTO}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers [post]
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
// @Description Get provider by ID
// @Tags providers
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} dtos.APIResponse{data=dtos.ProviderDTO}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/{id} [get]
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
// @Description Update provider by ID
// @Tags providers
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Param request body dtos.UpdateProviderRequestDTO true "Update provider request"
// @Success 200 {object} dtos.APIResponse{data=dtos.ProviderDTO}
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/{id} [put]
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
// @Tags providers
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Success 200 {object} dtos.APIResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /providers/{id} [delete]
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