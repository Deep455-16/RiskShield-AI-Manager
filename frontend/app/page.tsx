'use client'

import { useEffect, useState } from 'react'
import { AlertTriangle, TrendingUp, Shield, Activity, Users, CheckCircle, Zap, Lock, ArrowRight, ShieldCheck, CreditCard, Crosshair } from 'lucide-react'

interface DashboardData {
  active_incidents: number
  critical_risks: number
  policy_violations_24h: number
  ai_systems_monitored: number
  compliance_score: number
  transactions_today: number
  agents_active: number
  audit_events_24h: number
}

const MOCK_STATS: DashboardData = {
  active_incidents: 0,
  critical_risks: 0,
  policy_violations_24h: 0,
  ai_systems_monitored: 0,
  compliance_score: 100,
  transactions_today: 0,
  agents_active: 0,
  audit_events_24h: 0,
}

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardData>(MOCK_STATS)
  const [incidents, setIncidents] = useState<any[]>([])
  
  useEffect(() => {
    fetch('/api/v1/dashboard')
      .then(r => r.json())
      .then(d => setStats(d))
      .catch(() => setStats(MOCK_STATS))
      
    fetch('/api/v1/incidents')
      .then(r => r.json())
      .then(d => setIncidents(d.data || []))
      .catch(() => {})
  }, [])

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 p-8 pb-20">
      
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Executive Risk Overview</h1>
          <p className="text-slate-500 mt-1">TrustBank Corporate Security & Fraud Prevention</p>
        </div>
        <div className="flex gap-3">
          <div className="px-4 py-2 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-lg text-sm font-semibold flex items-center gap-2 shadow-sm">
            <CheckCircle className="w-4 h-4" />
            System Secure
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 mb-4">
            <h3 className="font-medium text-sm uppercase tracking-wide">Threat Level</h3>
            <Shield className="w-5 h-5 text-blue-600" />
          </div>
          <div className="text-3xl font-bold text-slate-900 mb-1">{stats.critical_risks > 2 ? 'ELEVATED' : 'NOMINAL'}</div>
          <div className="text-sm text-slate-500">{stats.critical_risks} active critical threats</div>
        </div>
        
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 mb-4">
            <h3 className="font-medium text-sm uppercase tracking-wide">Daily Volume Monitored</h3>
            <Activity className="w-5 h-5 text-blue-600" />
          </div>
          <div className="text-3xl font-bold text-slate-900 mb-1">{stats.transactions_today.toLocaleString()}</div>
          <div className="text-sm text-emerald-600 font-medium">+12.4% verified safe</div>
        </div>
        
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 mb-4">
            <h3 className="font-medium text-sm uppercase tracking-wide">AI Interventions</h3>
            <Zap className="w-5 h-5 text-amber-500" />
          </div>
          <div className="text-3xl font-bold text-slate-900 mb-1">{stats.policy_violations_24h}</div>
          <div className="text-sm text-slate-500">fraud attempts blocked (24h)</div>
        </div>

        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm relative overflow-hidden">
          <div className="absolute right-0 top-0 w-32 h-32 bg-blue-50 rounded-bl-full -mr-16 -mt-16 z-0"></div>
          <div className="relative z-10">
            <div className="flex items-center justify-between text-slate-500 mb-4">
              <h3 className="font-medium text-sm uppercase tracking-wide">Regulatory Score</h3>
              <CheckCircle className="w-5 h-5 text-emerald-600" />
            </div>
            <div className="text-3xl font-bold text-slate-900 mb-1">{stats.compliance_score}%</div>
            <div className="text-sm text-slate-500">Fully compliant with Basel III</div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* Left Column: AI Guardrails Explainer */}
        <div className="lg:col-span-2 space-y-8">
          <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
            <div className="px-6 py-5 border-b border-slate-100 bg-slate-50 flex justify-between items-center">
              <h2 className="text-lg font-semibold text-slate-900 flex items-center gap-2">
                <ShieldCheck className="w-5 h-5 text-blue-600" />
                AI Fraud Prevention Guardrails
              </h2>
            </div>
            <div className="p-6">
              <p className="text-slate-600 mb-6 text-sm leading-relaxed">
                The TrustBank AI Risk Manager operates continuously to protect corporate and retail accounts from direct cyber attacks, unauthorized access, and payment fraud. Our multi-dimensional AI scoring engine blocks malicious activity instantly.
              </p>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="p-4 bg-slate-50 border border-slate-100 rounded-lg">
                  <CreditCard className="w-6 h-6 text-blue-600 mb-3" />
                  <h4 className="font-semibold text-slate-900 mb-1">Velocity Tracking</h4>
                  <p className="text-xs text-slate-500">Detects rapid sequence transfers indicating account takeover attempts.</p>
                </div>
                <div className="p-4 bg-slate-50 border border-slate-100 rounded-lg">
                  <Crosshair className="w-6 h-6 text-blue-600 mb-3" />
                  <h4 className="font-semibold text-slate-900 mb-1">IP & Device Novelty</h4>
                  <p className="text-xs text-slate-500">Flags requests from unknown devices or blacklisted IP addresses automatically.</p>
                </div>
                <div className="p-4 bg-slate-50 border border-slate-100 rounded-lg">
                  <Lock className="w-6 h-6 text-blue-600 mb-3" />
                  <h4 className="font-semibold text-slate-900 mb-1">Model Drift Guard</h4>
                  <p className="text-xs text-slate-500">Monitors the AI's own decision accuracy to prevent adversarial manipulation.</p>
                </div>
              </div>
            </div>
          </div>
          
          <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
             <div className="px-6 py-5 border-b border-slate-100 flex justify-between items-center">
              <h2 className="text-lg font-semibold text-slate-900">Recent AI Interventions</h2>
            </div>
            <div className="divide-y divide-slate-100">
              {incidents.slice(0, 4).map(incident => (
                <div key={incident.id} className="p-4 px-6 hover:bg-slate-50 transition-colors">
                  <div className="flex justify-between items-start mb-1">
                    <span className="font-semibold text-slate-900">{incident.title}</span>
                    <span className={`px-2 py-1 rounded text-xs font-medium ${
                      incident.severity === 'CRITICAL' ? 'bg-red-100 text-red-700' :
                      incident.severity === 'HIGH' ? 'bg-orange-100 text-orange-700' :
                      'bg-amber-100 text-amber-700'
                    }`}>
                      {incident.severity}
                    </span>
                  </div>
                  <p className="text-sm text-slate-500">{incident.description}</p>
                  <div className="mt-3 flex items-center gap-4 text-xs text-slate-400 font-medium">
                    <span className="flex items-center gap-1">
                      <Shield className="w-3 h-3" /> Action Taken: Account Locked
                    </span>
                    <span>Status: {incident.status}</span>
                  </div>
                </div>
              ))}
              {incidents.length === 0 && (
                <div className="p-8 text-center text-slate-500 text-sm">No recent incidents detected.</div>
              )}
            </div>
          </div>
        </div>

        {/* Right Column */}
        <div className="space-y-6">
          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
            <h2 className="text-lg font-semibold text-slate-900 mb-4">Live Threat Feed</h2>
            <div className="space-y-4">
              <div className="flex gap-3">
                <div className="w-2 h-2 mt-2 rounded-full bg-red-500 shrink-0"></div>
                <div>
                  <p className="text-sm font-medium text-slate-900">High Velocity Transfer Blocked</p>
                  <p className="text-xs text-slate-500 mt-0.5">$45,000 to flagged offshore account. AI blocked via Policy #POL-104.</p>
                  <p className="text-xs text-slate-400 mt-1">2 mins ago</p>
                </div>
              </div>
              <div className="flex gap-3">
                <div className="w-2 h-2 mt-2 rounded-full bg-amber-500 shrink-0"></div>
                <div>
                  <p className="text-sm font-medium text-slate-900">Unusual Device Login</p>
                  <p className="text-xs text-slate-500 mt-0.5">Attempt from unrecognized IP range. Step-up authentication enforced.</p>
                  <p className="text-xs text-slate-400 mt-1">15 mins ago</p>
                </div>
              </div>
               <div className="flex gap-3">
                <div className="w-2 h-2 mt-2 rounded-full bg-emerald-500 shrink-0"></div>
                <div>
                  <p className="text-sm font-medium text-slate-900">System Integrity Check</p>
                  <p className="text-xs text-slate-500 mt-0.5">Fraud detection models re-calibrated. Accuracy 99.8%.</p>
                  <p className="text-xs text-slate-400 mt-1">1 hr ago</p>
                </div>
              </div>
            </div>
            <button className="w-full mt-6 py-2 bg-slate-50 hover:bg-slate-100 text-slate-700 text-sm font-medium rounded-lg border border-slate-200 transition-colors">
              View All Logs
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
