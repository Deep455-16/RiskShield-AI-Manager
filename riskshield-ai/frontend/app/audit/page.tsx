'use client'
import { useEffect, useState } from 'react'
export default function Audit() {
  const [logs, setLogs] = useState([])
  useEffect(() => {
    fetch('/api/v1/audit-logs').then(r => r.json()).then(d => setLogs(d.data || []))
  }, [])
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Audit Logs</h1>
      <p className="text-slate-400 text-sm mb-4">Tamper-evident SHA-256 chained audit trail</p>
      <div className="bg-slate-900 rounded-lg border border-slate-800 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-800 text-slate-400"><tr><th className="px-4 py-3 text-left">Action</th><th className="px-4 py-3 text-left">Resource</th><th className="px-4 py-3 text-left">Hash</th><th className="px-4 py-3 text-left">Time</th></tr></thead>
          <tbody>
            {logs.map((l: any) => (
              <tr key={l.id} className="border-t border-slate-800">
                <td className="px-4 py-3">{l.action}</td>
                <td className="px-4 py-3 text-slate-400">{l.resource_type}</td>
                <td className="px-4 py-3 font-mono text-xs text-slate-500">{l.hash?.slice(0,16)}...</td>
                <td className="px-4 py-3 text-slate-500 text-xs">{new Date(l.timestamp).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}