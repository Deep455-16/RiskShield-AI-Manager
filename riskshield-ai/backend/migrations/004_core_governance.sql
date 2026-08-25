-- Migration: Core Governance (Modules 1-7)

-- 1. AI System Registry (Update existing ai_systems table)
-- Existing ai_systems has risk_class. Let's add version, owner_id, deployment_env, risk_tier.
ALTER TABLE ai_systems ADD COLUMN IF NOT EXISTS version VARCHAR(50) DEFAULT '1.0';
ALTER TABLE ai_systems ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id);
ALTER TABLE ai_systems ADD COLUMN IF NOT EXISTS deployment_env VARCHAR(50) DEFAULT 'production';
ALTER TABLE ai_systems ADD COLUMN IF NOT EXISTS risk_tier VARCHAR(20) DEFAULT 'medium';

-- 2. Risk Engine Extension
CREATE TABLE IF NOT EXISTS risk_scores_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ai_system_id UUID NOT NULL REFERENCES ai_systems(id) ON DELETE CASCADE,
    score NUMERIC(5, 2) NOT NULL,
    computed_at TIMESTAMPTZ DEFAULT NOW(),
    breakdown JSONB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_risk_scores_history_sys ON risk_scores_history(ai_system_id, computed_at DESC);

-- 3. Cross-Framework Compliance Register (SoA)
-- Update existing compliance_controls table
ALTER TABLE compliance_controls ADD COLUMN IF NOT EXISTS framework VARCHAR(50) DEFAULT 'custom';
ALTER TABLE compliance_controls ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'not_implemented';
ALTER TABLE compliance_controls ADD COLUMN IF NOT EXISTS ai_system_id UUID REFERENCES ai_systems(id) ON DELETE CASCADE;
ALTER TABLE compliance_controls ADD COLUMN IF NOT EXISTS evidence_url TEXT;

-- 4. Assessment Workflows
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    ai_system_id UUID NOT NULL REFERENCES ai_systems(id) ON DELETE CASCADE,
    framework VARCHAR(50) NOT NULL,
    framework_version VARCHAR(50),
    status VARCHAR(50) NOT NULL, -- draft, in_review, completed, approved
    assessor_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assessment_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50),
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Continuous Monitoring
CREATE TABLE IF NOT EXISTS monitoring_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    ai_system_id UUID NOT NULL REFERENCES ai_systems(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    operator VARCHAR(10) NOT NULL, -- gt, lt, eq
    threshold NUMERIC(10, 2) NOT NULL,
    time_window VARCHAR(20) NOT NULL, -- e.g., '5m', '1h'
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Incident & CAPA Management
CREATE TABLE IF NOT EXISTS capa_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL, -- corrective, preventive
    description TEXT NOT NULL,
    owner_id UUID REFERENCES users(id),
    status VARCHAR(50) NOT NULL, -- open, in_progress, closed
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add SLA field to incidents
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS sla_deadline TIMESTAMPTZ;

-- 7. Reports (Audit trail for generated reports)
CREATE TABLE IF NOT EXISTS generated_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    scope VARCHAR(50) NOT NULL, -- framework, asset, org
    format VARCHAR(20) NOT NULL, -- json, pdf
    file_url TEXT,
    content_hash VARCHAR(64) NOT NULL, -- SHA-256
    generated_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
