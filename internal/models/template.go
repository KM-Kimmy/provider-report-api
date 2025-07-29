package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
)

type JSONFieldArray []string

func (j *JSONFieldArray) Scan(value interface{}) error {
    if value == nil {
        *j = JSONFieldArray{}
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    
    return json.Unmarshal(bytes, j)
}

func (j JSONFieldArray) Value() (driver.Value, error) {
    if len(j) == 0 {
        return "[]", nil
    }
    return json.Marshal(j)
}

type Template struct {
    ID             int             `json:"id" db:"id"`
    TemplateName   string          `json:"template_name" db:"template_name"`
    IsStandard     bool            `json:"is_standard" db:"is_standard"`
    Description    *string         `json:"description" db:"description"`
    HeaderFields   JSONFieldArray  `json:"header_fields" db:"header_fields"`
    DataFields     JSONFieldArray  `json:"data_fields" db:"data_fields"`
    SummaryFields  JSONFieldArray  `json:"summary_fields" db:"summary_fields"`
    FieldPositions *string         `json:"field_positions" db:"field_positions"`
    CreatedAt      time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
    CreatedBy      string          `json:"created_by" db:"created_by"`
    UpdatedBy      *string         `json:"updated_by" db:"updated_by"`
    IsDeleted      bool            `json:"is_deleted" db:"is_deleted"`
}

type CreateTemplateRequest struct {
    TemplateName   string   `json:"template_name" binding:"required"`
    IsStandard     bool     `json:"is_standard"`
    Description    string   `json:"description"`
    HeaderFields   []string `json:"header_fields"`
    DataFields     []string `json:"data_fields" binding:"required"`
    SummaryFields  []string `json:"summary_fields"`
    FieldPositions string   `json:"field_positions"`
}

type UpdateTemplateRequest struct {
    TemplateName   string   `json:"template_name" binding:"required"`
    IsStandard     bool     `json:"is_standard"`
    Description    string   `json:"description"`
    HeaderFields   []string `json:"header_fields"`
    DataFields     []string `json:"data_fields" binding:"required"`
    SummaryFields  []string `json:"summary_fields"`
    FieldPositions string   `json:"field_positions"`
}

type AvailableField struct {
    ID              int     `json:"id" db:"id"`
    FieldCode       string  `json:"field_code" db:"field_code"`
    FieldNameThai   string  `json:"field_name_thai" db:"field_name_thai"`
    FieldNameEng    string  `json:"field_name_eng" db:"field_name_eng"`
    FieldType       string  `json:"field_type" db:"field_type"`
    FieldCategory   string  `json:"field_category" db:"field_category"`
    DataSource      *string `json:"data_source" db:"data_source"`
    FormatExample   *string `json:"format_example" db:"format_example"`
    IsRequired      bool    `json:"is_required" db:"is_required"`
    IsActive        bool    `json:"is_active" db:"is_active"`
    SortOrder       *int    `json:"sort_order" db:"sort_order"`
    Description     *string `json:"description" db:"description"`
}