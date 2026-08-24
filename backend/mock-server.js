const http = require('http');

const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const data = {
  "/api/v1/dashboard": {
    active_incidents: 2,
    critical_risks: 5,
    policy_violations_24h: 18,
    ai_systems_monitored: 12,
    compliance_score: 92,
    transactions_today: 3450,
    agents_active: 8,
    audit_events_24h: 890
  },
  "/api/v1/risk/scores": {
    data: [
      { id: '1', asset_id: 'fdtct-001-xgb', asset_type: 'agent', overall_score: 82, risk_level: 'CRITICAL', timestamp: new Date(Date.now() - 300000).toISOString() },
      { id: '2', asset_id: 'crdt-scr-002', asset_type: 'model', overall_score: 67, risk_level: 'HIGH', timestamp: new Date(Date.now() - 900000).toISOString() },
      { id: '3', asset_id: 'pii-rdc-003', asset_type: 'agent', overall_score: 34, risk_level: 'MEDIUM', timestamp: new Date(Date.now() - 1800000).toISOString() },
      { id: '4', asset_id: 'anml-dtct-004', asset_type: 'system', overall_score: 71, risk_level: 'HIGH', timestamp: new Date(Date.now() - 3600000).toISOString() }
    ]
  },
  "/api/v1/ai-systems": {
    data: [
      { id: '1', name: 'Fraud Detection Engine', description: 'Real-time transaction fraud scoring', purpose: 'Fraud Prevention', owner: 'Risk Team', risk_class: 'high', approval_status: 'approved', deployment_status: 'production' },
      { id: '2', name: 'Customer Support Bot', description: 'Automated L1 support', purpose: 'Support', owner: 'CS Team', risk_class: 'medium', approval_status: 'approved', deployment_status: 'production' }
    ]
  },
  "/api/v1/transactions": {
    data: [
      { id: '1', transaction_id: 'TXN-98213', amount: 45000, currency: 'USD', status: 'approved', created_at: new Date(Date.now() - 120000).toISOString() },
      { id: '2', transaction_id: 'TXN-98214', amount: 87500, currency: 'EUR', status: 'blocked', created_at: new Date(Date.now() - 240000).toISOString() },
      { id: '3', transaction_id: 'TXN-98215', amount: 1250, currency: 'USD', status: 'review', created_at: new Date(Date.now() - 360000).toISOString() }
    ]
  },
  "/api/v1/agents": {
    data: [
      { id: '1', name: 'CreditScorer-v3', purpose: 'Automated credit limit adjustments', model: 'gpt-4', environment: 'production', risk_level: 'high', approval_status: 'approved' },
      { id: '2', name: 'SupportAgent-Gen2', purpose: 'Customer email resolution', model: 'claude-3-sonnet', environment: 'production', risk_level: 'medium', approval_status: 'approved' }
    ]
  },
  "/api/v1/incidents": {
    data: [
      { id: '1', title: 'Prompt Injection in SupportAgent', description: 'User attempted to bypass system prompt', severity: 'HIGH', status: 'OPEN' },
      { id: '2', title: 'Model Drift: FraudDetect-v2', description: 'Accuracy dropped by 4% in last 24h', severity: 'MEDIUM', status: 'RESOLVED' }
    ]
  },
  "/api/v1/policies": {
    data: [
      { id: '1', name: 'Block High Risk Transactions', description: 'Automatically block transactions with risk > 80', version: 3, status: 'active', action: 'BLOCK', priority: 100 },
      { id: '2', name: 'Review Novel Devices', description: 'Require review for txns from new devices over $1k', version: 1, status: 'active', action: 'REVIEW', priority: 50 }
    ]
  },
  "/api/v1/compliance/controls": {
    data: [
      { id: '1', control_id: 'SEC-1', title: 'Model Security', description: 'Protect AI models from attacks', status: 'IMPLEMENTED' },
      { id: '2', control_id: 'GOV-2', title: 'Risk Assessment', description: 'Conduct regular risk assessments', status: 'PARTIALLY_IMPLEMENTED' },
      { id: '3', control_id: 'FAIR-1', title: 'Bias Detection', description: 'Detect and mitigate algorithmic bias', status: 'NOT_IMPLEMENTED' }
    ]
  },
  "/api/v1/audit-logs": {
    data: [
      { id: '1', action: 'policy_updated', resource_type: 'policy', hash: 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', timestamp: new Date(Date.now() - 300000).toISOString() },
      { id: '2', action: 'agent_action_blocked', resource_type: 'agent_action', hash: '8a2b1c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b', timestamp: new Date(Date.now() - 600000).toISOString() }
    ]
  }
};

const server = http.createServer(async (req, res) => {
  // CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
  res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type,Authorization');

  if (req.method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }

  // Handle Auth
  if (req.url === '/api/v1/auth/login' && req.method === 'POST') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ access_token: 'mock-jwt-token' }));
    return;
  }

  // Handle Simulators
  if (req.url.startsWith('/api/v1/simulate/')) {
    await delay(1000);
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      risk_score: { overall_score: 85, risk_level: 'CRITICAL' },
      policy_decision: { decision: 'BLOCK' },
      before: { score: 25, level: 'LOW' },
      after: { score: 85, level: 'CRITICAL' },
      incident_created: 'INC-' + Math.floor(Math.random() * 10000)
    }));
    return;
  }

  // Standard GET endpoints
  const path = req.url.split('?')[0];
  if (data[path]) {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(data[path]));
  } else {
    // Default empty array for unknown GETs
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ data: [] }));
  }
});

const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Mock Backend server listening on port ${PORT}`);
});
