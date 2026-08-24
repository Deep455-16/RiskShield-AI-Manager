'use client'
import { useEffect, useState } from 'react'
export default function Compliance() {
  const [controls, setControls] = useState([])
  useEffect(() => {
    fetch('/api/v1/compliance/controls').then(r => r.json()).then(d => setControls(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Compliance Center</h1>
      <div className="space-y-2">
        {controls.map((c: any) => (
          <div key={c.id} className="bg-slate-900 rounded-lg border border-slate-800 p-4 flex items-center justify-between">
            <div>
              <div className="font-semibold">{c.control_id} — {c.title}</div>
              <div className="text-sm text-slate-400">{c.description}</div>
            </div>
            <span className={`px-2 py-1 rounded text-xs ${c.status==='IMPLEMENTED'?'bg-emerald-500/20 text-emerald-400':c.status==='PARTIALLY_IMPLEMENTED'?'bg-yellow-500/20 text-yellow-400':'bg-red-500/20 text-red-400'}`}>{c.status}</span>
          </div>
        ))}
      </div>
    </div>
  )
}