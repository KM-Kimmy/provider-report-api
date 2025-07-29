package dtos

type AvailableFieldDTO struct {
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

type FieldListResponseDTO struct {
    Fields []AvailableFieldDTO `json:"fields"`
    Total  int                 `json:"total"`
}

type FieldsByCategoryResponseDTO struct {
    Category string              `json:"category"`
    Fields   []AvailableFieldDTO `json:"fields"`
    Total    int                 `json:"total"`
}

type FieldCategoriesResponseDTO struct {
    Categories []string `json:"categories"`
    Total      int      `json:"total"`
}

type FieldValidationDTO struct {
    FieldCode string `json:"field_code"`
    IsValid   bool   `json:"is_valid"`
    Message   string `json:"message,omitempty"`
}

type FieldsValidationRequestDTO struct {
    FieldCodes []string `json:"field_codes" binding:"required,min=1"`
}

type FieldsValidationResponseDTO struct {
    Results []FieldValidationDTO `json:"results"`
    Summary struct {
        TotalFields  int `json:"total_fields"`
        ValidFields  int `json:"valid_fields"`
        InvalidFields int `json:"invalid_fields"`
    } `json:"summary"`
}

type CreateFieldRequestDTO struct {
    FieldCode       string  `json:"field_code" binding:"required"`
    FieldNameThai   string  `json:"field_name_thai" binding:"required"`
    FieldNameEng    string  `json:"field_name_eng" binding:"required"`
    FieldType       string  `json:"field_type" binding:"required,oneof=text numeric date boolean"`
    FieldCategory   string  `json:"field_category" binding:"required,oneof=header summary detail"`
    DataSource      *string `json:"data_source"`
    FormatExample   *string `json:"format_example"`
    IsRequired      bool    `json:"is_required"`
    SortOrder       *int    `json:"sort_order"`
    Description     *string `json:"description"`
}

type UpdateFieldRequestDTO struct {
    FieldNameThai   string  `json:"field_name_thai" binding:"required"`
    FieldNameEng    string  `json:"field_name_eng" binding:"required"`
    FieldType       string  `json:"field_type" binding:"required,oneof=text numeric date boolean"`
    FieldCategory   string  `json:"field_category" binding:"required,oneof=header summary detail"`
    DataSource      *string `json:"data_source"`
    FormatExample   *string `json:"format_example"`
    IsRequired      bool    `json:"is_required"`
    IsActive        bool    `json:"is_active"`
    SortOrder       *int    `json:"sort_order"`
    Description     *string `json:"description"`
}