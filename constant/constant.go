package constant

const DEFAULT_PAGE_SIZE = 20
const DEFAULT_LIMIT_RECORDS = 10000

const (
	TRANSACTION_TYPE_NEW         = "New"
	TRANSACTION_TYPE_ENDORSEMENT = "Endorsement"
)

const AGENT_IMPORT_INDEX = "client-agent-import"
const PRICING_IMPORT_INDEX = "client-pricing-import"
const DISCOUNT_IMPORT_INDEX = "client-discount-import"

// EXPORT INSURER
const (
	EXPORT_AGENT_ERROR        = "Error Message"
	EXPORT_AGENT_SEQUENCE     = "SEQUENCE"
	EXPORT_AGENT_INSURER_CODE = "Insurer Code"
	EXPORT_LEADER_AGENT_CODE  = "Agent Leader Code"
	EXPORT_LEADER_AGENT_NAME  = "Agent Leader Name"
	EXPORT_LEADER_MOBILE      = "Agent Leader Mobile"
	EXPORT_LEADER_EMAIL       = "Agent Leader Email"
	EXPORT_LEADER_REMARKS     = "Agent Leader Remarks"

	EXPORT_AGENT_AGENT_CODE        = "Agent Code"
	EXPORT_AGENT_AGENT_NAME        = "Agent Name"
	EXPORT_AGENT_MOBILE            = "Agent Mobile"
	EXPORT_AGENT_EMAIL             = "Agent Email"
	EXPORT_AGENT_REMARKS           = "Agent Remarks"
	EXPORT_AGENT_AGENT_LEADER_CODE = "Agent Leader Code"
)

const (
	CLIENT_MGMT_INDEX = "user-management-history"
	// client-management-history-paymentdetail
)

// EXPORT PROVIDER
const (
	EXPORT_ERROR         = "Error Message"
	EXPORT_SEQUENCE      = "SEQUENCE"
	EXPORT_PROVIDER_CODE = "Provider Code"
	EXPORT_FROMDATE      = "From Date"
	EXPORT_TODATE        = "To Date"
	EXPORT_ASOFDATE      = "As of Date"
	EXPORT_REMARKS       = "Remarks"

	//pricing
	EXPORT_PRICING_CATEGORY = "Pricing Category"
	EXPORT_PRICE            = "Price"

	//discount
	EXPORT_DISCOUNT_CATEGORY = "Discount Category"
	EXPORT_DISCOUNT          = "Discount %"
)
