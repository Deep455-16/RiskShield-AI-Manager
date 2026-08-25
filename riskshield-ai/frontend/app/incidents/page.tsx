'use client'

import { useEffect, useState } from 'react'
import { AlertTriangle, Clock, CheckCircle, ShieldAlert, Activity } from 'lucide-react'

export default function Incidents() {
  const [incidents, setIncidents] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/incidents')
      .then(r => r.json())
      .then(d => setIncidents(d.data || []))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="min-h-screen bg-slate-50 p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="bg-rose-700 p-2 rounded-lg">
              <ShieldAlert className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Incident & CAPA Management</h1>
              <p className="text-slate-500">Track anomalies, SLA deadlines, and Corrective/Preventive Actions</p>
            </div>
          </div>
        </div>
      </div>

      {loading ? (
        <div className="text-center text-slate-500 py-20">Loading incidents...</div>
      ) : incidents.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-12 text-center shadow-sm">
          <CheckCircle className="w-12 h-12 text-emerald-400 mx-auto mb-4" />
          <h2 className="text-lg font-bold text-slate-900 mb-2">No Active Incidents</h2>
          <p className="text-slate-500">Your AI systems are running smoothly. Any risk engine blocks or guardrail anomalies will appear here.</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {incidents.map((i: any) => (
            <div key={i.id} className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden flex flex-col md:flex-row">
              <div className={`w-2 ${i.status === 'OPEN' ? 'bg-rose-500' : i.status === 'RESOLVED' ? 'bg-emerald-500' : 'bg-amber-500'}`} />
              <div className="p-6 flex-1 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-lg font-bold text-slate-900">{i.title}</h3>
                    <span className={`px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wide
                      ${i.severity === 'CRITICAL' || i.severity === 'HIGH' ? 'bg-rose-100 text-rose-700' : 'bg-amber-100 text-amber-700'}`}>
                      {i.severity}
                    </span>
                  </div>
                  <p className="text-sm text-slate-600 mb-4">{i.description}</p>
                  
                  <div className="flex items-center gap-4 text-xs font-semibold">
                    <span className="flex items-center gap-1 text-slate-500"><Activity className="w-4 h-4" /> Category: {i.category || 'System'}</span>
                    {i.sla_deadline && (
                      <span className="flex items-center gap-1 text-rose-600 bg-rose-50 px-2 py-1 rounded"><Clock className="w-4 h-4" /> SLA: {new Date(i.sla_deadline).toLocaleString()}</span>
                    )}
                  </div>
                </div>
                <div className="text-right flex flex-col items-end gap-3">
                  <span className={`px-3 py-1.5 rounded-lg text-sm font-bold uppercase tracking-wide border
                    ${i.status === 'OPEN' ? 'bg-rose-50 text-rose-700 border-rose-200' : i.status === 'RESOLVED' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-amber-50 text-amber-700 border-amber-200'}`}>
                    {i.status}
                  </span>
                  <button className="text-sm font-semibold text-blue-700 hover:text-blue-800 bg-blue-50 px-3 py-1.5 rounded-lg border border-blue-200">
                    Manage CAPA
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}