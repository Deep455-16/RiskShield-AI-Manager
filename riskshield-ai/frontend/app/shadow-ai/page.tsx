'use client'

import { useEffect, useState } from 'react'
import { Eye, RefreshCw, Github, Cloud, CheckCircle, AlertTriangle } from 'lucide-react'

export default function ShadowAIPage() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [scanning, setScanning] = useState(false)

  const fetchItems = () => {
    setLoading(true)
    fetch('/api/v1/shadow-ai').then(r => r.json()).then(d => setItems(d.data || [])).finally(() => setLoading(false))
  }

  useEffect(() => { fetchItems() }, [])

  const runScan = async () => {
    setScanning(true)
    await fetch('/api/v1/shadow-ai/discover', { method: 'POST' })
    fetchItems()
    setScanning(false)
  }

  const sourceIcon = (source: string) => {
    if (source === 'github') return <Github className="w-4 h-4" />
    return <Cloud className="w-4 h-4" />
  }

  const sourceColor: Record<string, string> = {
    github: 'bg-slate-100 text-slate-700',
    aws_bedrock: 'bg-orange-100 text-orange-700',
    azure_ai: 'bg-blue-100 text-blue-700',
    huggingface: 'bg-yellow-100 text-yellow-700',
    vertex_ai: 'bg-green-100 text-green-700',
  }

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-slate-800 p-2 rounded-lg"><Eye className="w-6 h-6 text-white" /></div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Shadow AI Discovery</h1>
            <p className="text-slate-500">Detect unregistered AI models across cloud providers and code repositories</p>
          </div>
        </div>
        <button
          onClick={runScan}
          disabled={scanning}
          className="flex items-center gap-2 px-4 py-2 bg-slate-800 text-white font-semibold rounded-lg hover:bg-slate-700 disabled:opacity-50 transition-colors"
        >
          <RefreshCw className={`w-4 h-4 ${scanning ? 'animate-spin' : ''}`} />
          {scanning ? 'Scanning…' : 'Run Discovery Scan'}
        </button>
      </div>

      <div className="mb-6 flex gap-4">
        <div className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm flex-1 text-center">
          <p className="text-3xl font-bold text-slate-900">{(items as any[]).length}</p>
          <p className="text-sm text-slate-500">Total Discovered</p>
        </div>
        <div className="bg-white border border-amber-200 rounded-xl p-4 shadow-sm flex-1 text-center">
          <p className="text-3xl font-bold text-amber-700">{(items as any[]).filter((i: any) => i.status === 'unreviewed').length}</p>
          <p className="text-sm text-slate-500">Unreviewed</p>
        </div>
        <div className="bg-white border border-emerald-200 rounded-xl p-4 shadow-sm flex-1 text-center">
          <p className="text-3xl font-bold text-emerald-700">{(items as any[]).filter((i: any) => i.status === 'approved').length}</p>
          <p className="text-sm text-slate-500">Approved</p>
        </div>
      </div>

      {loading ? (
        <div className="text-center text-slate-500 py-20">Loading shadow AI inbox...</div>
      ) : items.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-12 text-center shadow-sm">
          <Eye className="w-12 h-12 text-slate-300 mx-auto mb-4" />
          <h2 className="text-lg font-bold text-slate-900 mb-2">Inbox is Empty</h2>
          <p className="text-slate-500 mb-6">Click "Run Discovery Scan" to scan your cloud providers and code repositories for unregistered AI models.</p>
          <button onClick={runScan} className="px-4 py-2 bg-slate-800 text-white font-semibold rounded-lg hover:bg-slate-700">Start Scan</button>
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200">
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Source</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Name</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">External ID</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Discovered</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Status</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold text-right">Action</th>
              </tr>
            </thead>
            <tbody>
              {(items as any[]).map((item: any) => (
                <tr key={item.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${sourceColor[item.source] || 'bg-slate-100 text-slate-700'}`}>
                      {sourceIcon(item.source)} {item.source}
                    </span>
                  </td>
                  <td className="px-6 py-4 font-semibold text-slate-900">{item.name}</td>
                  <td className="px-6 py-4 text-sm text-slate-600 font-mono">{item.external_id}</td>
                  <td className="px-6 py-4 text-sm text-slate-500">{new Date(item.discovered_at).toLocaleDateString()}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold uppercase
                      ${item.status === 'unreviewed' ? 'bg-amber-100 text-amber-700 border border-amber-200'
                      : item.status === 'approved' ? 'bg-emerald-100 text-emerald-700 border border-emerald-200'
                      : 'bg-slate-100 text-slate-600 border border-slate-200'}`}>
                      {item.status === 'unreviewed' ? <AlertTriangle className="w-3 h-3" /> : <CheckCircle className="w-3 h-3" />}
                      {item.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right space-x-2">
                    <button className="text-xs font-semibold text-blue-700 hover:underline">Register</button>
                    <button className="text-xs font-semibold text-slate-500 hover:underline">Ignore</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
