'use client'
import { useEffect, useState } from 'react'
export default function Agents() {
  const [agents, setAgents] = useState([])
  useEffect(() => {
    fetch('/api/v1/agents').then(r => r.json()).then(d => setAgents(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">AI Agent Registry</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {agents.map((a: any) => (
          <div key={a.id} className="bg-slate-900 rounded-lg border border-slate-800 p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold">{a.name}</h3>
              <span className={`text-xs px-2 py-1 rounded ${a.approval_status==='approved'?'bg-emerald-500/20 text-emerald-400':'bg-yellow-500/20 text-yellow-400'}`}>{a.approval_status}</span>
            </div>
            <p className="text-sm text-slate-400 mb-2">{a.purpose}</p>
            <div className="text-xs text-slate-500 space-y-1">
              <div>Model: {a.model}</div>
              <div>Environment: {a.environment}</div>
              <div>Risk Level: <span className={a.risk_level==='high'?'text-red-400':'text-slate-400'}>{a.risk_level}</span></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}