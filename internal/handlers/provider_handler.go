package handlers

import (
    "net/http"
    "strconv"
    "provider-report-api/internal/models"
    "provider-report-api/internal/services"
    "provider-report-api/internal/utils"
    "github.com/gin-gonic/gin"
)

type ProviderHandler struct {
    service *services.ProviderService
}

func NewProviderHandler(service *services.ProviderService) *ProviderHandler {
    return &ProviderHandler{service: service}
}

func (h *ProviderHandler) SearchProviders(c *gin.Context) {
    var req models.ProviderSearchRequest
    
    // Bind query parameters
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }
    
    // Search providers
    result, err := h.service.SearchProviders(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to search providers",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Providers retrieved successfully",
        Data:    result,
    })
}

func (h *ProviderHandler) GetProviderByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid provider ID",
            Error:   err.Error(),
        })
        return
    }
    
    provider, err := h.service.GetProviderByID(id)
    if err != nil {
        if err.Error() == "provider not found" {
            c.JSON(http.StatusNotFound, models.APIResponse{
                Success: false,
                Message: "Provider not found",
                Error:   err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to get provider",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Provider retrieved successfully",
        Data:    provider,
    })
}

func (h *ProviderHandler) ExportProviders(c *gin.Context) {
    var req models.ExportRequest
    
    // Bind query parameters
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }
    
    // Default format
    if req.Format == "" {
        req.Format = "pdf"
    }
    
    // Get providers data
    providers, err := h.service.GetProvidersForExport(req.ProviderSearchRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to get providers for export",
            Error:   err.Error(),
        })
        return
    }
    
    if len(providers) == 0 {
        c.JSON(http.StatusNotFound, models.APIResponse{
            Success: false,
            Message: "No providers found for export",
        })
        return
    }
    
    // Create export based on format
    switch req.Format {
    case "pdf":
        pdfBytes, err := utils.GeneratePDF(providers, req.HeaderFields, req.DetailFields)
        if err != nil {
            c.JSON(http.StatusInternalServerError, models.APIResponse{
                Success: false,
                Message: "Failed to generate PDF",
                Error:   err.Error(),
            })
            return
        }
        
        c.Header("Content-Type", "application/pdf")
        c.Header("Content-Disposition", "attachment; filename=providers_report.pdf")
        c.Data(http.StatusOK, "application/pdf", pdfBytes)
        
    case "excel":
        excelBytes, err := utils.GenerateExcel(providers, req.HeaderFields, req.DetailFields)
        if err != nil {
            c.JSON(http.StatusInternalServerError, models.APIResponse{
                Success: false,
                Message: "Failed to generate Excel",
                Error:   err.Error(),
            })
            return
        }
        
        c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
        c.Header("Content-Disposition", "attachment; filename=providers_report.xlsx")
        c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
        
    default:
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Unsupported export format. Use 'pdf' or 'excel'",
        })
    }
}