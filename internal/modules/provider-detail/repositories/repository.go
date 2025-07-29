package repositories

import (
    "fmt"
    "strings"

    "github.com/jmoiron/sqlx"
    "provider-report-api/internal/modules/provider-detail/dtos"
)

// ProviderRepository handles provider data operations
type ProviderRepository struct {
    db *sqlx.DB
}

func NewProviderRepository(db *sqlx.DB) *ProviderRepository {
    return &ProviderRepository{db: db}
}

func (r *ProviderRepository) Search(req dtos.ProviderSearchRequestDTO) ([]dtos.ProviderDTO, int64, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1

    baseQuery := `
        SELECT p.*
        FROM providers p
        WHERE 1=1
    `

    countQuery := `
        SELECT COUNT(*)
        FROM providers p
        WHERE 1=1
    `

    // Add search conditions
    if req.ProviderName != "" {
        conditions = append(conditions, fmt.Sprintf("(p.name_thai ILIKE $%d OR p.name_eng ILIKE $%d)", argIndex, argIndex))
        args = append(args, "%"+req.ProviderName+"%")
        argIndex++
    }

    if req.ProvinceName != "" {
        conditions = append(conditions, fmt.Sprintf("p.province ILIKE $%d", argIndex))
        args = append(args, "%"+req.ProvinceName+"%")
        argIndex++
    }

    if req.ProviderType != "" {
        conditions = append(conditions, fmt.Sprintf("p.provider_type = $%d", argIndex))
        args = append(args, req.ProviderType)
        argIndex++
    }

    if req.BusinessType != "" {
        conditions = append(conditions, fmt.Sprintf("p.business_type = $%d", argIndex))
        args = append(args, req.BusinessType)
        argIndex++
    }

    if req.IsTPANetwork != nil {
        conditions = append(conditions, fmt.Sprintf("p.is_tpa_network = $%d", argIndex))
        args = append(args, *req.IsTPANetwork)
        argIndex++
    }

    if req.CreatedFrom != nil {
        conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
        args = append(args, req.CreatedFrom.Format("2006-01-02"))
        argIndex++
    }

    if req.CreatedTo != nil {
        conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
        args = append(args, req.CreatedTo.Format("2006-01-02 23:59:59"))
        argIndex++
    }

    // Add conditions to queries
    if len(conditions) > 0 {
        conditionStr := " AND " + strings.Join(conditions, " AND ")
        baseQuery += conditionStr
        countQuery += conditionStr
    }

    // Get total count
    var total int64
    err := r.db.Get(&total, countQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get provider count: %w", err)
    }

    // Add pagination
    if req.Limit == 0 {
        req.Limit = 10
    }
    offset := (req.Page - 1) * req.Limit
    if offset < 0 {
        offset = 0
    }

    baseQuery += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, req.Limit, offset)

    // Execute query
    var providers []dtos.ProviderDTO
    err = r.db.Select(&providers, baseQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to search providers: %w", err)
    }

    return providers, total, nil
}

func (r *ProviderRepository) GetSummary(req dtos.ProviderSearchRequestDTO) (*dtos.ProviderSummaryDTO, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1

    query := `
        SELECT 
            COUNT(CASE WHEN provider_type = 'Hospital' THEN 1 END) as hospital,
            COUNT(CASE WHEN provider_type = 'Clinic' THEN 1 END) as clinic,
            COUNT(*) as grand_total
        FROM providers p
        WHERE 1=1
    `

    // Add same search conditions as in Search method
    if req.ProviderName != "" {
        conditions = append(conditions, fmt.Sprintf("(p.name_thai ILIKE $%d OR p.name_eng ILIKE $%d)", argIndex, argIndex))
        args = append(args, "%"+req.ProviderName+"%")
        argIndex++
    }

    if req.ProvinceName != "" {
        conditions = append(conditions, fmt.Sprintf("p.province ILIKE $%d", argIndex))
        args = append(args, "%"+req.ProvinceName+"%")
        argIndex++
    }

    if req.ProviderType != "" {
        conditions = append(conditions, fmt.Sprintf("p.provider_type = $%d", argIndex))
        args = append(args, req.ProviderType)
        argIndex++
    }

    if req.BusinessType != "" {
        conditions = append(conditions, fmt.Sprintf("p.business_type = $%d", argIndex))
        args = append(args, req.BusinessType)
        argIndex++
    }

    if req.IsTPANetwork != nil {
        conditions = append(conditions, fmt.Sprintf("p.is_tpa_network = $%d", argIndex))
        args = append(args, *req.IsTPANetwork)
        argIndex++
    }

    if req.CreatedFrom != nil {
        conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
        args = append(args, req.CreatedFrom.Format("2006-01-02"))
        argIndex++
    }

    if req.CreatedTo != nil {
        conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
        args = append(args, req.CreatedTo.Format("2006-01-02 23:59:59"))
        argIndex++
    }

    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    var summary dtos.ProviderSummaryDTO
    err := r.db.Get(&summary, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider summary: %w", err)
    }

    // Set type and province
    summary.Type = "Government"
    if req.ProvinceName != "" {
        summary.Province = req.ProvinceName
    }

    return &summary, nil
}

func (r *ProviderRepository) GetByID(id int) (*dtos.ProviderDTO, error) {
    var provider dtos.ProviderDTO
    err := r.db.Get(&provider, "SELECT * FROM providers WHERE id = $1", id)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider by ID: %w", err)
    }
    return &provider, nil
}

func (r *ProviderRepository) Create(provider *dtos.ProviderDTO) error {
    query := `
        INSERT INTO providers (
            provider_code, title_thai, name_thai, title_eng, name_eng,
            provider_type, register_status, business_type, bed_size,
            eligibility_method, province, region, country, provider_tax_id,
            wh_tax_percent, exempt_percent, wh_tax_exempt_from, wh_tax_exempt_to,
            opening_time, provider_status, building_no, village_no, lane_alley,
            road, sub_district, district, post_code, title_name, department,
            general_phone_no, direct_phone_no, email, email_to_list, email_cc_list,
            payment_method, payment_branch_id, payee_name, bank_account_number,
            bank_account_type, bank_branch_name, bank_name, is_tpa_network,
            has_incident, discount_categories, pricing_categories, created_by
        ) VALUES (
            :provider_code, :title_thai, :name_thai, :title_eng, :name_eng,
            :provider_type, :register_status, :business_type, :bed_size,
            :eligibility_method, :province, :region, :country, :provider_tax_id,
            :wh_tax_percent, :exempt_percent, :wh_tax_exempt_from, :wh_tax_exempt_to,
            :opening_time, :provider_status, :building_no, :village_no, :lane_alley,
            :road, :sub_district, :district, :post_code, :title_name, :department,
            :general_phone_no, :direct_phone_no, :email, :email_to_list, :email_cc_list,
            :payment_method, :payment_branch_id, :payee_name, :bank_account_number,
            :bank_account_type, :bank_branch_name, :bank_name, :is_tpa_network,
            :has_incident, :discount_categories, :pricing_categories, :created_by
        ) RETURNING id, created_at, updated_at
    `

    rows, err := r.db.NamedQuery(query, provider)
    if err != nil {
        return fmt.Errorf("failed to create provider: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)
        if err != nil {
            return fmt.Errorf("failed to scan created provider: %w", err)
        }
    }

    return nil
}

func (r *ProviderRepository) Update(provider *dtos.ProviderDTO) error {
    query := `
        UPDATE providers SET
            name_thai = :name_thai,
            name_eng = :name_eng,
            provider_type = :provider_type,
            business_type = :business_type,
            province = :province,
            general_phone_no = :general_phone_no,
            is_tpa_network = :is_tpa_network,
            provider_status = :provider_status,
            updated_by = :updated_by,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = :id
    `

    result, err := r.db.NamedExec(query, provider)
    if err != nil {
        return fmt.Errorf("failed to update provider: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("provider not found")
    }

    return nil
}

func (r *ProviderRepository) Delete(id int) error {
    query := `DELETE FROM providers WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete provider: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("provider not found")
    }

    return nil
}

func (r *ProviderRepository) GetProvinces() ([]string, error) {
    var provinces []string
    query := `SELECT DISTINCT province FROM providers WHERE province IS NOT NULL ORDER BY province`
    err := r.db.Select(&provinces, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get provinces: %w", err)
    }
    return provinces, nil
}

func (r *ProviderRepository) GetProviderTypes() ([]string, error) {
    var types []string
    query := `SELECT DISTINCT provider_type FROM providers WHERE provider_type IS NOT NULL ORDER BY provider_type`
    err := r.db.Select(&types, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider types: %w", err)
    }
    return types, nil
}

func (r *ProviderRepository) GetProviderStats() (map[string]interface{}, error) {
    query := `
        SELECT 
            COUNT(*) as total_providers,
            COUNT(CASE WHEN provider_type = 'Hospital' THEN 1 END) as total_hospitals,
            COUNT(CASE WHEN provider_type = 'Clinic' THEN 1 END) as total_clinics,
            COUNT(CASE WHEN is_tpa_network = true THEN 1 END) as tpa_network_providers,
            COUNT(CASE WHEN provider_status = 'Active' THEN 1 END) as active_providers,
            COUNT(CASE WHEN provider_status = 'Inactive' THEN 1 END) as inactive_providers
        FROM providers
    `

    row := r.db.QueryRow(query)
    
    var totalProviders, totalHospitals, totalClinics, tpaNetworkProviders, activeProviders, inactiveProviders int
    err := row.Scan(&totalProviders, &totalHospitals, &totalClinics, &tpaNetworkProviders, &activeProviders, &inactiveProviders)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider stats: %w", err)
    }

    stats := map[string]interface{}{
        "total_providers":         totalProviders,
        "total_hospitals":         totalHospitals,
        "total_clinics":          totalClinics,
        "tpa_network_providers":  tpaNetworkProviders,
        "active_providers":       activeProviders,
        "inactive_providers":     inactiveProviders,
    }

    return stats, nil
}

// TemplateRepository handles template data operations
type TemplateRepository struct {
    db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) *TemplateRepository {
    return &TemplateRepository{db: db}
}

func (r *TemplateRepository) GetAll() ([]dtos.TemplateDTO, error) {
    var templates []dtos.TemplateDTO
    query := `
        SELECT * FROM templates 
        WHERE is_deleted = false 
        ORDER BY is_standard DESC, created_at DESC
    `
    err := r.db.Select(&templates, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get templates: %w", err)
    }
    return templates, nil
}

func (r *TemplateRepository) GetByID(id int) (*dtos.TemplateDTO, error) {
    var template dtos.TemplateDTO
    query := `SELECT * FROM templates WHERE id = $1 AND is_deleted = false`
    err := r.db.Get(&template, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get template by ID: %w", err)
    }
    return &template, nil
}

func (r *TemplateRepository) Create(template *dtos.TemplateDTO) error {
    query := `
        INSERT INTO templates (
            template_name, is_standard, description, header_fields, 
            data_fields, summary_fields, field_positions, created_by
        ) VALUES (
            :template_name, :is_standard, :description, :header_fields,
            :data_fields, :summary_fields, :field_positions, :created_by
        ) RETURNING id, created_at, updated_at
    `

    rows, err := r.db.NamedQuery(query, template)
    if err != nil {
        return fmt.Errorf("failed to create template: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
        if err != nil {
            return fmt.Errorf("failed to scan created template: %w", err)
        }
    }

    return nil
}

func (r *TemplateRepository) Update(template *dtos.TemplateDTO) error {
    query := `
        UPDATE templates SET
            template_name = :template_name,
            is_standard = :is_standard,
            description = :description,
            header_fields = :header_fields,
            data_fields = :data_fields,
            summary_fields = :summary_fields,
            field_positions = :field_positions,
            updated_by = :updated_by,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = :id AND is_deleted = false
    `

    result, err := r.db.NamedExec(query, template)
    if err != nil {
        return fmt.Errorf("failed to update template: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("template not found")
    }

    return nil
}

func (r *TemplateRepository) Delete(id int) error {
    query := `UPDATE templates SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete template: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("template not found")
    }

    return nil
}

// ScheduleRepository handles schedule data operations
type ScheduleRepository struct {
    db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
    return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetAll() ([]dtos.ScheduleDTO, error) {
    var schedules []dtos.ScheduleDTO
    query := `
        SELECT s.*, t.template_name
        FROM schedules s
        LEFT JOIN templates t ON s.template_id = t.id
        WHERE s.is_deleted = false
        ORDER BY s.created_at DESC
    `
    err := r.db.Select(&schedules, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get schedules: %w", err)
    }
    return schedules, nil
}

func (r *ScheduleRepository) GetByID(id int) (*dtos.ScheduleDTO, error) {
    var schedule dtos.ScheduleDTO
    query := `
        SELECT s.*, t.template_name
        FROM schedules s
        LEFT JOIN templates t ON s.template_id = t.id
        WHERE s.id = $1 AND s.is_deleted = false
    `
    err := r.db.Get(&schedule, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
    }
    return &schedule, nil
}

func (r *ScheduleRepository) Create(schedule *dtos.ScheduleDTO) error {
    query := `
        INSERT INTO schedules (
            schedule_name, template_id, email_to, email_cc, email_bcc,
            frequency, schedule_days, start_date, end_date, start_time,
            timezone, search_criteria, export_format, created_by
        ) VALUES (
            :schedule_name, :template_id, :email_to, :email_cc, :email_bcc,
            :frequency, :schedule_days, :start_date, :end_date, :start_time,
            :timezone, :search_criteria, :export_format, :created_by
        ) RETURNING id, created_at, updated_at
    `

    rows, err := r.db.NamedQuery(query, schedule)
    if err != nil {
        return fmt.Errorf("failed to create schedule: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)
        if err != nil {
            return fmt.Errorf("failed to scan created schedule: %w", err)
        }
    }

    return nil
}

func (r *ScheduleRepository) Update(schedule *dtos.ScheduleDTO) error {
    query := `
        UPDATE schedules SET
            schedule_name = :schedule_name,
            template_id = :template_id,
            email_to = :email_to,
            email_cc = :email_cc,
            email_bcc = :email_bcc,
            frequency = :frequency,
            schedule_days = :schedule_days,
            start_date = :start_date,
            end_date = :end_date,
            start_time = :start_time,
            timezone = :timezone,
            is_active = :is_active,
            search_criteria = :search_criteria,
            export_format = :export_format,
            updated_by = :updated_by,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = :id AND is_deleted = false
    `

    result, err := r.db.NamedExec(query, schedule)
    if err != nil {
        return fmt.Errorf("failed to update schedule: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("schedule not found")
    }

    return nil
}

func (r *ScheduleRepository) Delete(id int) error {
    query := `UPDATE schedules SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete schedule: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("schedule not found")
    }

    return nil
}

func (r *ScheduleRepository) GetActiveSchedules() ([]dtos.ScheduleDTO, error) {
    var schedules []dtos.ScheduleDTO
    query := `
        SELECT s.*, t.template_name
        FROM schedules s
        LEFT JOIN templates t ON s.template_id = t.id
        WHERE s.is_active = true AND s.is_deleted = false
        AND (s.end_date IS NULL OR s.end_date >= CURRENT_DATE)
        ORDER BY s.next_run_at ASC
    `
    err := r.db.Select(&schedules, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get active schedules: %w", err)
    }
    return schedules, nil
}

func (r *ScheduleRepository) UpdateLastRun(id int) error {
    query := `
        UPDATE schedules SET
            last_run_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `
    _, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to update last run: %w", err)
    }
    return nil
}

// LogRepository handles log data operations
type LogRepository struct {
    db *sqlx.DB
}

func NewLogRepository(db *sqlx.DB) *LogRepository {
    return &LogRepository{db: db}
}

func (r *LogRepository) GetSentReportLogs(req dtos.LogSearchRequestDTO) ([]dtos.SentReportLogDTO, int64, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1

    baseQuery := `
        SELECT 
            l.*,
            t.template_name,
            s.schedule_name
        FROM sent_report_logs l
        LEFT JOIN templates t ON l.template_id = t.id
        LEFT JOIN schedules s ON l.schedule_id = s.id
        WHERE 1=1
    `

    countQuery := `
        SELECT COUNT(*)
        FROM sent_report_logs l
        WHERE 1=1
    `

    // Add search conditions
    if req.TemplateID != nil {
        conditions = append(conditions, fmt.Sprintf("l.template_id = $%d", argIndex))
        args = append(args, *req.TemplateID)
        argIndex++
    }

    if req.ScheduleID != nil {
        conditions = append(conditions, fmt.Sprintf("l.schedule_id = $%d", argIndex))
        args = append(args, *req.ScheduleID)
        argIndex++
    }

    if req.Status != "" {
        conditions = append(conditions, fmt.Sprintf("l.status = $%d", argIndex))
        args = append(args, req.Status)
        argIndex++
    }

    if req.DateFrom != nil {
        conditions = append(conditions, fmt.Sprintf("l.sent_at >= $%d", argIndex))
        args = append(args, req.DateFrom.Format("2006-01-02"))
        argIndex++
    }

    if req.DateTo != nil {
        conditions = append(conditions, fmt.Sprintf("l.sent_at <= $%d", argIndex))
        args = append(args, req.DateTo.Format("2006-01-02 23:59:59"))
        argIndex++
    }

    // Add conditions to queries
    if len(conditions) > 0 {
        conditionStr := " AND " + strings.Join(conditions, " AND ")
        baseQuery += conditionStr
        countQuery += conditionStr
    }

    // Get total count
    var total int64
    err := r.db.Get(&total, countQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get log count: %w", err)
    }

    // Add pagination
    if req.Limit == 0 {
        req.Limit = 10
    }
    offset := (req.Page - 1) * req.Limit
    if offset < 0 {
        offset = 0
    }

    baseQuery += fmt.Sprintf(" ORDER BY l.sent_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, req.Limit, offset)

    // Execute query
    var logs []dtos.SentReportLogDTO
    err = r.db.Select(&logs, baseQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get sent report logs: %w", err)
    }

    return logs, total, nil
}

func (r *LogRepository) GetByID(id int) (*dtos.SentReportLogDTO, error) {
    var log dtos.SentReportLogDTO
    query := `
        SELECT 
            l.*,
            t.template_name,
            s.schedule_name
        FROM sent_report_logs l
        LEFT JOIN templates t ON l.template_id = t.id
        LEFT JOIN schedules s ON l.schedule_id = s.id
        WHERE l.id = $1
    `
    err := r.db.Get(&log, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get sent report log by ID: %w", err)
    }
    return &log, nil
}

func (r *LogRepository) Create(log *dtos.SentReportLogDTO) error {
    query := `
        INSERT INTO sent_report_logs (
            template_id, schedule_id, recipients, subject, file_name,
            file_size_kb, export_format, total_records, status,
            error_message, retry_count, execution_time_ms
        ) VALUES (
            :template_id, :schedule_id, :recipients, :subject, :file_name,
            :file_size_kb, :export_format, :total_records, :status,
            :error_message, :retry_count, :execution_time_ms
        ) RETURNING id, sent_at
    `

    rows, err := r.db.NamedQuery(query, log)
    if err != nil {
        return fmt.Errorf("failed to create sent report log: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&log.ID, &log.SentAt)
        if err != nil {
            return fmt.Errorf("failed to scan created log: %w", err)
        }
    }

    return nil
}

func (r *LogRepository) UpdateStatus(id int, status string, errorMessage *string) error {
    query := `
        UPDATE sent_report_logs SET
            status = $2,
            error_message = $3
        WHERE id = $1
    `
    _, err := r.db.Exec(query, id, status, errorMessage)
    if err != nil {
        return fmt.Errorf("failed to update log status: %w", err)
    }
    return nil
}

// FieldRepository handles field data operations
type FieldRepository struct {
    db *sqlx.DB
}

func NewFieldRepository(db *sqlx.DB) *FieldRepository {
    return &FieldRepository{db: db}
}

func (r *FieldRepository) GetAllFields() ([]dtos.AvailableFieldDTO, error) {
    var fields []dtos.AvailableFieldDTO
    query := `
        SELECT * FROM available_fields 
        WHERE is_active = true 
        ORDER BY field_category, sort_order
    `
    err := r.db.Select(&fields, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get available fields: %w", err)
    }
    return fields, nil
}

func (r *FieldRepository) GetFieldsByCategory(category string) ([]dtos.AvailableFieldDTO, error) {
    var fields []dtos.AvailableFieldDTO
    query := `
        SELECT * FROM available_fields 
        WHERE field_category = $1 AND is_active = true 
        ORDER BY sort_order
    `
    err := r.db.Select(&fields, query, category)
    if err != nil {
        return nil, fmt.Errorf("failed to get fields by category: %w", err)
    }
    return fields, nil
}

func (r *FieldRepository) GetFieldByCode(fieldCode string) (*dtos.AvailableFieldDTO, error) {
    var field dtos.AvailableFieldDTO
    query := `SELECT * FROM available_fields WHERE field_code = $1 AND is_active = true`
    err := r.db.Get(&field, query, fieldCode)
    if err != nil {
        return nil, fmt.Errorf("failed to get field by code: %w", err)
    }
    return &field, nil
}

func (r *FieldRepository) ValidateFields(fieldCodes []string) ([]dtos.FieldValidationDTO, error) {
    var results []dtos.FieldValidationDTO
    
    for _, code := range fieldCodes {
        var count int
        query := `SELECT COUNT(*) FROM available_fields WHERE field_code = $1 AND is_active = true`
        err := r.db.Get(&count, query, code)
        if err != nil {
            results = append(results, dtos.FieldValidationDTO{
                FieldCode: code,
                IsValid:   false,
                Message:   "Database error: " + err.Error(),
            })
            continue
        }
        
        if count > 0 {
            results = append(results, dtos.FieldValidationDTO{
                FieldCode: code,
                IsValid:   true,
                Message:   "Field is valid",
            })
        } else {
            results = append(results, dtos.FieldValidationDTO{
                FieldCode: code,
                IsValid:   false,
                Message:   "Field code not found or inactive",
            })
        }
    }
    
    return results, nil
}

func (r *FieldRepository) GetFieldCategories() ([]string, error) {
    var categories []string
    query := `
        SELECT DISTINCT field_category 
        FROM available_fields 
        WHERE is_active = true 
        ORDER BY field_category
    `
    err := r.db.Select(&categories, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get field categories: %w", err)
    }
    return categories, nil
}

func (r *FieldRepository) CreateField(field *dtos.AvailableFieldDTO) error {
    query := `
        INSERT INTO available_fields (
            field_code, field_name_thai, field_name_eng, field_type,
            field_category, data_source, format_example, is_required,
            sort_order, description
        ) VALUES (
            :field_code, :field_name_thai, :field_name_eng, :field_type,
            :field_category, :data_source, :format_example, :is_required,
            :sort_order, :description
        ) RETURNING id
    `

    rows, err := r.db.NamedQuery(query, field)
    if err != nil {
        return fmt.Errorf("failed to create field: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&field.ID)
        if err != nil {
            return fmt.Errorf("failed to scan created field: %w", err)
        }
    }

    return nil
}

func (r *FieldRepository) UpdateField(field *dtos.AvailableFieldDTO) error {
    query := `
        UPDATE available_fields SET
            field_name_thai = :field_name_thai,
            field_name_eng = :field_name_eng,
            field_type = :field_type,
            field_category = :field_category,
            data_source = :data_source,
            format_example = :format_example,
            is_required = :is_required,
            is_active = :is_active,
            sort_order = :sort_order,
            description = :description
        WHERE id = :id
    `

    result, err := r.db.NamedExec(query, field)
    if err != nil {
        return fmt.Errorf("failed to update field: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("field not found")
    }

    return nil
}

func (r *FieldRepository) DeleteField(id int) error {
    query := `UPDATE available_fields SET is_active = false WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete field: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("field not found")
    }

    return nil
}

func (r *FieldRepository) ExistsFieldCode(code string) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM available_fields WHERE field_code = $1`
    err := r.db.Get(&count, query, code)
    if err != nil {
        return false, fmt.Errorf("failed to check field code existence: %w", err)
    }
    return count > 0, nil
}

func (r *FieldRepository) GetRequiredFields() ([]dtos.AvailableFieldDTO, error) {
    var fields []dtos.AvailableFieldDTO
    query := `
        SELECT * FROM available_fields 
        WHERE is_required = true AND is_active = true 
        ORDER BY field_category, sort_order
    `
    err := r.db.Select(&fields, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get required fields: %w", err)
    }
    return fields, nil
}

func (r *FieldRepository) GetFieldsByType(fieldType string) ([]dtos.AvailableFieldDTO, error) {
    var fields []dtos.AvailableFieldDTO
    query := `
        SELECT * FROM available_fields 
        WHERE field_type = $1 AND is_active = true 
        ORDER BY field_category, sort_order
    `
    err := r.db.Select(&fields, query, fieldType)
    if err != nil {
        return nil, fmt.Errorf("failed to get fields by type: %w", err)
    }
    return fields, nil
}

func (r *FieldRepository) GetFieldsForExport(fieldCodes []string) ([]dtos.AvailableFieldDTO, error) {
    if len(fieldCodes) == 0 {
        return []dtos.AvailableFieldDTO{}, nil
    }

    // Create placeholders for the IN clause
    placeholders := make([]string, len(fieldCodes))
    args := make([]interface{}, len(fieldCodes))
    for i, code := range fieldCodes {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
        args[i] = code
    }

    query := fmt.Sprintf(`
        SELECT * FROM available_fields 
        WHERE field_code IN (%s) AND is_active = true 
        ORDER BY field_category, sort_order
    `, strings.Join(placeholders, ","))

    var fields []dtos.AvailableFieldDTO
    err := r.db.Select(&fields, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to get fields for export: %w", err)
    }
    return fields, nil
}