package utils

import (
    "bytes"
    "fmt"
    "time"
    "provider-report-api/internal/models"
    "github.com/jung-kurt/gofpdf"
)

func GeneratePDF(providers []models.Provider, headerFields []string, detailFields []string) ([]byte, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    
    // Set font
    pdf.SetFont("Arial", "B", 16)
    
    // Title
    pdf.Cell(190, 10, "Provider Details Report")
    pdf.Ln(15)
    
    // Header info
    pdf.SetFont("Arial", "", 10)
    pdf.Cell(190, 5, fmt.Sprintf("Generated Date: %s", time.Now().Format("2006-01-02 15:04:05")))
    pdf.Ln(5)
    pdf.Cell(190, 5, fmt.Sprintf("Total Records: %d", len(providers)))
    pdf.Ln(10)
    
    // Table header
    pdf.SetFont("Arial", "B", 8)
    pdf.SetFillColor(200, 200, 200)
    
    // Define column headers based on selected detail fields
    headers := getColumnHeaders(detailFields)
    colWidths := getColumnWidths(len(headers))
    
    for i, header := range headers {
        pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
    }
    pdf.Ln(8)
    
    // Table data
    pdf.SetFont("Arial", "", 7)
    pdf.SetFillColor(255, 255, 255)
    
    for _, provider := range providers {
        values := getRowValues(provider, detailFields)
        
        for i, value := range values {
            if len(value) > 25 {
                value = value[:22] + "..."
            }
            pdf.CellFormat(colWidths[i], 6, value, "1", 0, "L", false, 0, "")
        }
        pdf.Ln(6)
        
        // Add new page if needed
        if pdf.GetY() > 270 {
            pdf.AddPage()
            
            // Re-add headers
            pdf.SetFont("Arial", "B", 8)
            pdf.SetFillColor(200, 200, 200)
            for i, header := range headers {
                pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
            }
            pdf.Ln(8)
            pdf.SetFont("Arial", "", 7)
        }
    }
    
    var buf bytes.Buffer
    err := pdf.Output(&buf)
    if err != nil {
        return nil, fmt.Errorf("failed to generate PDF: %w", err)
    }
    
    return buf.Bytes(), nil
}

func getColumnHeaders(detailFields []string) []string {
    fieldMap := map[string]string{
        "D1":  "ID",
        "D2":  "Type",
        "D4":  "Name (TH)",
        "D5":  "Name (EN)",
        "D6":  "Phone",
        "D8":  "District",
        "D9":  "Province",
        "D11": "Code",
        "D13": "Email",
        "D14": "TPA Network",
        "D15": "Status",
        "D17": "Country",
    }
    
    var headers []string
    if len(detailFields) == 0 {
        // Default fields
        headers = []string{"ID", "Name (TH)", "Type", "Province", "Phone", "Status"}
    } else {
        for _, field := range detailFields {
            if header, exists := fieldMap[field]; exists {
                headers = append(headers, header)
            }
        }
    }
    
    return headers
}

func getColumnWidths(numCols int) []float64 {
    availableWidth := 190.0 // A4 width minus margins
    
    if numCols <= 3 {
        return []float64{30, 80, 80}[:numCols]
    } else if numCols <= 6 {
        colWidth := availableWidth / float64(numCols)
        widths := make([]float64, numCols)
        for i := range widths {
            widths[i] = colWidth
        }
        return widths
    } else {
        // For many columns, use smaller width
        colWidth := availableWidth / float64(numCols)
        widths := make([]float64, numCols)
        for i := range widths {
            widths[i] = colWidth
        }
        return widths
    }
}

func getRowValues(provider models.Provider, detailFields []string) []string {
    fieldMap := map[string]string{
        "D1":  fmt.Sprintf("%d", provider.ID),
        "D2":  provider.ProviderType,
        "D4":  provider.ProviderNameTH,
        "D5":  provider.ProviderNameEN,
        "D6":  provider.TelephoneNumber,
        "D8":  provider.District,
        "D9":  provider.Province,
        "D11": provider.ProviderCode,
        "D13": provider.Email,
        "D14": func() string {
            if provider.IsTPANetwork {
                return "Yes"
            }
            return "No"
        }(),
        "D15": provider.ProviderStatus,
        "D17": provider.Country,
    }
    
    var values []string
    if len(detailFields) == 0 {
        // Default values
        values = []string{
            fmt.Sprintf("%d", provider.ID),
            provider.ProviderNameTH,
            provider.ProviderType,
            provider.Province,
            provider.TelephoneNumber,
            provider.ProviderStatus,
        }
    } else {
        for _, field := range detailFields {
            if value, exists := fieldMap[field]; exists {
                values = append(values, value)
            } else {
                values = append(values, "")
            }
        }
    }
    
    return values
}