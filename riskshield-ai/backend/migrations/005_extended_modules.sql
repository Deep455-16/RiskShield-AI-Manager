-- Migration: Extended Modules (Phase 2 — Modules 8-13)

-- 8. Agent Registry & Runtime Extensions
ALTER TABLE agents ADD COLUMN IF NOT EXISTS agent_type VARCHAR(50) DEFAULT 'assistant';
ALTER TABLE agents ADD COLUMN IF NOT EXISTS capability_scope JSONB DEFAULT '{}';
ALTER TABLE agents ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'active';

CREATE TABLE IF NOT EXISTS agent_behavior_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    outcome VARCHAR(50) NOT NULL, -- ok, error, blocked, anomaly
    payload JSONB DEFAULT '{}',
    ts TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_agent_behavior_logs_agent ON agent_behavior_logs(agent_id, ts DESC);

CREATE TABLE IF NOT EXISTS guardrail_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    rule TEXT NOT NULL,
    action VARCHAR(20) NOT NULL, -- block, warn, log
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS agent_anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    severity VARCHAR(20) NOT NULL, -- low, medium, high, critical
    description TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 9. Shadow AI Discovery
CREATE TABLE IF NOT EXISTS shadow_ai_inbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    source VARCHAR(100) NOT NULL, -- aws_bedrock, azure_ai, github, huggingface, etc.
    external_id TEXT NOT NULL,
    name TEXT,
    details JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'unreviewed', -- unreviewed, approved, ignored
    discovered_at TIMESTAMPTZ DEFAULT NOW(),
    content_hash TEXT GENERATED ALWAYS AS (encode(sha256((org_id::text || source || external_id)::bytea), 'hex')) STORED,
    UNIQUE(content_hash)
);
CREATE INDEX IF NOT EXISTS idx_shadow_ai_inbox_org ON shadow_ai_inbox(org_id);

-- 10. Vendor Risk Management
-- Update existing vendors table
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS contract_meta JSONB DEFAULT '{}';
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS renewal_date TIMESTAMPTZ;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS risk_tier VARCHAR(20) DEFAULT 'medium';

CREATE TABLE IF NOT EXISTS vendor_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    score NUMERIC(5,2),
    findings JSONB DEFAULT '{}',
    assessed_by UUID REFERENCES users(id),
    assessed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 11. Policy Management (extend existing policies table)
ALTER TABLE policies ADD COLUMN IF NOT EXISTS body TEXT;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS framework_refs JSONB DEFAULT '[]';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS asset_refs JSONB DEFAULT '[]';
ALTER TABLE policies ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;

CREATE TABLE IF NOT EXISTS policy_attestations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    attested_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(policy_id, user_id)
);

-- 12. Model Cards
CREATE TABLE IF NOT EXISTS model_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ai_system_id UUID NOT NULL REFERENCES ai_systems(id) ON DELETE CASCADE,
    version VARCHAR(50) NOT NULL,
    intended_use TEXT,
    limitations TEXT,
    performance_metrics JSONB DEFAULT '{}',
    training_data TEXT,
    ethical_considerations TEXT,
    maintenance_info TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 13. Model Testing
CREATE TABLE IF NOT EXISTS test_suites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ai_system_id UUID NOT NULL REFERENCES ai_systems(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- functional, bias, robustness, security, custom
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suite_id UUID NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    input TEXT NOT NULL,
    expected_outcome TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS test_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suite_id UUID NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    score NUMERIC(5,2),
    run_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS test_case_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    case_id UUID NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    result VARCHAR(20) NOT NULL -- pass, fail, error, skip
);
