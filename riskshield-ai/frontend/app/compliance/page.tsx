'use client'

import { useEffect, useState } from 'react'
import { CheckCircle, ShieldCheck, AlertCircle, FileText, Download } from 'lucide-react'

export default function Compliance() {
  const [controls, setControls] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/compliance/controls')
      .then(r => r.json())
      .then(d => setControls(d.data || []))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="min-h-screen bg-slate-50 p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="bg-emerald-700 p-2 rounded-lg">
              <ShieldCheck className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Compliance Register</h1>
              <p className="text-slate-500">Statement of Applicability (EU AI Act, NIST AI RMF, ISO 42001)</p>
            </div>
          </div>
        </div>
        <button className="flex items-center gap-2 bg-white border border-slate-200 text-slate-700 px-4 py-2 rounded-lg font-semibold shadow-sm hover:bg-slate-50 transition-colors">
          <Download className="w-4 h-4" /> Export SOA Report
        </button>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200">
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Framework</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Control Ref</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Description</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Status</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold text-right">Evidence</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-slate-500">Loading compliance data...</td>
                </tr>
              ) : controls.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-12 text-center">
                    <FileText className="w-8 h-8 text-slate-300 mx-auto mb-2" />
                    <p className="text-slate-500 font-medium">No compliance controls found.</p>
                  </td>
                </tr>
              ) : (
                controls.map((c: any) => (
                  <tr key={c.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-6 py-4 font-semibold text-slate-700 whitespace-nowrap">
                      {c.framework === 'eu_ai_act' ? 'EU AI Act' : c.framework === 'nist_ai_rmf' ? 'NIST AI RMF' : c.framework}
                    </td>
                    <td className="px-6 py-4 text-sm font-medium text-slate-900 whitespace-nowrap">{c.control_ref}</td>
                    <td className="px-6 py-4 text-sm text-slate-600 max-w-md truncate">{c.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wide
                        ${c.status === 'implemented' ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' 
                        : c.status === 'partial' ? 'bg-amber-100 text-amber-700 border border-amber-200' 
                        : 'bg-rose-100 text-rose-700 border border-rose-200'}`}>
                        {c.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right whitespace-nowrap">
                      {c.evidence_url ? (
                        <a href={c.evidence_url} className="text-sm font-semibold text-blue-600 hover:underline">View Evidence</a>
                      ) : (
                        <span className="text-sm text-slate-400">Missing</span>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}