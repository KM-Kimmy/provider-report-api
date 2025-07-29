package dtos

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
)

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
    if value == nil {
        *j = JSONMap{}
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    
    return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
    if len(j) == 0 {
        return "{}", nil
    }
    return json.Marshal(j)
}

type ScheduleDTO struct {
    ID             int             `json:"id" db:"id"`
    ScheduleName   string          `json:"schedule_name" db:"schedule_name"`
    TemplateID     int             `json:"template_id" db:"template_id"`
    EmailTo        string          `json:"email_to" db:"email_to"`
    EmailCC        *string         `json:"email_cc" db:"email_cc"`
    EmailBCC       *string         `json:"email_bcc" db:"email_bcc"`
    Frequency      string          `json:"frequency" db:"frequency"`
    ScheduleDays   JSONFieldArray  `json:"schedule_days" db:"schedule_days"`
    StartDate      time.Time       `json:"start_date" db:"start_date"`
    EndDate        *time.Time      `json:"end_date" db:"end_date"`
    StartTime      string          `json:"start_time" db:"start_time"`
    Timezone       string          `json:"timezone" db:"timezone"`
    IsActive       bool            `json:"is_active" db:"is_active"`
    LastRunAt      *time.Time      `json:"last_run_at" db:"last_run_at"`
    NextRunAt      *time.Time      `json:"next_run_at" db:"next_run_at"`
    SearchCriteria JSONMap         `json:"search_criteria" db:"search_criteria"`
    ExportFormat   string          `json:"export_format" db:"export_format"`
    CreatedAt      time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
    CreatedBy      string          `json:"created_by" db:"created_by"`
    UpdatedBy      *string         `json:"updated_by" db:"updated_by"`
    IsDeleted      bool            `json:"is_deleted" db:"is_deleted"`
    
    // Joined fields
    TemplateName   string          `json:"template_name" db:"template_name"`
}

type CreateScheduleRequestDTO struct {
    ScheduleName   string                 `json:"schedule_name" binding:"required"`
    TemplateID     int                    `json:"template_id" binding:"required"`
    EmailTo        string                 `json:"email_to" binding:"required,email"`
    EmailCC        string                 `json:"email_cc"`
    EmailBCC       string                 `json:"email_bcc"`
    Frequency      string                 `json:"frequency" binding:"required,oneof=daily weekly monthly"`
    ScheduleDays   []string               `json:"schedule_days"`
    StartDate      time.Time              `json:"start_date" binding:"required"`
    EndDate        *time.Time             `json:"end_date"`
    StartTime      string                 `json:"start_time" binding:"required"`
    Timezone       string                 `json:"timezone"`
    SearchCriteria map[string]interface{} `json:"search_criteria"`
    ExportFormat   string                 `json:"export_format" binding:"oneof=excel pdf word"`
}

type UpdateScheduleRequestDTO struct {
    ScheduleName   string                 `json:"schedule_name" binding:"required"`
    TemplateID     int                    `json:"template_id" binding:"required"`
    EmailTo        string                 `json:"email_to" binding:"required,email"`
    EmailCC        string                 `json:"email_cc"`
    EmailBCC       string                 `json:"email_bcc"`
    Frequency      string                 `json:"frequency" binding:"required,oneof=daily weekly monthly"`
    ScheduleDays   []string               `json:"schedule_days"`
    StartDate      time.Time              `json:"start_date" binding:"required"`
    EndDate        *time.Time             `json:"end_date"`
    StartTime      string                 `json:"start_time" binding:"required"`
    Timezone       string                 `json:"timezone"`
    IsActive       bool                   `json:"is_active"`
    SearchCriteria map[string]interface{} `json:"search_criteria"`
    ExportFormat   string                 `json:"export_format" binding:"oneof=excel pdf word"`
}

type ScheduleListResponseDTO struct {
    Schedules []ScheduleDTO `json:"schedules"`
    Total     int           `json:"total"`
}

type RunScheduleResponseDTO struct {
    Message      string    `json:"message"`
    ExecutedAt   time.Time `json:"executed_at"`
    Recipients   string    `json:"recipients"`
    RecordCount  int       `json:"record_count"`
    FileSize     string    `json:"file_size"`
    Status       string    `json:"status"`
}