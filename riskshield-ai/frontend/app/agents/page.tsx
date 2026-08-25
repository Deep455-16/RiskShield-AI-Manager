'use client'

import { useEffect, useState } from 'react'
import { Bot, Zap, AlertTriangle, Search, XCircle, CheckCircle, Activity } from 'lucide-react'

export default function AgentsPage() {
  const [agents, setAgents] = useState<any[]>([])
  const [logs, setLogs] = useState<Record<string, any[]>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/agents').then(r => r.json()).then(d => setAgents(d.data || [])).finally(() => setLoading(false))
  }, [])

  const killSwitch = async (id: string) => {
    if (!confirm('Are you sure you want to suspend this agent? This will immediately stop all its operations.')) return
    await fetch(`/api/v1/agents/${id}/kill-switch`, { method: 'POST' })
    setAgents((prev: any[]) => prev.map((a: any) => a.id === id ? { ...a, status: 'suspended' } : a))
  }

  const loadLogs = async (id: string) => {
    const r = await fetch(`/api/v1/agents/${id}/behavior-logs`)
    const d = await r.json()
    setLogs(prev => ({ ...prev, [id]: d.data || [] }))
  }

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-indigo-700 p-2 rounded-lg"><Bot className="w-6 h-6 text-white" /></div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Agent Governance</h1>
            <p className="text-slate-500">Monitor, audit, and control all AI agents in real time</p>
          </div>
        </div>
      </div>

      {loading ? <div className="text-center text-slate-500 py-20">Loading agents...</div> : (
        <div className="space-y-6">
          {(agents as any[]).map((agent: any) => (
            <div key={agent.id} className={`bg-white rounded-xl border shadow-sm overflow-hidden ${agent.status === 'suspended' ? 'border-rose-200 opacity-70' : 'border-slate-200'}`}>
              <div className="px-6 py-4 border-b border-slate-100 bg-slate-50 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Bot className="w-5 h-5 text-indigo-600" />
                  <div>
                    <h3 className="font-bold text-slate-900">{agent.name}</h3>
                    <p className="text-xs text-slate-500">{agent.purpose} · v{agent.version || '1.0'} · {agent.environment}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wide
                    ${agent.status === 'suspended' ? 'bg-rose-100 text-rose-700 border border-rose-200'
                    : agent.risk_level === 'high' ? 'bg-amber-100 text-amber-700 border border-amber-200'
                    : 'bg-emerald-100 text-emerald-700 border border-emerald-200'}`}>
                    {agent.status === 'suspended' ? 'SUSPENDED' : agent.risk_level || 'active'} Risk
                  </span>
                  {agent.status !== 'suspended' && (
                    <button
                      onClick={() => killSwitch(agent.id)}
                      className="flex items-center gap-1 px-3 py-1.5 bg-rose-600 text-white text-xs font-bold rounded-lg hover:bg-rose-700 transition-colors"
                    >
                      <XCircle className="w-3.5 h-3.5" /> Kill Switch
                    </button>
                  )}
                  {agent.status === 'suspended' && (
                    <span className="flex items-center gap-1 px-3 py-1.5 bg-rose-50 text-rose-700 text-xs font-bold rounded-lg border border-rose-200">
                      <XCircle className="w-3.5 h-3.5" /> Suspended
                    </span>
                  )}
                </div>
              </div>

              <div className="p-6 grid grid-cols-4 gap-4 mb-4">
                <div><p className="text-xs text-slate-400 uppercase font-semibold mb-1">Model</p><p className="text-sm font-medium text-slate-900">{agent.model || 'N/A'}</p></div>
                <div><p className="text-xs text-slate-400 uppercase font-semibold mb-1">Approval</p>
                  <span className={`text-xs font-bold ${agent.approval_status === 'approved' ? 'text-emerald-700' : 'text-amber-700'}`}>{agent.approval_status}</span>
                </div>
                <div><p className="text-xs text-slate-400 uppercase font-semibold mb-1">Tools</p><p className="text-sm font-medium text-slate-900">{(agent.tools || []).length} configured</p></div>
                <div>
                  <button onClick={() => loadLogs(agent.id)} className="text-xs font-semibold text-indigo-600 hover:underline flex items-center gap-1 mt-2">
                    <Activity className="w-3.5 h-3.5" /> View Behavior Logs
                  </button>
                </div>
              </div>

              {logs[agent.id] && (
                <div className="mx-6 mb-6 bg-slate-50 rounded-lg border border-slate-200 p-4">
                  <h4 className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3">Recent Behavior Logs</h4>
                  {logs[agent.id].length === 0 ? (
                    <p className="text-sm text-slate-500">No behavior logs recorded yet.</p>
                  ) : (
                    <div className="space-y-2">
                      {logs[agent.id].slice(0, 5).map((log: any, i: number) => (
                        <div key={i} className="flex items-center gap-3 text-sm">
                          <span className={`px-2 py-0.5 rounded text-xs font-bold
                            ${log.outcome === 'ok' ? 'bg-emerald-100 text-emerald-700' : log.outcome === 'blocked' ? 'bg-rose-100 text-rose-700' : 'bg-amber-100 text-amber-700'}`}>
                            {log.outcome.toUpperCase()}
                          </span>
                          <span className="text-slate-600">{new Date(log.ts).toLocaleString()}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}