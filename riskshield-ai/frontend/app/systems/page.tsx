'use client'

import { useEffect, useState } from 'react'
import { Server, Shield, Activity, PlusCircle, ServerCog } from 'lucide-react'

export default function Systems() {
  const [systems, setSystems] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/ai-systems')
      .then(r => r.json())
      .then(d => setSystems(d.data || []))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="min-h-screen bg-slate-50 p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="bg-blue-900 p-2 rounded-lg">
              <Server className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-900 tracking-tight">AI System Registry</h1>
              <p className="text-slate-500">Inventory and governance of deployed AI models</p>
            </div>
          </div>
        </div>
        <button className="flex items-center gap-2 bg-blue-700 hover:bg-blue-800 text-white px-4 py-2 rounded-lg font-semibold transition-colors">
          <PlusCircle className="w-4 h-4" /> Register AI System
        </button>
      </div>

      {loading ? (
        <div className="text-center text-slate-500 py-20">Loading registry...</div>
      ) : systems.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-12 text-center shadow-sm">
          <ServerCog className="w-12 h-12 text-slate-300 mx-auto mb-4" />
          <h2 className="text-lg font-bold text-slate-900 mb-2">No AI Systems Registered</h2>
          <p className="text-slate-500 max-w-sm mx-auto mb-6">You haven't added any AI systems to the registry yet. Register your first model to begin tracking governance and compliance.</p>
          <button className="flex items-center gap-2 bg-blue-700 hover:bg-blue-800 text-white px-4 py-2 rounded-lg font-semibold transition-colors mx-auto">
            <PlusCircle className="w-4 h-4" /> Register System
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {systems.map((s: any) => (
            <div key={s.id} className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
              <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50">
                <div className="flex items-center gap-3">
                  <ServerCog className="w-5 h-5 text-blue-700" />
                  <h3 className="font-bold text-slate-900">{s.name}</h3>
                </div>
                <span className={`text-xs px-2 py-1 rounded-full font-bold uppercase tracking-wide
                  ${s.risk_tier === 'high' ? 'bg-rose-100 text-rose-700 border border-rose-200' 
                  : s.risk_tier === 'minimal' ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' 
                  : 'bg-amber-100 text-amber-700 border border-amber-200'}`}>
                  {s.risk_tier || 'medium'} Risk
                </span>
              </div>
              <div className="p-6">
                <p className="text-sm text-slate-600 mb-6">{s.description}</p>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-xs text-slate-400 uppercase font-semibold mb-1">Version</p>
                    <p className="text-sm font-medium text-slate-900">{s.version || '1.0.0'}</p>
                  </div>
                  <div>
                    <p className="text-xs text-slate-400 uppercase font-semibold mb-1">Deployment Env</p>
                    <p className="text-sm font-medium text-slate-900 capitalize">{s.deployment_env || 'production'}</p>
                  </div>
                  <div className="col-span-2">
                    <p className="text-xs text-slate-400 uppercase font-semibold mb-1">Purpose</p>
                    <p className="text-sm font-medium text-slate-900">{s.purpose}</p>
                  </div>
                </div>
              </div>
              <div className="px-6 py-3 bg-slate-50 border-t border-slate-100 flex gap-4">
                <button className="text-sm font-semibold text-blue-700 hover:text-blue-800">View Model Card</button>
                <button className="text-sm font-semibold text-slate-600 hover:text-slate-800">Manage Compliance</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}