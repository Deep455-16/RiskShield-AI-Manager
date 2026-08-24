'use client'
import { useEffect, useState } from 'react'
export default function Systems() {
  const [systems, setSystems] = useState([])
  useEffect(() => {
    fetch('/api/v1/ai-systems').then(r => r.json()).then(d => setSystems(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">AI Systems Registry</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {systems.map((s: any) => (
          <div key={s.id} className="bg-slate-900 rounded-lg border border-slate-800 p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold">{s.name}</h3>
              <span className={`text-xs px-2 py-1 rounded ${s.approval_status==='approved'?'bg-emerald-500/20 text-emerald-400':'bg-yellow-500/20 text-yellow-400'}`}>{s.approval_status}</span>
            </div>
            <p className="text-sm text-slate-400 mb-2">{s.description}</p>
            <div className="text-xs text-slate-500 space-y-1">
              <div>Purpose: {s.purpose}</div>
              <div>Owner: {s.owner}</div>
              <div>Risk Class: <span className={s.risk_class==='high'?'text-red-400':'text-slate-400'}>{s.risk_class}</span></div>
              <div>Deployment: {s.deployment_status}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}