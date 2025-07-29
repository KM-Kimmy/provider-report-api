package models

import (
    "time"
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
    return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
    if value == nil {
        *s = nil
        return nil
    }
    
    switch v := value.(type) {
    case []byte:
        return json.Unmarshal(v, s)
    case string:
        return json.Unmarshal([]byte(v), s)
    }
    return fmt.Errorf("cannot scan %T into StringArray", value)
}

type ReportTemplate struct {
    ID              int         `json:"id" db:"id"`
    Name            string      `json:"name" db:"name"`
    Description     string      `json:"description" db:"description"`
    ReportType      string      `json:"report_type" db:"report_type"`
    IsDefault       bool        `json:"is_default" db:"is_default"`
    HeaderFields    StringArray `json:"header_fields" db:"header_fields"`
    DetailFields    StringArray `json:"detail_fields" db:"detail_fields"`
    CreatedBy       string      `json:"created_by" db:"created_by"`
    CreatedAt       time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

type FieldDefinition struct {
    FieldCode   string `json:"field_code"`
    FieldName   string `json:"field_name"`
    Description string `json:"description"`
    DataType    string `json:"data_type"`
    Format      string `json:"format,omitempty"`
    Example     string `json:"example"`
    Category    string `json:"category"`
}

type ProviderSearchRequest struct {
    Name            string   `json:"name" form:"name"`
    Province        string   `json:"province" form:"province"`
    ProviderType    string   `json:"provider_type" form:"provider_type"`
    Status          string   `json:"status" form:"status"`
    DateFrom        string   `json:"date_from" form:"date_from"`
    DateTo          string   `json:"date_to" form:"date_to"`
    Page            int      `json:"page" form:"page"`
    Limit           int      `json:"limit" form:"limit"`
    SortBy          string   `json:"sort_by" form:"sort_by"`
    SortOrder       string   `json:"sort_order" form:"sort_order"`
    TemplateID      *int     `json:"template_id" form:"template_id"`
    HeaderFields    []string `json:"header_fields" form:"header_fields"`
    DetailFields    []string `json:"detail_fields" form:"detail_fields"`
}

type ExportRequest struct {
    Format         string `json:"format" form:"format"`
    TemplateID     *int   `json:"template_id" form:"template_id"`
    TemplateName   string `json:"template_name" form:"template_name"`
    SaveAsTemplate bool   `json:"save_as_template" form:"save_as_template"`
    ProviderSearchRequest
}

type TemplateRequest struct {
    Name           string   `json:"name"`
    Description    string   `json:"description"`
    ReportType     string   `json:"report_type"`
    IsDefault      bool     `json:"is_default"`
    HeaderFields   []string `json:"header_fields"`
    DetailFields   []string `json:"detail_fields"`
    CreatedBy      string   `json:"created_by"`
}

type FieldsResponse struct {
    HeaderFields []FieldDefinition `json:"header_fields"`
    DetailFields []FieldDefinition `json:"detail_fields"`
}

type ProviderResponse struct {
    Data       []Provider `json:"data"`
    Total      int64      `json:"total"`
    Page       int        `json:"page"`
    Limit      int        `json:"limit"`
    TotalPages int64      `json:"total_pages"`
    Message    string     `json:"message,omitempty"`
}
