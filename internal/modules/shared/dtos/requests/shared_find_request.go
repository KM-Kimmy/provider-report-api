package request

type ProvinceQueryDto struct {
	RegionId *string `form:"regionId"`
}

type DistrictQueryDto struct {
	ProvinceId *string `form:"provinceId"`
}
type SubDistrictQueryDto struct {
	DistrictId *string `form:"districtId"`
}

type SaveSearchDto struct {
	SubModule string `json:"subModule" form:"subModule" validate:"required"`
}

type RegionQueryDto struct {
	CountryId *int `form:"countryId"`
}

type CountryQueryDto struct {
	CountryName *string `form:"countryName"`
	CountryCode *string `form:"countryCode"`
}

type AdminfeeQueryDto struct {
	InsurerId *int `form:"insurerId"`
}

type NetworkQueryDto struct {
	NetworkName *string `form:"networkName"`
}

type AuditLog struct {
	Id        string `json:"id"`        // Unique log Id
	Action    string `json:"action"`    // add, edit, delete
	Username  string `json:"username"`  // username of the user performing the action
	Data      any    `json:"data"`      // The data affected (can be before/after for edits)
	Timestamp string `json:"timestamp"` // Action timestamp
	Level     string `json:"level"`     // info, error, warn, debug
	Service   string `json:"service"`   // which microservice/produced the log
	SubModule string `json:"subModule"` // query from sub-module
}