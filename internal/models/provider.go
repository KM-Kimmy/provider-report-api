package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
)

type JSONStringArray []string

func (j *JSONStringArray) Scan(value interface{}) error {
    if value == nil {
        *j = JSONStringArray{}
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    
    return json.Unmarshal(bytes, j)
}

func (j JSONStringArray) Value() (driver.Value, error) {
    if len(j) == 0 {
        return "[]", nil
    }
    return json.Marshal(j)
}

type Provider struct {
    ID                  int               `json:"id" db:"id"`
    ProviderCode        string            `json:"provider_code" db:"provider_code"`
    TitleThai           *string           `json:"title_thai" db:"title_thai"`
    NameThai            string            `json:"name_thai" db:"name_thai"`
    TitleEng            *string           `json:"title_eng" db:"title_eng"`
    NameEng             *string           `json:"name_eng" db:"name_eng"`
    ProviderType        string            `json:"provider_type" db:"provider_type"`
    RegisterStatus      *string           `json:"register_status" db:"register_status"`
    BusinessType        *string           `json:"business_type" db:"business_type"`
    BedSize             *string           `json:"bed_size" db:"bed_size"`
    EligibilityMethod   *string           `json:"eligibility_method" db:"eligibility_method"`
    Province            string            `json:"province" db:"province"`
    Region              *string           `json:"region" db:"region"`
    Country             *string           `json:"country" db:"country"`
    ProviderTaxID       *string           `json:"provider_tax_id" db:"provider_tax_id"`
    WHTaxPercent        *float64          `json:"wh_tax_percent" db:"wh_tax_percent"`
    ExemptPercent       *float64          `json:"exempt_percent" db:"exempt_percent"`
    WHTaxExemptFrom     *time.Time        `json:"wh_tax_exempt_from" db:"wh_tax_exempt_from"`
    WHTaxExemptTo       *time.Time        `json:"wh_tax_exempt_to" db:"wh_tax_exempt_to"`
    OpeningTime         *string           `json:"opening_time" db:"opening_time"`
    ProviderStatus      string            `json:"provider_status" db:"provider_status"`
    BuildingNo          *string           `json:"building_no" db:"building_no"`
    VillageNo           *string           `json:"village_no" db:"village_no"`
    LaneAlley           *string           `json:"lane_alley" db:"lane_alley"`
    Road                *string           `json:"road" db:"road"`
    SubDistrict         *string           `json:"sub_district" db:"sub_district"`
    District            *string           `json:"district" db:"district"`
    PostCode            *string           `json:"post_code" db:"post_code"`
    TitleName           *string           `json:"title_name" db:"title_name"`
    Department          *string           `json:"department" db:"department"`
    GeneralPhoneNo      *string           `json:"general_phone_no" db:"general_phone_no"`
    DirectPhoneNo       *string           `json:"direct_phone_no" db:"direct_phone_no"`
    Email               *string           `json:"email" db:"email"`
    EmailToList         *string           `json:"email_to_list" db:"email_to_list"`
    EmailCCList         *string           `json:"email_cc_list" db:"email_cc_list"`
    PaymentMethod       *string           `json:"payment_method" db:"payment_method"`
    PaymentBranchID     *string           `json:"payment_branch_id" db:"payment_branch_id"`
    PayeeName           *string           `json:"payee_name" db:"payee_name"`
    BankAccountNumber   *string           `json:"bank_account_number" db:"bank_account_number"`
    BankAccountType     *string           `json:"bank_account_type" db:"bank_account_type"`
    BankBranchName      *string           `json:"bank_branch_name" db:"bank_branch_name"`
    BankName            *string           `json:"bank_name" db:"bank_name"`
    IsTPANetwork        bool              `json:"is_tpa_network" db:"is_tpa_network"`
    HasIncident         bool              `json:"has_incident" db:"has_incident"`
    DiscountCategories  JSONStringArray   `json:"discount_categories" db:"discount_categories"`
    PricingCategories   JSONStringArray   `json:"pricing_categories" db:"pricing_categories"`
    CreatedAt           time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at" db:"updated_at"`
    CreatedBy           *string           `json:"created_by" db:"created_by"`
    UpdatedBy           *string           `json:"updated_by" db:"updated_by"`
}

type ProviderReportRequest struct {
    SearchParams ProviderSearchRequest `json:"search_params"`
    TemplateID   *int                  `json:"template_id"`
    FormatType   string                `json:"format_type"` // excel, pdf, word
    CustomFields []string              `json:"custom_fields,omitempty"`
}

type ProviderSummary struct {
    Type       string `json:"type"`
    Hospital   int    `json:"hospital"`
    Clinic     int    `json:"clinic"`
    GrandTotal int    `json:"grand_total"`
    Province   string `json:"province,omitempty"`
}

type ProviderReportData struct {
    Header   map[string]interface{} `json:"header"`
    Summary  ProviderSummary        `json:"summary"`
    Providers []Provider            `json:"providers"`
    Total    int64                  `json:"total"`
}