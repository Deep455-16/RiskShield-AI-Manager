-- Seed demo organization and frameworks
INSERT INTO organizations (id, name, slug) VALUES
('11111111-1111-1111-1111-111111111111', 'Demo Fintech Corp', 'demo-fintech')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO compliance_frameworks (id, org_id, name, version, description) VALUES
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'NIST AI RMF', '1.0', 'NIST AI Risk Management Framework'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'EU AI Act', '2024/1689', 'European Union Artificial Intelligence Act'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'ISO/IEC 42001', '2023', 'AI Management Systems');

INSERT INTO compliance_controls (id, org_id, control_id, title, description, status) VALUES
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'GOV-1', 'AI System Inventory', 'Maintain inventory of all AI systems', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'GOV-2', 'Risk Assessment', 'Conduct regular risk assessments', 'PARTIALLY_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'GOV-3', 'Human Oversight', 'Ensure human oversight for high-risk AI', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'GOV-4', 'Data Governance', 'Implement data governance practices', 'NOT_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'GOV-5', 'Transparency', 'Ensure AI system transparency', 'PARTIALLY_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'SEC-1', 'Model Security', 'Protect AI models from attacks', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'SEC-2', 'Prompt Injection Defense', 'Implement prompt injection detection', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'SEC-3', 'PII Protection', 'Protect personally identifiable information', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'SEC-4', 'Agent Security', 'Secure AI agent operations', 'PARTIALLY_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'FAIR-1', 'Bias Detection', 'Detect and mitigate algorithmic bias', 'NOT_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'FAIR-2', 'Fairness Metrics', 'Measure fairness across demographic groups', 'NOT_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'MON-1', 'Model Monitoring', 'Continuously monitor model performance', 'IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'MON-2', 'Drift Detection', 'Detect data and concept drift', 'PARTIALLY_IMPLEMENTED'),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'MON-3', 'Incident Response', 'Establish AI incident response procedures', 'IMPLEMENTED');
