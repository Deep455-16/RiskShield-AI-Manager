# RiskShield AI — Architecture

## Design Principles

1. **Deterministic over Probabilistic** — Risk scores and policy decisions must be reproducible
2. **Evidence-Based** — Every score must reference measurable evidence
3. **Human Accountability** — AI recommends, humans decide on high-impact actions
4. **Tenant Isolation** — Server-side enforcement, never trust client-side auth
5. **Immutable Audit** — Tamper-evident log chain with cryptographic hashing

## Component Architecture

### Risk Engine
- 7 configurable dimensions: Fraud, Security, Privacy, Compliance, Fairness, Reliability, Operational
- Weighted sum with configurable organization-specific weights
- Risk bands: LOW (0-24), MEDIUM (25-49), HIGH (50-74), CRITICAL (75-100)
- Every score contains: score, confidence, factors, evidence, timestamp, version

### Policy Engine
- Deterministic condition evaluation
- Actions: allow, review, block
- Versioned, auditable, immutable after approval
- Default: CRITICAL → block, HIGH → review

### Agent Security
- Zero-trust: no inherited permissions
- Per-action risk scoring
- Tool allowlist, API allowlist, budget limits
- Kill-switch capability

### Audit System
- SHA-256 hash chain: each entry includes previous hash
- Every security-sensitive state change is logged
- No secrets or raw PII in logs
