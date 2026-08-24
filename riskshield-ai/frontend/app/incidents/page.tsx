'use client'
import { useEffect, useState } from 'react'
export default function Incidents() {
  const [incidents, setIncidents] = useState([])
  useEffect(() => {
    fetch('/api/v1/incidents').then(r => r.json()).then(d => setIncidents(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Incident Management</h1>
      <div className="space-y-3">
        {incidents.map((i: any) => (
          <div key={i.id} className="bg-slate-900 rounded-lg border border-slate-800 p-4 flex items-center justify-between">
            <div>
              <h3 className="font-semibold">{i.title}</h3>
              <p className="text-sm text-slate-400">{i.description}</p>
            </div>
            <div className="text-right">
              <span className={`px-2 py-1 rounded text-xs ${i.status==='OPEN'?'bg-red-500/20 text-red-400':i.status==='RESOLVED'?'bg-emerald-500/20 text-emerald-400':'bg-yellow-500/20 text-yellow-400'}`}>{i.status}</span>
              <div className="text-xs text-slate-500 mt-1">{i.severity}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}