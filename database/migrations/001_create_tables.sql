-- Providers table with all fields
CREATE TABLE IF NOT EXISTS providers (
    id SERIAL PRIMARY KEY,
    
    -- Basic Info
    provider_type VARCHAR(100) NOT NULL DEFAULT 'Hospital',
    status VARCHAR(50) NOT NULL DEFAULT 'Private',
    provider_name_th VARCHAR(255) NOT NULL,
    provider_name_en VARCHAR(255),
    telephone_number VARCHAR(50),
    address TEXT,
    district VARCHAR(100),
    province VARCHAR(100),
    post_code VARCHAR(20),
    
    -- Extended Info
    provider_code VARCHAR(50) UNIQUE,
    title_th VARCHAR(100),
    name_th VARCHAR(255),
    title_en VARCHAR(100),
    name_en VARCHAR(255),
    register_status VARCHAR(100),
    business_type VARCHAR(100),
    bed_size VARCHAR(50),
    eligibility_method VARCHAR(100),
    region VARCHAR(100),
    country VARCHAR(100) DEFAULT 'ประเทศไทย',
    provider_tax_id VARCHAR(20),
    wh_tax_percent DECIMAL(5,2) DEFAULT 0,
    exempt_percent DECIMAL(5,2) DEFAULT 0,
    wh_tax_exempt_from DATE,
    wh_tax_exempt_to DATE,
    opening_time VARCHAR(100),
    provider_status VARCHAR(50) DEFAULT 'Active',
    building_no VARCHAR(100),
    village_no VARCHAR(100),
    lane_alley VARCHAR(100),
    road VARCHAR(100),
    sub_district VARCHAR(100),
    
    -- Contact Info
    contact_title_name VARCHAR(255),
    contact_department VARCHAR(100),
    general_phone_no VARCHAR(50),
    direct_phone_no VARCHAR(50),
    email VARCHAR(255),
    email_to_list TEXT,
    email_cc_list TEXT,
    
    -- Payment Info
    payment_method VARCHAR(50),
    payment_branch_id VARCHAR(20),
    payee_name VARCHAR(255),
    bank_account_number VARCHAR(50),
    bank_account_type VARCHAR(50),
    bank_branch_name VARCHAR(100),
    bank_name VARCHAR(100),
    
    -- Network & Incident
    is_tpa_network BOOLEAN DEFAULT FALSE,
    has_incident BOOLEAN DEFAULT FALSE,
    
    -- Discount Categories
    discount_category1 VARCHAR(100),
    discount_category2 VARCHAR(100),
    discount_category3 VARCHAR(100),
    discount_category4 VARCHAR(100),
    discount_category5 VARCHAR(100),
    
    -- Pricing Categories
    pricing_category1 VARCHAR(100),
    pricing_category2 VARCHAR(100),
    pricing_category3 VARCHAR(100),
    pricing_category4 VARCHAR(100),
    pricing_category5 VARCHAR(100),
    
    -- System fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Report Templates table
CREATE TABLE IF NOT EXISTS report_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    report_type VARCHAR(50) NOT NULL DEFAULT 'provider_detail',
    is_default BOOLEAN DEFAULT FALSE,
    header_fields JSONB,
    detail_fields JSONB,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled Reports table
CREATE TABLE IF NOT EXISTS scheduled_reports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_id INTEGER REFERENCES report_templates(id),
    cron_expression VARCHAR(100) NOT NULL,
    email_to JSONB,
    email_cc JSONB,
    search_criteria JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_run TIMESTAMP,
    next_run TIMESTAMP
);

-- Report Logs table
CREATE TABLE IF NOT EXISTS report_logs (
    id SERIAL PRIMARY KEY,
    report_type VARCHAR(50) NOT NULL,
    template_id INTEGER,
    scheduled_report_id INTEGER,
    status VARCHAR(50) NOT NULL, -- success, failed, running
    file_path VARCHAR(500),
    file_size BIGINT,
    export_format VARCHAR(20),
    search_criteria JSONB,
    error_message TEXT,
    execution_time INTEGER, -- seconds
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_providers_name ON providers(provider_name_th);
CREATE INDEX IF NOT EXISTS idx_providers_province ON providers(province);
CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(provider_type);
CREATE INDEX IF NOT EXISTS idx_providers_status ON providers(provider_status);
CREATE INDEX IF NOT EXISTS idx_providers_created_at ON providers(created_at);
CREATE INDEX IF NOT EXISTS idx_report_templates_type ON report_templates(report_type);
CREATE INDEX IF NOT EXISTS idx_scheduled_reports_active ON scheduled_reports(is_active);
CREATE INDEX IF NOT EXISTS idx_report_logs_created_at ON report_logs(created_at);

-- Sample data
INSERT INTO providers (
    provider_name_th, provider_name_en, provider_type, status, province, district,
    telephone_number, address, post_code, provider_code, business_type,
    email, is_tpa_network, provider_status
) VALUES 
(
    'โรงพยาบาลบำรุงราษฎร์', 'Bumrungrad Hospital', 'Hospital', 'Private',
    'กรุงเทพมหานคร', 'สาทร', '02-667-1000', '33 ถนนสุขุมวิท แขวงคลองเตย',
    '10110', 'PVR001', 'Private', 'info@bumrungrad.com', TRUE, 'Active'
),
(
    'โรงพยาบาลสมิติเวช', 'Samitivej Hospital', 'Hospital', 'Private',
    'กรุงเทพมหานคร', 'ปทุมวัน', '02-022-2222', '488 ถนนพระราม 1',
    '10330', 'PVR002', 'Private', 'info@samitivej.co.th', TRUE, 'Active'
),
(
    'คลินิกเวชกรรมทั่วไป', 'General Medical Clinic', 'Clinic', 'Private',
    'เชียงใหม่', 'เมืองเชียงใหม่', '053-123-456', '123 ถนนนิมมานเหมินท์',
    '50200', 'PVR003', 'Private', 'info@clinic.com', FALSE, 'Active'
);

-- Default template
INSERT INTO report_templates (
    name, description, report_type, is_default, 
    header_fields, detail_fields, created_by
) VALUES (
    'Standard Provider Report', 
    'Default template for provider detail report',
    'provider_detail',
    TRUE,
    '["H1", "H2"]'::jsonb,
    '["D1", "D4", "D5", "D6", "D8", "D9", "D11", "D13", "D14", "D17"]'::jsonb,
    'system'
);