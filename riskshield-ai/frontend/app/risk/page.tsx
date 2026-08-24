'use client'
import { useEffect, useState } from 'react'

const MOCK_SCORES = [
  { id: '1', asset_id: 'fdtct-001-xgb', asset_type: 'agent',      overall_score: 82, risk_level: 'CRITICAL', timestamp: new Date(Date.now() - 300000).toISOString() },
  { id: '2', asset_id: 'crdt-scr-002',  asset_type: 'model',      overall_score: 67, risk_level: 'HIGH',     timestamp: new Date(Date.now() - 900000).toISOString() },
  { id: '3', asset_id: 'pii-rdc-003',   asset_type: 'agent',      overall_score: 34, risk_level: 'MEDIUM',   timestamp: new Date(Date.now() - 1800000).toISOString() },
  { id: '4', asset_id: 'anml-dtct-004', asset_type: 'system',     overall_score: 71, risk_level: 'HIGH',     timestamp: new Date(Date.now() - 3600000).toISOString() },
  { id: '5', asset_id: 'churn-pr-005',  asset_type: 'model',      overall_score: 18, risk_level: 'LOW',      timestamp: new Date(Date.now() - 7200000).toISOString() },
  { id: '6', asset_id: 'sent-anl-006',  asset_type: 'agent',      overall_score: 22, risk_level: 'LOW',      timestamp: new Date(Date.now() - 10800000).toISOString() },
  { id: '7', asset_id: 'api-gw-007',    asset_type: 'api',        overall_score: 55, risk_level: 'MEDIUM',   timestamp: new Date(Date.now() - 14400000).toISOString() },
  { id: '8', asset_id: 'data-pp-008',   asset_type: 'pipeline',   overall_score: 41, risk_level: 'MEDIUM',   timestamp: new Date(Date.now() - 18000000).toISOString() },
]

const levelColors: Record<string, string> = {
  CRITICAL: 'bg-red-500/20 text-red-400',
  HIGH:     'bg-orange-500/20 text-orange-400',
  MEDIUM:   'bg-yellow-500/20 text-yellow-400',
  LOW:      'bg-emerald-500/20 text-emerald-400',
}

export default function RiskCenter() {
  const [scores, setScores] = useState(MOCK_SCORES)

  useEffect(() => {
    fetch('/api/v1/risk/scores')
      .then(r => r.ok ? r.json() : null)
      .then(d => { if (d?.data?.length) setScores(d.data) })
      .catch(() => {})
  }, [])

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-2">Risk Center</h1>
      <p className="text-slate-400 text-sm mb-6">7-dimension weighted risk scores across all monitored assets</p>
      <div className="bg-slate-900 rounded-xl border border-slate-800 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-800 text-slate-400">
            <tr>
              <th className="px-4 py-3 text-left">Asset ID</th>
              <th className="px-4 py-3 text-left">Type</th>
              <th className="px-4 py-3 text-left">Score</th>
              <th className="px-4 py-3 text-left">Risk Level</th>
              <th className="px-4 py-3 text-left">Assessed</th>
            </tr>
          </thead>
          <tbody>
            {scores.map((s: any) => (
              <tr key={s.id} className="border-t border-slate-800 hover:bg-slate-800/50 transition-colors">
                <td className="px-4 py-3 font-mono text-xs text-slate-300">{s.asset_id?.slice(0, 14)}</td>
                <td className="px-4 py-3 capitalize text-slate-300">{s.asset_type}</td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-slate-100">{s.overall_score}</span>
                    <div className="w-20 bg-slate-700 rounded-full h-1.5">
                      <div
                        className={`h-1.5 rounded-full ${s.overall_score >= 75 ? 'bg-red-500' : s.overall_score >= 50 ? 'bg-orange-500' : s.overall_score >= 25 ? 'bg-yellow-500' : 'bg-emerald-500'}`}
                        style={{ width: `${s.overall_score}%` }}
                      />
                    </div>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${levelColors[s.risk_level] ?? 'bg-slate-700 text-slate-400'}`}>
                    {s.risk_level}
                  </span>
                </td>
                <td className="px-4 py-3 text-slate-500 text-xs">{new Date(s.timestamp).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}