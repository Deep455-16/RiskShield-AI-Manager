'use client'

import { useEffect, useState } from 'react'
import { Lock, Download } from 'lucide-react'

export default function AuditPage() {
  const [logs, setLogs] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/audit-logs').then(r => r.json()).then(d => setLogs(d.data || [])).finally(() => setLoading(false))
  }, [])

  const exportCSV = () => {
    window.open('/api/v1/audit-logs/export', '_blank')
  }

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-slate-700 p-2 rounded-lg"><Lock className="w-6 h-6 text-white" /></div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Audit Trail</h1>
            <p className="text-slate-500">Immutable log of every governance action and decision</p>
          </div>
        </div>
        <button onClick={exportCSV} className="flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 text-slate-700 font-semibold rounded-lg hover:bg-slate-50 shadow-sm transition-colors">
          <Download className="w-4 h-4" /> Export CSV
        </button>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Timestamp</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Action</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Resource Type</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Resource ID</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Actor</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-slate-500">Loading audit logs...</td></tr>
            ) : logs.length === 0 ? (
              <tr><td colSpan={5} className="px-6 py-12 text-center text-slate-400">No audit events recorded yet.</td></tr>
            ) : (
              (logs as any[]).map((log: any) => (
                <tr key={log.id} className="border-b border-slate-100 hover:bg-slate-50 text-sm">
                  <td className="px-6 py-3 text-slate-500 whitespace-nowrap">{new Date(log.timestamp || log.created_at).toLocaleString()}</td>
                  <td className="px-6 py-3">
                    <span className="px-2 py-0.5 rounded bg-slate-100 text-slate-700 text-xs font-mono font-semibold">{log.action}</span>
                  </td>
                  <td className="px-6 py-3 text-slate-600 capitalize">{(log.resource_type || '').replace(/_/g, ' ')}</td>
                  <td className="px-6 py-3 text-slate-500 font-mono text-xs">{log.resource_id ? String(log.resource_id).slice(0, 12) + '…' : '—'}</td>
                  <td className="px-6 py-3 text-slate-600 font-mono text-xs">{log.actor_id ? String(log.actor_id).slice(0, 12) + '…' : 'system'}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}