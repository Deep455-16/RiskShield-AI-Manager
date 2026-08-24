CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'VIEWER',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, email)
);

CREATE TABLE ai_systems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    purpose TEXT,
    owner TEXT,
    risk_class TEXT,
    data_class TEXT,
    deployment_status TEXT DEFAULT 'development',
    approval_status TEXT DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    version TEXT,
    type TEXT,
    framework TEXT,
    owner TEXT,
    status TEXT DEFAULT 'registered',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    owner TEXT,
    purpose TEXT,
    model TEXT,
    version TEXT,
    environment TEXT,
    risk_level TEXT DEFAULT 'medium',
    approval_status TEXT DEFAULT 'pending',
    tools JSONB DEFAULT '[]',
    permissions JSONB DEFAULT '[]',
    allowed_apis JSONB DEFAULT '[]',
    data_access_scope JSONB DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE agent_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    tool TEXT,
    target TEXT,
    risk_score FLOAT,
    decision TEXT,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    transaction_id TEXT NOT NULL,
    amount FLOAT,
    currency TEXT,
    merchant_id TEXT,
    customer_id TEXT,
    device_id TEXT,
    ip_address TEXT,
    location TEXT,
    velocity INT DEFAULT 0,
    model_probability FLOAT DEFAULT 0,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE risk_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL,
    asset_type TEXT NOT NULL,
    overall_score FLOAT NOT NULL,
    risk_level TEXT NOT NULL,
    confidence FLOAT,
    dimensions JSONB,
    contributing_factors JSONB,
    evidence JSONB,
    scoring_version TEXT,
    model_version TEXT,
    policy_version TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_risk_scores_org_asset ON risk_scores(org_id, asset_id, timestamp DESC);

CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    condition JSONB,
    action TEXT NOT NULL,
    priority INT DEFAULT 0,
    version INT DEFAULT 1,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    status TEXT DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    severity TEXT,
    category TEXT,
    affected_system_id UUID,
    affected_user_id TEXT,
    root_cause TEXT,
    evidence JSONB,
    risk_score FLOAT,
    owner TEXT,
    sla_hours INT,
    status TEXT DEFAULT 'OPEN',
    remediation TEXT,
    capa TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE compliance_frameworks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    version TEXT,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE compliance_controls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    framework_id UUID REFERENCES compliance_frameworks(id),
    control_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'NOT_IMPLEMENTED',
    evidence JSONB,
    mapped_frameworks JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id UUID,
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    ip_address TEXT,
    request_id TEXT,
    before_state JSONB,
    after_state JSONB,
    hash TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_org ON audit_logs(org_id, timestamp DESC);

CREATE TABLE approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    resource_type TEXT NOT NULL,
    resource_id UUID NOT NULL,
    risk_score_id UUID REFERENCES risk_scores(id),
    decision TEXT,
    status TEXT DEFAULT 'pending',
    reviewer_id UUID REFERENCES users(id),
    reviewer_comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP
);

CREATE TABLE vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    service TEXT,
    ai_capability TEXT,
    data_access TEXT,
    region TEXT,
    compliance_certs JSONB,
    security_posture TEXT,
    incident_history JSONB,
    contract_expiry DATE,
    risk_score FLOAT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    payload JSONB,
    result JSONB,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
