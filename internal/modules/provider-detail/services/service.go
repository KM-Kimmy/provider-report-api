package services

import (
    "bytes"
    "fmt"
    "net/smtp"
    "strconv"
    "time"

    "github.com/xuri/excelize/v2"
    "provider-report-api/configs"
    "provider-report-api/internal/modules/provider-detail/dtos"
    "provider-report-api/internal/modules/provider-detail/repositories"
)

// ProviderService handles provider business logic
type ProviderService struct {
    providerRepo *repositories.ProviderRepository
    exportService *ExportService
    fieldRepo    *repositories.FieldRepository
}

func NewProviderService(providerRepo *repositories.ProviderRepository, exportService *ExportService, fieldRepo *repositories.FieldRepository) *ProviderService {
    return &ProviderService{
        providerRepo:  providerRepo,
        exportService: exportService,
        fieldRepo:     fieldRepo,
    }
}

func (s *ProviderService) SearchProviders(req dtos.ProviderSearchRequestDTO) ([]dtos.ProviderDTO, int64, error) {
    return s.providerRepo.Search(req)
}

func (s *ProviderService) GetProviderSummary(req dtos.ProviderSearchRequestDTO) (*dtos.ProviderSummaryDTO, error) {
    return s.providerRepo.GetSummary(req)
}

func (s *ProviderService) GenerateReport(req dtos.ProviderReportRequestDTO) (*dtos.ProviderReportDataDTO, error) {
    // Get provider data
    providers, total, err := s.providerRepo.Search(req.SearchParams)
    if err != nil {
        return nil, fmt.Errorf("failed to search providers: %w", err)
    }

    // Get summary
    summary, err := s.providerRepo.GetSummary(req.SearchParams)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider summary: %w", err)
    }

    // Generate header information
    header := map[string]interface{}{
        "generated_at":    time.Now(),
        "criteria":        req.SearchParams,
        "total_records":   total,
        "template_id":     req.TemplateID,
        "format_type":     req.FormatType,
    }

    return &dtos.ProviderReportDataDTO{
        Header:    header,
        Summary:   *summary,
        Providers: providers,
        Total:     total,
    }, nil
}

func (s *ProviderService) ExportReport(req dtos.ProviderReportRequestDTO) ([]byte, string, string, error) {
    // Generate report data
    reportData, err := s.GenerateReport(req)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to generate report: %w", err)
    }

    // Get fields for export
    var fields []dtos.AvailableFieldDTO
    if len(req.CustomFields) > 0 {
        fields, err = s.fieldRepo.GetFieldsForExport(req.CustomFields)
        if err != nil {
            return nil, "", "", fmt.Errorf("failed to get custom fields: %w", err)
        }
    } else {
        // Use all available fields
        fields, err = s.fieldRepo.GetAllFields()
        if err != nil {
            return nil, "", "", fmt.Errorf("failed to get all fields: %w", err)
        }
    }

    // Export based on format
    switch req.FormatType {
    case "excel":
        return s.exportService.ExportToExcel(reportData, fields)
    case "pdf":
        return s.exportService.ExportToPDF(reportData, fields)
    case "word":
        return s.exportService.ExportToWord(reportData, fields)
    default:
        return s.exportService.ExportToExcel(reportData, fields)
    }
}

func (s *ProviderService) GetProvinces() ([]string, error) {
    return s.providerRepo.GetProvinces()
}

func (s *ProviderService) GetProviderTypes() ([]string, error) {
    return s.providerRepo.GetProviderTypes()
}

func (s *ProviderService) GetProviderStats() (map[string]interface{}, error) {
    return s.providerRepo.GetProviderStats()
}

func (s *ProviderService) CreateProvider(req dtos.CreateProviderRequestDTO) (*dtos.ProviderDTO, error) {
    provider := &dtos.ProviderDTO{
        ProviderCode:   req.ProviderCode,
        NameThai:       req.NameThai,
        NameEng:        req.NameEng,
        ProviderType:   req.ProviderType,
        BusinessType:   req.BusinessType,
        Province:       req.Province,
        GeneralPhoneNo: req.GeneralPhoneNo,
        IsTPANetwork:   req.IsTPANetwork,
        ProviderStatus: req.ProviderStatus,
    }

    err := s.providerRepo.Create(provider)
    if err != nil {
        return nil, fmt.Errorf("failed to create provider: %w", err)
    }

    return provider, nil
}

func (s *ProviderService) GetProviderByID(id int) (*dtos.ProviderDTO, error) {
    return s.providerRepo.GetByID(id)
}

func (s *ProviderService) UpdateProvider(id int, req dtos.UpdateProviderRequestDTO) (*dtos.ProviderDTO, error) {
    provider, err := s.providerRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("provider not found: %w", err)
    }

    // Update fields
    provider.NameThai = req.NameThai
    provider.NameEng = req.NameEng
    provider.ProviderType = req.ProviderType
    provider.BusinessType = req.BusinessType
    provider.Province = req.Province
    provider.GeneralPhoneNo = req.GeneralPhoneNo
    provider.IsTPANetwork = req.IsTPANetwork
    provider.ProviderStatus = req.ProviderStatus

    err = s.providerRepo.Update(provider)
    if err != nil {
        return nil, fmt.Errorf("failed to update provider: %w", err)
    }

    return provider, nil
}

func (s *ProviderService) DeleteProvider(id int) error {
    return s.providerRepo.Delete(id)
}

// TemplateService handles template business logic
type TemplateService struct {
    templateRepo *repositories.TemplateRepository
    fieldRepo    *repositories.FieldRepository
}

func NewTemplateService(templateRepo *repositories.TemplateRepository, fieldRepo *repositories.FieldRepository) *TemplateService {
    return &TemplateService{
        templateRepo: templateRepo,
        fieldRepo:    fieldRepo,
    }
}

func (s *TemplateService) GetAllTemplates() ([]dtos.TemplateDTO, error) {
    return s.templateRepo.GetAll()
}

func (s *TemplateService) GetTemplate(id int) (*dtos.TemplateDTO, error) {
    return s.templateRepo.GetByID(id)
}

func (s *TemplateService) CreateTemplate(req dtos.CreateTemplateRequestDTO) (*dtos.TemplateDTO, error) {
    // Validate fields
    allFields := append(req.HeaderFields, req.DataFields...)
    allFields = append(allFields, req.SummaryFields...)
    
    validationResults, err := s.fieldRepo.ValidateFields(allFields)
    if err != nil {
        return nil, fmt.Errorf("failed to validate fields: %w", err)
    }

    for _, result := range validationResults {
        if !result.IsValid {
            return nil, fmt.Errorf("invalid field code: %s - %s", result.FieldCode, result.Message)
        }
    }

    template := &dtos.TemplateDTO{
        TemplateName:   req.TemplateName,
        IsStandard:     req.IsStandard,
        Description:    &req.Description,
        HeaderFields:   dtos.JSONFieldArray(req.HeaderFields),
        DataFields:     dtos.JSONFieldArray(req.DataFields),
        SummaryFields:  dtos.JSONFieldArray(req.SummaryFields),
        FieldPositions: &req.FieldPositions,
        CreatedBy:      "system", // TODO: Get from context
    }

    err = s.templateRepo.Create(template)
    if err != nil {
        return nil, fmt.Errorf("failed to create template: %w", err)
    }

    return template, nil
}

func (s *TemplateService) UpdateTemplate(id int, req dtos.UpdateTemplateRequestDTO) (*dtos.TemplateDTO, error) {
    template, err := s.templateRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("template not found: %w", err)
    }

    // Validate fields
    allFields := append(req.HeaderFields, req.DataFields...)
    allFields = append(allFields, req.SummaryFields...)
    
    validationResults, err := s.fieldRepo.ValidateFields(allFields)
    if err != nil {
        return nil, fmt.Errorf("failed to validate fields: %w", err)
    }

    for _, result := range validationResults {
        if !result.IsValid {
            return nil, fmt.Errorf("invalid field code: %s - %s", result.FieldCode, result.Message)
        }
    }

    // Update fields
    template.TemplateName = req.TemplateName
    template.IsStandard = req.IsStandard
    template.Description = &req.Description
    template.HeaderFields = dtos.JSONFieldArray(req.HeaderFields)
    template.DataFields = dtos.JSONFieldArray(req.DataFields)
    template.SummaryFields = dtos.JSONFieldArray(req.SummaryFields)
    template.FieldPositions = &req.FieldPositions
    template.UpdatedBy = stringPtr("system") // TODO: Get from context

    err = s.templateRepo.Update(template)
    if err != nil {
        return nil, fmt.Errorf("failed to update template: %w", err)
    }

    return template, nil
}

func (s *TemplateService) DeleteTemplate(id int) error {
    return s.templateRepo.Delete(id)
}

// ScheduleService handles schedule business logic
type ScheduleService struct {
    scheduleRepo  *repositories.ScheduleRepository
    templateRepo  *repositories.TemplateRepository
    emailService  *EmailService
}

func NewScheduleService(scheduleRepo *repositories.ScheduleRepository, templateRepo *repositories.TemplateRepository, emailService *EmailService) *ScheduleService {
    return &ScheduleService{
        scheduleRepo: scheduleRepo,
        templateRepo: templateRepo,
        emailService: emailService,
    }
}

func (s *ScheduleService) GetAllSchedules() ([]dtos.ScheduleDTO, error) {
    return s.scheduleRepo.GetAll()
}

func (s *ScheduleService) GetSchedule(id int) (*dtos.ScheduleDTO, error) {
    return s.scheduleRepo.GetByID(id)
}

func (s *ScheduleService) CreateSchedule(req dtos.CreateScheduleRequestDTO) (*dtos.ScheduleDTO, error) {
    // Validate template exists
    _, err := s.templateRepo.GetByID(req.TemplateID)
    if err != nil {
        return nil, fmt.Errorf("template not found: %w", err)
    }

    // Convert search criteria to JSON
    searchCriteria := dtos.JSONMap(req.SearchCriteria)

    schedule := &dtos.ScheduleDTO{
        ScheduleName:   req.ScheduleName,
        TemplateID:     req.TemplateID,
        EmailTo:        req.EmailTo,
        EmailCC:        &req.EmailCC,
        EmailBCC:       &req.EmailBCC,
        Frequency:      req.Frequency,
        ScheduleDays:   dtos.JSONFieldArray(req.ScheduleDays),
        StartDate:      req.StartDate,
        EndDate:        req.EndDate,
        StartTime:      req.StartTime,
        Timezone:       req.Timezone,
        IsActive:       true,
        SearchCriteria: searchCriteria,
        ExportFormat:   req.ExportFormat,
        CreatedBy:      "system", // TODO: Get from context
    }

    err = s.scheduleRepo.Create(schedule)
    if err != nil {
        return nil, fmt.Errorf("failed to create schedule: %w", err)
    }

    return schedule, nil
}

func (s *ScheduleService) UpdateSchedule(id int, req dtos.UpdateScheduleRequestDTO) (*dtos.ScheduleDTO, error) {
    schedule, err := s.scheduleRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("schedule not found: %w", err)
    }

    // Validate template exists
    _, err = s.templateRepo.GetByID(req.TemplateID)
    if err != nil {
        return nil, fmt.Errorf("template not found: %w", err)
    }

    // Convert search criteria to JSON
    searchCriteria := dtos.JSONMap(req.SearchCriteria)

    // Update fields
    schedule.ScheduleName = req.ScheduleName
    schedule.TemplateID = req.TemplateID
    schedule.EmailTo = req.EmailTo
    schedule.EmailCC = &req.EmailCC
    schedule.EmailBCC = &req.EmailBCC
    schedule.Frequency = req.Frequency
    schedule.ScheduleDays = dtos.JSONFieldArray(req.ScheduleDays)
    schedule.StartDate = req.StartDate
    schedule.EndDate = req.EndDate
    schedule.StartTime = req.StartTime
    schedule.Timezone = req.Timezone
    schedule.IsActive = req.IsActive
    schedule.SearchCriteria = searchCriteria
    schedule.ExportFormat = req.ExportFormat
    schedule.UpdatedBy = stringPtr("system") // TODO: Get from context

    err = s.scheduleRepo.Update(schedule)
    if err != nil {
        return nil, fmt.Errorf("failed to update schedule: %w", err)
    }

    return schedule, nil
}

func (s *ScheduleService) DeleteSchedule(id int) error {
    return s.scheduleRepo.Delete(id)
}

func (s *ScheduleService) RunSchedule(id int) (*dtos.RunScheduleResponseDTO, error) {
    schedule, err := s.scheduleRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("schedule not found: %w", err)
    }

    // TODO: Implement actual schedule execution
    // This would involve:
    // 1. Converting SearchCriteria to ProviderSearchRequestDTO
    // 2. Generating report
    // 3. Sending email
    // 4. Logging the result

    // Update last run time
    err = s.scheduleRepo.UpdateLastRun(id)
    if err != nil {
        return nil, fmt.Errorf("failed to update last run: %w", err)
    }

    return &dtos.RunScheduleResponseDTO{
        Message:     "Schedule executed successfully",
        ExecutedAt:  time.Now(),
        Recipients:  schedule.EmailTo,
        RecordCount: 0, // TODO: Get actual count
        FileSize:    "0 KB", // TODO: Get actual size
        Status:      "success",
    }, nil
}

// LogService handles log business logic
type LogService struct {
    logRepo *repositories.LogRepository
}

func NewLogService(logRepo *repositories.LogRepository) *LogService {
    return &LogService{
        logRepo: logRepo,
    }
}

func (s *LogService) GetSentReportLogs(req dtos.LogSearchRequestDTO) (*dtos.LogListResponseDTO, error) {
    logs, total, err := s.logRepo.GetSentReportLogs(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get sent report logs: %w", err)
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

    return &dtos.LogListResponseDTO{
        Logs:       logs,
        Total:      total,
        Page:       req.Page,
        Limit:      req.Limit,
        TotalPages: totalPages,
    }, nil
}

func (s *LogService) GetSentReportLog(id int) (*dtos.SentReportLogDTO, error) {
    return s.logRepo.GetByID(id)
}

func (s *LogService) CreateLog(req dtos.CreateLogRequestDTO) (*dtos.SentReportLogDTO, error) {
    log := &dtos.SentReportLogDTO{
        TemplateID:      req.TemplateID,
        ScheduleID:      req.ScheduleID,
        Recipients:      req.Recipients,
        Subject:         req.Subject,
        FileName:        req.FileName,
        FileSizeKB:      req.FileSizeKB,
        ExportFormat:    req.ExportFormat,
        TotalRecords:    req.TotalRecords,
        Status:          req.Status,
        ErrorMessage:    req.ErrorMessage,
        RetryCount:      req.RetryCount,
        ExecutionTimeMs: req.ExecutionTimeMs,
    }

    err := s.logRepo.Create(log)
    if err != nil {
        return nil, fmt.Errorf("failed to create log: %w", err)
    }

    return log, nil
}

// ExportService handles file export operations
type ExportService struct{}

func NewExportService() *ExportService {
    return &ExportService{}
}

func (s *ExportService) ExportToExcel(data *dtos.ProviderReportDataDTO, fields []dtos.AvailableFieldDTO) ([]byte, string, string, error) {
    f := excelize.NewFile()
    defer f.Close()

    // Create header row
    headerRow := 1
    for i, field := range fields {
        colName, _ := excelize.ColumnNumberToName(i + 1)
        f.SetCellValue("Sheet1", colName+strconv.Itoa(headerRow), field.FieldNameEng)
    }

    // Add data rows
    for rowIdx, provider := range data.Providers {
        dataRow := rowIdx + 2
        for i, field := range fields {
            colName, _ := excelize.ColumnNumberToName(i + 1)
            value := s.getProviderFieldValue(provider, field.FieldCode)
            f.SetCellValue("Sheet1", colName+strconv.Itoa(dataRow), value)
        }
    }

    // Save to buffer
    var buf bytes.Buffer
    err := f.Write(&buf)
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to write Excel file: %w", err)
    }

    filename := fmt.Sprintf("provider_report_%s.xlsx", time.Now().Format("20060102_150405"))
    contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

    return buf.Bytes(), filename, contentType, nil
}

func (s *ExportService) ExportToPDF(data *dtos.ProviderReportDataDTO, fields []dtos.AvailableFieldDTO) ([]byte, string, string, error) {
    // TODO: Implement PDF export
    return nil, "", "", fmt.Errorf("PDF export not implemented yet")
}

func (s *ExportService) ExportToWord(data *dtos.ProviderReportDataDTO, fields []dtos.AvailableFieldDTO) ([]byte, string, string, error) {
    // TODO: Implement Word export
    return nil, "", "", fmt.Errorf("Word export not implemented yet")
}

func (s *ExportService) getProviderFieldValue(provider dtos.ProviderDTO, fieldCode string) interface{} {
    switch fieldCode {
    case "provider_code":
        return provider.ProviderCode
    case "name_thai":
        return provider.NameThai
    case "name_eng":
        if provider.NameEng != nil {
            return *provider.NameEng
        }
        return ""
    case "provider_type":
        return provider.ProviderType
    case "business_type":
        if provider.BusinessType != nil {
            return *provider.BusinessType
        }
        return ""
    case "province":
        return provider.Province
    case "general_phone_no":
        if provider.GeneralPhoneNo != nil {
            return *provider.GeneralPhoneNo
        }
        return ""
    case "is_tpa_network":
        return provider.IsTPANetwork
    case "provider_status":
        return provider.ProviderStatus
    case "created_at":
        return provider.CreatedAt.Format("2006-01-02 15:04:05")
    default:
        return ""
    }
}

// EmailService handles email operations
type EmailService struct {
    config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
    return &EmailService{
        config: cfg,
    }
}

func (s *EmailService) SendEmail(to, subject, body string, attachment []byte, filename string) error {
    // TODO: Implement actual email sending
    // This is a placeholder implementation
    
    auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
    
    msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
    
    err := smtp.SendMail(
        s.config.SMTPHost+":"+s.config.SMTPPort,
        auth,
        s.config.SMTPFrom,
        []string{to},
        []byte(msg),
    )
    
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    
    return nil
}

func (s *EmailService) SendScheduledReport(schedule dtos.ScheduleDTO, reportData []byte, filename string) error {
    subject := fmt.Sprintf("Scheduled Report: %s", schedule.ScheduleName)
    body := fmt.Sprintf("This is an automated report generated at %s", time.Now().Format("2006-01-02 15:04:05"))
    
    return s.SendEmail(schedule.EmailTo, subject, body, reportData, filename)
}

// Helper functions
func stringPtr(s string) *string {
    return &s
}