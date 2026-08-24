# RiskShield AI

> **Detect AI risk. Explain it. Control it. Prevent it.**

RiskShield AI is a production-grade, open-source AI Risk Management and AI Governance platform designed for fintech and enterprise environments. It continuously monitors, evaluates, explains, and controls risks arising from AI models, AI agents, payment/fraud systems, APIs, data, and autonomous decisions.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    RiskShield AI Platform                    │
├─────────────────────────────────────────────────────────────┤
│  Frontend (Next.js 15 + React + TypeScript + Tailwind)    │
├─────────────────────────────────────────────────────────────┤
│  Go Backend (Chi Router + PostgreSQL + JWT + Argon2id)      │
│  ├── Auth & RBAC (7 roles, tenant isolation)                │
│  ├── Risk Engine (deterministic weighted scoring)           │
│  ├── Policy Engine (deterministic rule evaluation)          │
│  ├── AI System Registry                                     │
│  ├── Agent Registry & Action Monitoring                     │
│  ├── Incident & CAPA Management                             │
│  ├── Compliance (NIST AI RMF, EU AI Act, ISO 42001)         │
│  ├── Audit Logs (SHA-256 tamper-evident chaining)           │
│  └── Approval Workflows (human-in-the-loop)                 │
├─────────────────────────────────────────────────────────────┤
│  Python AI Service (FastAPI)                                │
│  ├── PII Detection & Redaction                              │
│  ├── Prompt Injection Detection                           │
│  ├── Fraud Assessment                                     │
│  ├── Fairness Evaluation                                  │
│  ├── Model Drift Detection                                │
│  └── LLM Evaluation                                       │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL 16  │  Redis  │  Docker Compose                 │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Clone and start
git clone https://github.com/your-org/riskshield-ai.git
cd riskshield-ai

# Start all services
docker compose up --build

# Access:
# Dashboard:  http://localhost:3000
# API Docs:   http://localhost:8080/health
# AI Service: http://localhost:8001/health
```

**Demo Login:**
- Email: `admin@riskshield.demo`
- Password: `DemoAdmin123!`

## Core Features

| Feature | Description |
|---------|-------------|
| **Risk Engine** | Deterministic 7-dimension weighted scoring with evidence |
| **Policy Engine** | Immutable, versioned, auditable policy enforcement |
| **Human-in-the-Loop** | Approval workflows for high-impact decisions |
| **Agent Governance** | Zero-trust agent registry with per-action risk scoring |
| **Payment Simulator** | Interactive fintech risk demonstration |
| **Attack Simulator** | 8 attack types with before/after risk visualization |
| **PII Detection** | Pattern-based detection with automatic redaction |
| **Prompt Injection** | Multi-layer defense with deterministic patterns |
| **Compliance** | NIST AI RMF, EU AI Act, ISO 42001 mappings |
| **Audit Logs** | Tamper-evident SHA-256 chained logging |
| **Multi-Tenancy** | Organization-scoped data with server-side isolation |
| **RBAC** | 7 roles with 30+ granular permissions |

## Security Principles

- **Zero Trust** — Never trust, always verify
- **Least Privilege** — Agents inherit nothing by default
- **Defense in Depth** — Layered security controls
- **Fail Closed** — Default deny on policy failure
- **Deterministic Authority** — LLM analyzes, policy engine decides

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 15, React, TypeScript, Tailwind CSS, Recharts |
| Backend | Go 1.23, Chi, pgx/v5, JWT, Argon2id |
| AI Service | Python 3.12, FastAPI, scikit-learn |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Auth | JWT (15min access / 7day refresh), Argon2id |

## API Endpoints

```
GET  /health
GET  /ready
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout

GET  /api/v1/ai-systems
POST /api/v1/ai-systems
GET  /api/v1/models
POST /api/v1/models
GET  /api/v1/agents
POST /api/v1/agents
POST /api/v1/agents/{id}/actions

GET  /api/v1/transactions
POST /api/v1/transactions
POST /api/v1/transactions/{id}/assess

GET  /api/v1/risk/scores
POST /api/v1/risk/assess
GET  /api/v1/policies
POST /api/v1/policies

GET  /api/v1/incidents
POST /api/v1/incidents
GET  /api/v1/compliance/frameworks
GET  /api/v1/compliance/controls
GET  /api/v1/audit-logs
GET  /api/v1/approvals

POST /api/v1/simulate/payment
POST /api/v1/simulate/attack
POST /api/v1/copilot/explain
GET  /api/v1/dashboard
```

## License

Apache 2.0

## Attribution

Inspired by the architecture and principles of:
- [Assuro](https://github.com/YASSERRMD/Assuro) — AI governance platform
- [AIAF Sentry](https://github.com/mbwika/AI-Assurance-Framework) — AI assurance framework
- [GRITS](https://github.com/X-Scale-AI/GRITS) — AI agent security framework
- [OpenDeRisk](https://github.com/derisk-ai/OpenDerisk) — AI-native risk intelligence
