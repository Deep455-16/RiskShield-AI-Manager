'use client'

import { useEffect, useState } from 'react'
import { Building2, PlusCircle, ShieldCheck, AlertTriangle } from 'lucide-react'

export default function VendorsPage() {
  const [vendors, setVendors] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/vendors').then(r => r.json()).then(d => setVendors(d.data || [])).finally(() => setLoading(false))
  }, [])

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-teal-700 p-2 rounded-lg"><Building2 className="w-6 h-6 text-white" /></div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Vendor Risk Management</h1>
            <p className="text-slate-500">Track AI, data, and platform vendors with risk-tiered assessments</p>
          </div>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 bg-teal-700 text-white font-semibold rounded-lg hover:bg-teal-800 transition-colors">
          <PlusCircle className="w-4 h-4" /> Add Vendor
        </button>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Vendor</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Type</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Risk Tier</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Renewal Date</th>
              <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-slate-500">Loading vendors...</td></tr>
            ) : vendors.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center">
                  <Building2 className="w-8 h-8 text-slate-300 mx-auto mb-2" />
                  <p className="text-slate-500">No vendors registered. Click "Add Vendor" to begin tracking third-party AI providers.</p>
                </td>
              </tr>
            ) : (
              (vendors as any[]).map((v: any) => (
                <tr key={v.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-6 py-4 font-semibold text-slate-900">{v.name}</td>
                  <td className="px-6 py-4 text-sm text-slate-600 capitalize">{(v.type || '').replace(/_/g, ' ')}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold uppercase
                      ${v.risk_tier === 'high' ? 'bg-rose-100 text-rose-700 border border-rose-200'
                      : v.risk_tier === 'low' ? 'bg-emerald-100 text-emerald-700 border border-emerald-200'
                      : 'bg-amber-100 text-amber-700 border border-amber-200'}`}>
                      {v.risk_tier === 'high' ? <AlertTriangle className="w-3 h-3" /> : <ShieldCheck className="w-3 h-3" />}
                      {v.risk_tier}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-slate-600">{v.renewal_date ? new Date(v.renewal_date).toLocaleDateString() : '—'}</td>
                  <td className="px-6 py-4 text-right space-x-3">
                    <button className="text-xs font-semibold text-blue-700 hover:underline">Assess</button>
                    <button className="text-xs font-semibold text-slate-500 hover:underline">Edit</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
