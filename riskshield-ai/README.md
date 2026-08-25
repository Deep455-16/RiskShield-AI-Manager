<div align="center"> 🛡️ RiskShield AI
Detect AI risk. Explain it. Control it. Prevent it.
A production-grade, open-source AI Risk Management & Governance platform for fintech and enterprise environments.
</div> RiskShield AI continuously monitors, evaluates, explains, and controls risk arising from AI models, AI agents, payment/fraud systems, APIs, data pipelines, and autonomous decisions — giving fintech and enterprise teams a single pane of glass for AI governance.

## 🏗️ Architecture
```mermaid
flowchart TB 
UI["🖥️ Client Layer<br/>Next.js 15 · React · TypeScript · Tailwind · Recharts"] 
GO["🔀 API Gateway<br/>Go 1.23 Backend — Chi Router<br/>JWT (15m/7d) · Argon2id · rate limiting"] 
subgraph CORE["⚙️ Core Governance Services (Go)"] 
AUTH["🔐 Auth & RBAC<br/>7 roles · tenant isolation"] 
RISK["🎯 Risk Engine<br/>deterministic 7-dim scoring"] 
POLICY["📜 Policy Engine<br/>versioned, immutable rules"] 
REG["🤖 AI & Agent Registry<br/>action monitoring"] 
INC["🚨 Incident & CAPA"] 
COMP["📋 Compliance Mapper<br/>NIST · EU AI Act · ISO 42001"] 
AUDIT["🔗 Audit Log<br/>SHA-256 chained"] 
APPR["🙋 Approval Workflows<br/>human-in-the-loop"] 
end 
subgraph AISVC["🧠 AI Analysis Service (Python · FastAPI)"] 
PII["🕵️ PII Detection"] 
INJ["🛡️ Prompt Injection"] 
FRAUD["💳 Fraud Assessment"] 
FAIR["⚖️ Fairness Evaluation"] 
DRIFT["📉 Model Drift"] 
EVAL["🧪 LLM Evaluation"] 
end 
PG[("🐘 PostgreSQL 16<br/>system of record")] 
REDIS[("🔴 Redis 7<br/>cache · sessions")] 

UI -->|HTTPS / REST| GO 
GO --> CORE 
CORE <-->|internal API| AISVC 
CORE <--> PG 
CORE <--> REDIS 
style UI fill:#0f2744,stroke:#4cc9f0,color:#fff 
style GO fill:#0a2426,stroke:#00c2c2,color:#fff 
style CORE fill:#171029,stroke:#b58cff,color:#fff 
style AISVC fill:#0f261e,stroke:#5cf0a0,color:#fff 
style PG fill:#26170d,stroke:#f0a860,color:#fff 
style REDIS fill:#26170d,stroke:#f0a860,color:#fff
```

⚖️ **Deterministic authority, always.** The AI service analyzes and scores — the Policy Engine is the sole decision-maker. This keeps every allow/block/escalate outcome explainable, reproducible, and auditable.
Request flow: Client → Go API Gateway (auth + routing) → Core Governance Services → (optionally) AI Service for scoring/detection → Policy Engine renders a deterministic decision → Audit Log chains the outcome → PostgreSQL persists state, Redis handles caching/session data.

| Layer | Responsibility |
| --- | --- |
| 🖥️ Client | Next.js 15 dashboard, simulators, and compliance UI |
| 🔀 API Gateway | Go/Chi backend — authentication, routing, rate limiting |
| ⚙️ Core Governance Services | Auth/RBAC, Risk Engine, Policy Engine, Registry, Incidents, Compliance, Audit, Approvals |
| 🧠 AI Analysis Service | Stateless Python/FastAPI microservices for PII, injection, fraud, fairness, drift, and LLM evaluation |
| 💾 Data Layer | PostgreSQL 16 (system of record) + Redis 7 (cache, sessions, rate limits) |

## 🚀 Quick Start
```bash
# Clone and start
git clone https://github.com/your-org/riskshield-ai.git 
cd riskshield-ai 
# Start all services 
docker compose up --build
```

| Service | URL |
| --- | --- |
| 🖥️ Dashboard | http://localhost:3000 |
| ⚙️ API Health | http://localhost:8080/health |
| 🧠 AI Service Health | http://localhost:8001/health |

**Demo Login**
Email: `admin@riskshield.demo`
Password: `DemoAdmin123!`

## ✨ Core Features
| Feature | Description |
| --- | --- |
| 🎯 Risk Engine | Deterministic 7-dimension weighted scoring with full evidence trail |
| 📜 Policy Engine | Immutable, versioned, auditable policy enforcement |
| 🙋 Human-in-the-Loop | Approval workflows for high-impact / high-risk decisions |
| 🤖 Agent Governance | Zero-trust agent registry with per-action risk scoring |
| 💳 Payment Simulator | Interactive fintech risk demonstration sandbox |
| ⚔️ Attack Simulator | 8 attack types with before/after risk visualization |
| 🕵️ PII Detection | Pattern-based detection with automatic redaction |
| 🛡️ Prompt Injection Defense | Multi-layer, deterministic pattern detection |
| 📋 Compliance Mapping | NIST AI RMF, EU AI Act, ISO 42001 |
| 🔗 Tamper-Evident Audit Logs | SHA-256 chained logging |
| 🏢 Multi-Tenancy | Organization-scoped data with server-side isolation |
| 🔐 RBAC | 7 roles, 30+ granular permissions |

## 🧩 Feature Deep Dive

### Core Governance
<details> <summary><strong>🗂️ AI System Registry</strong></summary> Register every AI system, pipeline, or model in your organization with structured metadata — name, version, description, owner, deployment environment, and purpose. Each system is assigned a risk classification tier based on EU AI Act categories (minimal, limited, high, unacceptable) and linked to its underlying assets.
</details> 
<details> <summary><strong>🎯 Risk Engine</strong></summary> Computes a weighted risk score for each asset based on assessment findings, incident history, and monitoring data. Scores are stored historically so you can track how risk evolves over time. The dashboard surfaces the current risk tier distribution (critical, high, medium, low, unknown) across all assets in your organization.
</details> 
<details> <summary><strong>📋 Cross-Framework Compliance Register</strong></summary> A Statement of Applicability (SoA) system that maps assets and controls to multiple regulatory frameworks simultaneously. Each control record tracks implementation status (implemented, partially implemented, not applicable, not implemented) and links to evidence artefacts. Compliance score is calculated automatically as:
score = (implemented + 0.5 × partial) / total × 100
Supported frameworks: EU AI Act 2024 · NIST AI Risk Management Framework 1.0 · ISO/IEC 42001:2023
</details> 
<details> <summary><strong>📝 Assessment Workflows</strong></summary> Structured questionnaires tied to a specific AI system and framework version. Each assessment moves through a defined lifecycle: draft → in-review → completed → approved. Assessors attach evidence files (documents, screenshots, API responses), and the platform tracks which version of the framework each assessment was conducted against.
</details> 
<details> <summary><strong>📡 Continuous Monitoring</strong></summary> Define threshold-based monitoring rules for metrics such as performance scores, data drift indicators, and uptime. When a metric breaches a threshold, the platform creates an alert and can trigger notifications or webhook deliveries to your existing tooling.
</details> 
<details> <summary><strong>🚨 Incident and CAPA Management</strong></summary> Log AI-related incidents with severity levels, root cause categories, and affected systems. Attach corrective and preventive action (CAPA) records to each incident, assign owners, and track resolution status. Unresolved incidents older than the configured SLA automatically surface in the analytics dashboard.
</details> 
<details> <summary><strong>📄 Report Generation</strong></summary> Generate compliance reports in JSON or PDF format. Each report includes a SHA-256 hash of its content so you can prove to auditors the document has not been modified after generation. Reports can be scoped to a single framework, a single asset, or the entire organization.
</details> 

### Extended Modules
<details> <summary><strong>🤖 Agent Registry & Runtime</strong></summary> Register autonomous AI agents alongside human-operated AI systems. The registry tracks agent type (autonomous, assistant, pipeline, custom), status, and capability scope. The runtime layer adds:
Behaviour logging with outcome classification (ok, error, blocked, anomaly)
Guardrail policies with configurable actions (block, warn, log)
Anomaly records with severity tiers (low, medium, high, critical) and resolution tracking
A kill-switch endpoint that immediately suspends an agent and writes a timestamped audit entry
</details> 
<details> <summary><strong>🕵️ Shadow AI Discovery</strong></summary> Organizations routinely run AI systems that were never formally registered. The discovery module scans cloud AI platforms and code repositories via a connector framework. Discovered items are written to a deduplication inbox using SHA-256 keys on (org, source, external-id) tuples, preventing duplicate entries on repeated scans. Reconcile inbox items against the asset registry to identify genuine shadow AI.
Built-in connectors:
| Connector | What it scans |
| --- | --- |
| AWS Bedrock | Foundation models in the configured region |
| Azure AI | Models deployed in the configured workspace |
| GCP Vertex AI | Models in the configured project and location |
| GitHub | Repos with Python/Jupyter files or AI keywords (ml, llm, gpt, bert, embedding, inference) |
| HuggingFace Hub | Public models associated with the configured organization |

All connectors implement a static fallback that returns representative models when credentials are absent, keeping tests deterministic without mocking.
</details> 
<details> <summary><strong>🏢 Vendor Risk Management</strong></summary> Maintain an inventory of third-party AI vendors categorized by type (ai_provider, data_provider, platform, tools, custom). Create vendor assessments with risk scoring — the platform automatically updates the vendor's risk tier to match the latest assessment result. Track renewal dates, contract metadata, and open assessment findings per vendor.
</details> 
<details> <summary><strong>📜 Policy Management</strong></summary> Create and manage AI governance policies through a structured workflow: draft → in-review → approved → archived. Each approved policy version records attestations from designated users, providing a signed acceptance trail. Policies link to relevant frameworks, assets, and controls for full traceability.
</details> 
<details> <summary><strong>🪪 Model Cards</strong></summary> Generate and version structured model documentation following the model card standard — intended use, limitations, performance metrics, training data, ethical considerations, and maintenance information. Cards are versioned and linked to the corresponding AI system and asset records.
</details> 
<details> <summary><strong>🧪 Model Testing</strong></summary> Define test suites by category (functional, bias, robustness, security, custom) and attach individual test cases with expected outcomes. When a test run completes, the platform calculates a percentage score (passed / total × 100) and stores per-case results (pass, fail, error, skip). Failed test runs can trigger alerts and feed into risk score adjustments.
</details> 
<details> <summary><strong>🙋 Approval Gates</strong></summary> Configure multi-stakeholder approval workflows for high-stakes decisions such as deploying a new AI system or approving a risk exception. Each workflow specifies required approvers and a minimum approval count. The gate transitions automatically — any single rejection marks the request as rejected; reaching the required approval count marks it as approved. All decisions are immutable and written to the audit log.
</details> 
<details> <summary><strong>✅ Task Management</strong></summary> Create and assign remediation and review tasks linked to any resource in the system (asset, incident, finding, vendor). Tasks carry priority levels (critical, high, medium, low), due dates, status (open, in-progress, blocked, done, cancelled), and a comment thread. The task list is sorted by priority using a deterministic weight function.
</details> 
<details> <summary><strong>📊 Analytics Dashboard</strong></summary> A single API endpoint aggregates metrics across all modules into a structured snapshot:
Risk tier distribution from the latest score per asset
Compliance score from the SoA control status breakdown
Incident trend by month for the past six months
Task summary by status and overdue count
Scalar counts for active agents, open shadow AI findings, active vendors, and approved policies
</details> 
<details> <summary><strong>🔗 Audit Log</strong></summary> Every state change writes a structured event with actor, target type, target ID, action, and a JSON payload. The log supports filtering by action, target type, target ID, actor, and date range. Export the filtered result as CSV for submission to external auditors or SIEM ingestion.
</details> 
<details> <summary><strong>🔐 Role-Based Access Control</strong></summary> In addition to built-in system roles (owner, admin, member), create custom org-scoped roles and assign fine-grained permissions. 30+ permission scopes cover every domain in the platform.
</details> 

## 🔒 Security Principles
- **Zero Trust** — never trust, always verify
- **Least Privilege** — agents inherit nothing by default
- **Defense in Depth** — layered security controls
- **Fail Closed** — default deny on policy failure
- **Deterministic Authority** — LLM analyzes, policy engine decides

## 🧰 Tech Stack
| Layer | Technology |
| --- | --- |
| Frontend | Next.js 15, React, TypeScript, Tailwind CSS, Recharts |
| Backend | Go 1.23, Chi, pgx/v5, JWT, Argon2id |
| AI Service | Python 3.12, FastAPI, scikit-learn |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Auth | JWT (15 min access / 7 day refresh), Argon2id |

## 📡 API Endpoints
<details> <summary><strong>Click to expand full API reference</strong></summary>
GET /health 
GET /ready 
POST /api/v1/auth/register 
POST /api/v1/auth/login 
POST /api/v1/auth/refresh 
POST /api/v1/auth/logout 
GET /api/v1/ai-systems 
POST /api/v1/ai-systems 
GET /api/v1/models 
POST /api/v1/models 
GET /api/v1/agents 
POST /api/v1/agents 
POST /api/v1/agents/{id}/actions 
GET /api/v1/transactions 
POST /api/v1/transactions 
POST /api/v1/transactions/{id}/assess 
GET /api/v1/risk/scores 
POST /api/v1/risk/assess 
GET /api/v1/policies 
POST /api/v1/policies 
GET /api/v1/incidents 
POST /api/v1/incidents 
GET /api/v1/compliance/frameworks 
GET /api/v1/compliance/controls 
GET /api/v1/audit-logs 
GET /api/v1/approvals 
POST /api/v1/simulate/payment 
POST /api/v1/simulate/attack 
POST /api/v1/copilot/explain 
GET /api/v1/dashboard
</details> 

## 📄 License
Apache 2.0

## 🙏 Attribution
Inspired by the architecture and principles of:
- [Assuro](https://github.com/YASSERRMD/Assuro) — AI governance platform
- [AIAF Sentry](https://github.com/mbwika/AI-Assurance-Framework) — AI assurance framework
- [GRITS](https://github.com/X-Scale-AI/GRITS) — AI agent security framework
- [OpenDeRisk](https://github.com/derisk-ai/OpenDerisk) — AI-native risk intelligence
