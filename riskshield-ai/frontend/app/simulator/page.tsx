'use client'

import { useState } from 'react'
import { AlertTriangle, Shield, Play, ArrowRight } from 'lucide-react'

export default function Simulator() {
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)

  const runSimulation = async (type: string) => {
    setLoading(true)
    const res = await fetch('/api/v1/simulate/attack', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ attack_type: type }),
    })
    const data = await res.json()
    setResult(data)
    setLoading(false)
  }

  const runPaymentSim = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    const form = new FormData(e.currentTarget)
    const body = {
      amount: parseFloat(form.get('amount') as string),
      currency: form.get('currency'),
      velocity: parseInt(form.get('velocity') as string),
      device_novelty: form.get('device_novelty') === 'on',
      location_novelty: form.get('location_novelty') === 'on',
      merchant_risk: parseFloat(form.get('merchant_risk') as string),
      ip_reputation: parseFloat(form.get('ip_reputation') as string),
      model_confidence: parseFloat(form.get('model_confidence') as string),
    }
    const res = await fetch('/api/v1/simulate/payment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const data = await res.json()
    setResult(data)
    setLoading(false)
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-2">Risk Simulator</h1>
      <p className="text-amber-400 text-sm font-medium mb-8">⚠️ SIMULATION — NO REAL PAYMENT PROCESSING</p>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-slate-900 rounded-lg border border-slate-800 p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <Shield className="w-5 h-5 text-emerald-500" />
            Payment Risk Simulator
          </h2>
          <form onSubmit={runPaymentSim} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm text-slate-400 mb-1">Amount</label>
                <input name="amount" type="number" defaultValue="87500" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" />
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1">Currency</label>
                <select name="currency" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm">
                  <option value="INR">INR</option>
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm text-slate-400 mb-1">Velocity (txns/hour)</label>
              <input name="velocity" type="number" defaultValue="5" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm text-slate-400 mb-1">Merchant Risk (0-1)</label>
                <input name="merchant_risk" type="number" step="0.1" max="1" min="0" defaultValue="0.6" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" />
              </div>
              <div>
                <label className="block text-sm text-slate-400 mb-1">IP Reputation (0-1)</label>
                <input name="ip_reputation" type="number" step="0.1" max="1" min="0" defaultValue="0.4" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" />
              </div>
            </div>
            <div>
              <label className="block text-sm text-slate-400 mb-1">Model Confidence (0-1)</label>
              <input name="model_confidence" type="number" step="0.01" max="1" min="0" defaultValue="0.81" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" />
            </div>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 text-sm">
                <input name="device_novelty" type="checkbox" className="rounded bg-slate-800 border-slate-700" />
                New Device
              </label>
              <label className="flex items-center gap-2 text-sm">
                <input name="location_novelty" type="checkbox" className="rounded bg-slate-800 border-slate-700" />
                New Location
              </label>
            </div>
            <button type="submit" disabled={loading} className="w-full bg-emerald-600 hover:bg-emerald-500 text-white font-medium py-2 rounded-lg transition-colors flex items-center justify-center gap-2">
              <Play size={16} /> {loading ? 'Running...' : 'Run Risk Assessment'}
            </button>
          </form>
        </div>

        <div className="bg-slate-900 rounded-lg border border-slate-800 p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <AlertTriangle className="w-5 h-5 text-red-500" />
            Attack Simulator
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {['prompt_injection','pii_leakage','suspicious_transaction','privilege_escalation','model_drift','fairness_disparity','api_abuse','excessive_spending'].map((type) => (
              <button key={type} onClick={() => runSimulation(type)} disabled={loading}
                className="p-3 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-lg text-sm text-left transition-colors capitalize">
                {type.replace(/_/g, ' ')}
              </button>
            ))}
          </div>
        </div>
      </div>

      {result && (
        <div className="mt-8 bg-slate-900 rounded-lg border border-slate-800 p-6">
          <h3 className="text-lg font-semibold mb-4">Simulation Results</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div className="bg-slate-800 rounded-lg p-4">
              <div className="text-sm text-slate-400">Risk Score</div>
              <div className={`text-3xl font-bold ${getScoreColor(result.risk_score?.overall_score || result.after?.score || 0)}`}>
                {Math.round(result.risk_score?.overall_score || result.after?.score || 0)}
              </div>
            </div>
            <div className="bg-slate-800 rounded-lg p-4">
              <div className="text-sm text-slate-400">Risk Level</div>
              <div className="text-3xl font-bold text-slate-200">{result.risk_score?.risk_level || result.after?.level || 'UNKNOWN'}</div>
            </div>
            <div className="bg-slate-800 rounded-lg p-4">
              <div className="text-sm text-slate-400">Decision</div>
              <div className={`text-3xl font-bold ${getDecisionColor(result.policy_decision?.decision || 'ALLOW')}`}>
                {result.policy_decision?.decision || 'ALLOW'}
              </div>
            </div>
          </div>
          {result.before && result.after && (
            <div className="flex items-center gap-4 bg-slate-800 rounded-lg p-4">
              <div className="text-center"><div className="text-sm text-slate-400">Before</div><div className="text-xl font-bold text-emerald-400">{result.before.score} {result.before.level}</div></div>
              <ArrowRight className="text-slate-500" />
              <div className="text-center"><div className="text-sm text-slate-400">After Attack</div><div className="text-xl font-bold text-red-400">{result.after.score} {result.after.level}</div></div>
              <div className="flex-1 text-right"><div className="text-sm text-slate-400">Incident Created</div><div className="text-sm font-mono text-slate-300">{result.incident_created}</div></div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function getScoreColor(score: number): string {
  if (score >= 75) return 'text-red-500'
  if (score >= 50) return 'text-orange-500'
  if (score >= 25) return 'text-yellow-500'
  return 'text-emerald-500'
}
function getDecisionColor(d: string): string {
  if (d === 'BLOCK') return 'text-red-500'
  if (d === 'REVIEW') return 'text-orange-500'
  return 'text-emerald-500'
}
