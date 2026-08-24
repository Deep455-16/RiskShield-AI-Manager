# Security Model

## Authentication
- Argon2id password hashing (m=65536, t=3, p=4)
- JWT access tokens (15 min expiry)
- JWT refresh tokens (7 day expiry)
- Token revocation support

## Authorization
- RBAC with 7 roles: SUPER_ADMIN, ORG_ADMIN, RISK_MANAGER, SECURITY_ANALYST, AI_ENGINEER, AUDITOR, VIEWER
- Tenant isolation enforced server-side in every query

## Data Protection
- PII detection and redaction
- Sensitive data never logged
- Masked display: `**** **** **** 1111`

## Network Security
- CORS allowlist
- Rate limiting ready
- No secrets in source code
- Non-root Docker containers

## Threat Model
See THREAT_MODEL.md for detailed analysis of:
- Malicious users, compromised accounts
- Prompt injection, data exfiltration
- Privilege escalation, API abuse
- Insider threats, compromised vendors
