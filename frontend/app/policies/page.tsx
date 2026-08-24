'use client'
import { useEffect, useState } from 'react'
export default function Policies() {
  const [policies, setPolicies] = useState([])
  useEffect(() => {
    fetch('/api/v1/policies').then(r => r.json()).then(d => setPolicies(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Policy Engine</h1>
      <div className="space-y-3">
        {policies.map((p: any) => (
          <div key={p.id} className="bg-slate-900 rounded-lg border border-slate-800 p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold">{p.name}</h3>
              <div className="flex gap-2">
                <span className="text-xs px-2 py-1 rounded bg-slate-800 text-slate-400">v{p.version}</span>
                <span className={`text-xs px-2 py-1 rounded ${p.status==='active'?'bg-emerald-500/20 text-emerald-400':'bg-slate-700 text-slate-400'}`}>{p.status}</span>
              </div>
            </div>
            <p className="text-sm text-slate-400">{p.description}</p>
            <div className="text-xs text-slate-500 mt-2">Action: <span className="font-mono text-slate-300">{p.action}</span> | Priority: {p.priority}</div>
          </div>
        ))}
      </div>
    </div>
  )
}