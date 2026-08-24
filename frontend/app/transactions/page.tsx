'use client'

import { useEffect, useState } from 'react'
import { AlertCircle, CheckCircle, Search, Filter, ArrowUpRight, ArrowDownRight, ShieldAlert } from 'lucide-react'

export default function Transactions() {
  const [txs, setTxs] = useState([])

  useEffect(() => {
    fetch('/api/v1/transactions')
      .then(r => r.json())
      .then(d => setTxs(d.data || []))
      .catch(() => {})
  }, [])

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 p-8 pb-20">
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Transaction Fraud Monitoring</h1>
          <p className="text-slate-500 mt-1">Real-time AI analysis of banking ledger movements</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input 
              type="text" 
              placeholder="Search Txn ID..." 
              className="pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:outline-none focus:border-blue-500 w-64 shadow-sm"
            />
          </div>
          <button className="px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-medium hover:bg-slate-50 flex items-center gap-2 shadow-sm">
            <Filter className="w-4 h-4" /> Filter
          </button>
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm flex justify-between items-center">
          <div>
            <p className="text-sm font-medium text-slate-500 uppercase tracking-wide">Processed Today</p>
            <p className="text-3xl font-bold text-slate-900 mt-1">12,450</p>
          </div>
          <div className="w-12 h-12 bg-blue-50 text-blue-600 rounded-full flex items-center justify-center">
            <ArrowUpRight className="w-6 h-6" />
          </div>
        </div>
        <div className="bg-white p-6 rounded-xl border border-rose-100 shadow-sm flex justify-between items-center relative overflow-hidden">
          <div className="absolute top-0 right-0 w-2 h-full bg-rose-500"></div>
          <div>
            <p className="text-sm font-medium text-slate-500 uppercase tracking-wide">AI Blocked</p>
            <p className="text-3xl font-bold text-slate-900 mt-1">42</p>
          </div>
          <div className="w-12 h-12 bg-rose-50 text-rose-600 rounded-full flex items-center justify-center">
            <ShieldAlert className="w-6 h-6" />
          </div>
        </div>
        <div className="bg-white p-6 rounded-xl border border-emerald-100 shadow-sm flex justify-between items-center relative overflow-hidden">
          <div className="absolute top-0 right-0 w-2 h-full bg-emerald-500"></div>
          <div>
            <p className="text-sm font-medium text-slate-500 uppercase tracking-wide">Clearance Rate</p>
            <p className="text-3xl font-bold text-slate-900 mt-1">99.6%</p>
          </div>
          <div className="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-full flex items-center justify-center">
            <CheckCircle className="w-6 h-6" />
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100 bg-slate-50">
          <h2 className="font-semibold text-slate-900">Live Transaction Feed</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-left">
            <thead className="bg-white border-b border-slate-100 text-slate-500 uppercase text-xs font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Transaction ID</th>
                <th className="px-6 py-4">Amount</th>
                <th className="px-6 py-4">Risk Score</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">AI Guardrail Insight</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {txs.map((t: any) => (
                <tr key={t.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 font-mono text-slate-600 text-xs">{t.id.split('-')[0]}...</td>
                  <td className="px-6 py-4">
                    <span className={`font-semibold ${t.amount > 10000 ? 'text-rose-600' : 'text-slate-900'}`}>
                      {new Intl.NumberFormat('en-US', { style: 'currency', currency: t.currency || 'USD' }).format(t.amount)}
                    </span>
                    {t.amount > 10000 && <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-rose-100 text-rose-700">HIGH VALUE</span>}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <div className="w-16 h-2 bg-slate-100 rounded-full overflow-hidden">
                        <div 
                          className={`h-full ${t.risk_score > 75 ? 'bg-rose-500' : t.risk_score > 30 ? 'bg-amber-500' : 'bg-emerald-500'}`}
                          style={{ width: `${Math.min(100, t.risk_score || 0)}%` }}
                        ></div>
                      </div>
                      <span className="font-medium text-slate-700">{Math.round(t.risk_score || 0)}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold border
                      ${t.status === 'blocked' ? 'bg-rose-50 text-rose-700 border-rose-200' : 
                        t.status === 'approved' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 
                        'bg-amber-50 text-amber-700 border-amber-200'}`
                    }>
                      {t.status === 'blocked' && <AlertCircle className="w-3 h-3" />}
                      {t.status === 'approved' && <CheckCircle className="w-3 h-3" />}
                      {t.status?.toUpperCase() || 'UNKNOWN'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-500 text-xs">
                    {t.status === 'blocked' ? (
                      <span className="text-rose-600 font-medium flex items-center gap-1">
                         <ShieldAlert className="w-3.5 h-3.5" /> Blocked: Velocity & IP Novelty Detected
                      </span>
                    ) : t.amount > 10000 ? (
                      "Cleared after secondary model review"
                    ) : (
                      "Standard low-risk clearance"
                    )}
                  </td>
                </tr>
              ))}
              {txs.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-6 py-12 text-center text-slate-500">
                    No transactions found. Use the Simulator to generate some.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}