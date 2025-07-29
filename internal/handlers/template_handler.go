package handlers

import (
    "net/http"
    "strconv"
    "provider-report-api/internal/models"
    "provider-report-api/internal/services"
    "github.com/gin-gonic/gin"
)

type TemplateHandler struct {
    service *services.TemplateService
}

func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
    return &TemplateHandler{service: service}
}

func (h *TemplateHandler) GetAvailableFields(c *gin.Context) {
    fields, err := h.service.GetAvailableFields()
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to get available fields",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Available fields retrieved successfully",
        Data:    fields,
    })
}

func (h *TemplateHandler) GetTemplates(c *gin.Context) {
    reportType := c.Query("report_type")
    
    templates, err := h.service.GetTemplates(reportType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to get templates",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Templates retrieved successfully",
        Data:    templates,
    })
}

func (h *TemplateHandler) GetTemplate(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid template ID",
            Error:   err.Error(),
        })
        return
    }
    
    template, err := h.service.GetTemplate(id)
    if err != nil {
        if err.Error() == "template not found" {
            c.JSON(http.StatusNotFound, models.APIResponse{
                Success: false,
                Message: "Template not found",
                Error:   err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to get template",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Template retrieved successfully",
        Data:    template,
    })
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
    var req models.TemplateRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }
    
    // Validation
    if req.Name == "" {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Template name is required",
        })
        return
    }
    
    if req.ReportType == "" {
        req.ReportType = "provider_detail"
    }
    
    if req.CreatedBy == "" {
        req.CreatedBy = "api_user"
    }
    
    template, err := h.service.CreateTemplate(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to create template",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, models.APIResponse{
        Success: true,
        Message: "Template created successfully",
        Data:    template,
    })
}

func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid template ID",
            Error:   err.Error(),
        })
        return
    }
    
    var req models.TemplateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }
    
    template, err := h.service.UpdateTemplate(id, req)
    if err != nil {
        if err.Error() == "template not found" {
            c.JSON(http.StatusNotFound, models.APIResponse{
                Success: false,
                Message: "Template not found",
                Error:   err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to update template",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Template updated successfully",
        Data:    template,
    })
}

func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid template ID",
            Error:   err.Error(),
        })
        return
    }
    
    err = h.service.DeleteTemplate(id)
    if err != nil {
        if err.Error() == "template not found" {
            c.JSON(http.StatusNotFound, models.APIResponse{
                Success: false,
                Message: "Template not found",
                Error:   err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to delete template",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Template deleted successfully",
    })
}