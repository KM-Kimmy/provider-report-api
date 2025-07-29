package models

import "time"

type SentReportLog struct {
    ID              int        `json:"id" db:"id"`
    TemplateID      int        `json:"template_id" db:"template_id"`
    ScheduleID      *int       `json:"schedule_id" db:"schedule_id"`
    Recipients      string     `json:"recipients" db:"recipients"`
    Subject         *string    `json:"subject" db:"subject"`
    FileName        *string    `json:"file_name" db:"file_name"`
    FileSizeKB      *int       `json:"file_size_kb" db:"file_size_kb"`
    ExportFormat    *string    `json:"export_format" db:"export_format"`
    TotalRecords    *int       `json:"total_records" db:"total_records"`
    SentAt          time.Time  `json:"sent_at" db:"sent_at"`
    Status          string     `json:"status" db:"status"`
    ErrorMessage    *string    `json:"error_message" db:"error_message"`
    RetryCount      int        `json:"retry_count" db:"retry_count"`
    ExecutionTimeMs *int       `json:"execution_time_ms" db:"execution_time_ms"`
    
    // Joined fields
    TemplateName    string     `json:"template_name" db:"template_name"`
    ScheduleName    *string    `json:"schedule_name" db:"schedule_name"`
}

type LogSearchRequest struct {
    TemplateID   *int       `json:"template_id" form:"template_id"`
    ScheduleID   *int       `json:"schedule_id" form:"schedule_id"`
    Status       string     `json:"status" form:"status"`
    DateFrom     *time.Time `json:"date_from" form:"date_from"`
    DateTo       *time.Time `json:"date_to" form:"date_to"`
    Page         int        `json:"page" form:"page"`
    Limit        int        `json:"limit" form:"limit"`
}