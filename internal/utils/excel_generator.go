package utils

import (
    "bytes"
    "fmt"
    "strconv"
    "time"
    "provider-report-api/internal/models"
    "github.com/xuri/excelize/v2"
)

func GenerateExcel(providers []models.Provider, headerFields []string, detailFields []string) ([]byte, error) {
    f := excelize.NewFile()
    
    sheetName := "Provider Report"
    index, err := f.NewSheet(sheetName)
    if err != nil {
        return nil, fmt.Errorf("failed to create sheet: %w", err)
    }
    
    f.SetActiveSheet(index)
    
    // Style for header
    headerStyle, err := f.NewStyle(&excelize.Style{
        Font: &excelize.Font{
            Bold: true,
            Size: 12,
        },
        Alignment: &excelize.Alignment{
            Horizontal: "center",
            Vertical:   "center",
        },
        Fill: excelize.Fill{
            Type:    "pattern",
            Color:   []string{"#E6E6FA"},
            Pattern: 1,
        },
        Border: []excelize.Border{
            {Type: "left", Color: "000000", Style: 1},
            {Type: "top", Color: "000000", Style: 1},
            {Type: "bottom", Color: "000000", Style: 1},
            {Type: "right", Color: "000000", Style: 1},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create header style: %w", err)
    }
    
    // Style for data
    dataStyle, err := f.NewStyle(&excelize.Style{
        Border: []excelize.Border{
            {Type: "left", Color: "000000", Style: 1},
            {Type: "top", Color: "000000", Style: 1},
            {Type: "bottom", Color: "000000", Style: 1},
            {Type: "right", Color: "000000", Style: 1},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create data style: %w", err)
    }
    
    row := 1
    
    // Report title
    f.SetCellValue(sheetName, "A1", "Provider Details Report")
    titleStyle, _ := f.NewStyle(&excelize.Style{
        Font: &excelize.Font{Bold: true, Size: 16},
    })
    f.SetCellStyle(sheetName, "A1", "A1", titleStyle)
    row += 2
    
    // Report metadata
    f.SetCellValue(sheetName, "A"+strconv.Itoa(row), "Generated Date:")
    f.SetCellValue(sheetName, "B"+strconv.Itoa(row), time.Now().Format("2006-01-02 15:04:05"))
    row++
    
    f.SetCellValue(sheetName, "A"+strconv.Itoa(row), "Total Records:")
    f.SetCellValue(sheetName, "B"+strconv.Itoa(row), len(providers))
    row += 2
    
    // Column headers
    headers := getExcelColumnHeaders(detailFields)
    for i, header := range headers {
        col := string(rune('A' + i))
        cell := col + strconv.Itoa(row)
        f.SetCellValue(sheetName, cell, header)
        f.SetCellStyle(sheetName, cell, cell, headerStyle)
    }
    row++
    
    // Data rows
    for _, provider := range providers {
        values := getExcelRowValues(provider, detailFields)
        
        for i, value := range values {
            col := string(rune('A' + i))
            cell := col + strconv.Itoa(row)
            f.SetCellValue(sheetName, cell, value)
            f.SetCellStyle(sheetName, cell, cell, dataStyle)
        }
        row++
    }
    
    // Auto-fit columns
    for i := 0; i < len(headers); i++ {
        col := string(rune('A' + i))
        f.SetColWidth(sheetName, col, col, 15)
    }
    
    // Delete default sheet
    f.DeleteSheet("Sheet1")
    
    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, fmt.Errorf("failed to write Excel file: %w", err)
    }
    
    return buf.Bytes(), nil
}

func getExcelColumnHeaders(detailFields []string) []string {
    fieldMap := map[string]string{
        "D1":  "ID",
        "D2":  "Provider Type",
        "D3":  "Status",
        "D4":  "Provider Name (TH)",
        "D5":  "Provider Name (EN)",
        "D6":  "Telephone",
        "D7":  "Address",
        "D8":  "District",
        "D9":  "Province",
        "D10": "Post Code",
        "D11": "Provider Code",
        "D12": "Business Type",
        "D13": "Email",
        "D14": "TPA Network",
        "D15": "Provider Status",
        "D16": "Region",
        "D17": "Country",
        "D18": "Created Date",
    }
    
    var headers []string
    if len(detailFields) == 0 {
        // Default fields
        headers = []string{
            "ID", "Provider Name (TH)", "Provider Type", 
            "Province", "Telephone", "Email", "Status",
        }
    } else {
        for _, field := range detailFields {
            if header, exists := fieldMap[field]; exists {
                headers = append(headers, header)
            }
        }
    }
    
    return headers
}

func getExcelRowValues(provider models.Provider, detailFields []string) []interface{} {
    fieldMap := map[string]interface{}{
        "D1":  provider.ID,
        "D2":  provider.ProviderType,
        "D3":  provider.Status,
        "D4":  provider.ProviderNameTH,
        "D5":  provider.ProviderNameEN,
        "D6":  provider.TelephoneNumber,
        "D7":  provider.Address,
        "D8":  provider.District,
        "D9":  provider.Province,
        "D10": provider.PostCode,
        "D11": provider.ProviderCode,
        "D12": provider.BusinessType,
        "D13": provider.Email,
        "D14": func() string {
            if provider.IsTPANetwork {
                return "Yes"
            }
            return "No"
        }(),
        "D15": provider.ProviderStatus,
        "D16": provider.Region,
        "D17": provider.Country,
        "D18": provider.CreatedAt.Format("2006-01-02"),
    }
    
    var values []interface{}
    if len(detailFields) == 0 {
        // Default values
        values = []interface{}{
            provider.ID,
            provider.ProviderNameTH,
            provider.ProviderType,
            provider.Province,
            provider.TelephoneNumber,
            provider.Email,
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