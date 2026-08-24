'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import {
  Database, Play, CheckCircle, AlertTriangle, XCircle,
  Eye, BarChart2, RefreshCw, Shield
} from 'lucide-react'

interface DatasetMeta {
  id: string
  name: string
  source: string
  status: string
  row_count: number
  column_count: number
  quality_score: number
  has_fraud_labels: boolean
  last_scanned_at: string
}

function StatusBadge({ status }: { status: string }) {
  if (status === 'AVAILABLE') return (
    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
      <CheckCircle className="w-3 h-3" /> AVAILABLE
    </span>
  )
  return (
    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold bg-rose-50 text-rose-700 border border-rose-200">
      <XCircle className="w-3 h-3" /> NOT AVAILABLE
    </span>
  )
}

function QualityBar({ score }: { score: number }) {
  const color = score >= 90 ? 'bg-emerald-500' : score >= 70 ? 'bg-amber-500' : 'bg-rose-500'
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-2 bg-slate-100 rounded-full overflow-hidden">
        <div className={`h-full ${color} transition-all duration-500`} style={{ width: `${score}%` }} />
      </div>
      <span className="text-sm font-semibold text-slate-700 w-12 text-right">{score.toFixed(0)}/100</span>
    </div>
  )
}

export default function DataHubPage() {
  const [datasets, setDatasets] = useState<DatasetMeta[]>([])
  const [loading, setLoading] = useState(true)
  const [replaying, setReplaying] = useState<Record<string, boolean>>({})

  useEffect(() => {
    fetch('/api/v1/datasets')
      .then(r => r.json())
      .then(d => setDatasets(d.data || []))
      .finally(() => setLoading(false))
  }, [])

  const startReplay = async (id: string, speed = 5) => {
    setReplaying(r => ({ ...r, [id]: true }))
    await fetch(`/api/v1/datasets/${id}/replay`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'start', speed })
    })
  }

  const stopReplay = async (id: string) => {
    await fetch(`/api/v1/datasets/${id}/replay`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'stop' })
    })
    setReplaying(r => ({ ...r, [id]: false }))
  }

  const datasetDescriptions: Record<string, { desc: string; useCase: string; purpose: string }> = {
    sfindset: {
      desc: 'SFinDSet for Systematic Detection of FinCrimes',
      useCase: 'Financial crime analysis, suspicious transaction detection, fraud analytics',
      purpose: 'Transaction Fraud Stream',
    },
    global_bank: {
      desc: 'Synthetic Global Bank Transactions Dataset',
      useCase: 'Transaction simulation, customer behavior, fraud analysis, real-time monitoring',
      purpose: 'Transaction Replay Engine',
    },
    bank_marketing: {
      desc: 'UCI Bank Marketing Dataset — Ethical AI & Fairness Evaluation',
      useCase: 'Fairness metrics, model evaluation, demographic parity, ethical AI',
      purpose: 'Ethical AI Module',
    },
  }

  return (
    <div className="min-h-screen bg-slate-50 p-8 pb-20">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <div className="bg-blue-900 p-2 rounded-lg">
            <Database className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Data Hub</h1>
            <p className="text-slate-500">Local banking datasets powering the RiskShield AI engine</p>
          </div>
        </div>
        <div className="mt-4 flex items-center gap-2 px-4 py-2 bg-blue-50 border border-blue-200 rounded-lg text-sm text-blue-700 w-fit">
          <Shield className="w-4 h-4" />
          All data is processed locally. No external API calls or Kaggle credentials required.
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-6 mb-10">
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm text-center">
          <p className="text-3xl font-bold text-slate-900">{datasets.filter(d => d.status === 'AVAILABLE').length}</p>
          <p className="text-sm text-slate-500 mt-1">Datasets Available</p>
        </div>
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm text-center">
          <p className="text-3xl font-bold text-slate-900">
            {datasets.filter(d => d.status === 'AVAILABLE').reduce((acc, d) => acc + d.row_count, 0).toLocaleString()}
          </p>
          <p className="text-sm text-slate-500 mt-1">Total Records Available</p>
        </div>
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm text-center">
          <p className="text-3xl font-bold text-slate-900">
            {datasets.filter(d => d.has_fraud_labels).length}
          </p>
          <p className="text-sm text-slate-500 mt-1">Datasets with Fraud Labels</p>
        </div>
      </div>

      {/* Dataset Cards */}
      {loading ? (
        <div className="text-center text-slate-500 py-20">Scanning local datasets...</div>
      ) : (
        <div className="space-y-6">
          {datasets.map(ds => {
            const info = datasetDescriptions[ds.id] || { desc: ds.source, useCase: 'General analysis', purpose: 'Analysis' }
            const isBankMarketing = ds.id === 'bank_marketing'
            const isReplaying = replaying[ds.id]

            return (
              <div key={ds.id} className={`bg-white rounded-xl border shadow-sm overflow-hidden ${ds.status !== 'AVAILABLE' ? 'opacity-60' : ''}`}>
                <div className="px-6 py-5 border-b border-slate-100 bg-slate-50 flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Database className="w-5 h-5 text-blue-600" />
                    <div>
                      <h2 className="text-lg font-bold text-slate-900">{ds.name}</h2>
                      <p className="text-sm text-slate-500">{info.desc}</p>
                    </div>
                  </div>
                  <StatusBadge status={ds.status} />
                </div>

                <div className="p-6">
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mb-6">
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wide font-medium mb-1">Rows</p>
                      <p className="text-xl font-bold text-slate-900">{ds.row_count.toLocaleString()}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wide font-medium mb-1">Columns</p>
                      <p className="text-xl font-bold text-slate-900">{ds.column_count}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wide font-medium mb-1">Fraud Labels</p>
                      <p className="text-xl font-bold text-slate-900">{ds.has_fraud_labels ? 'Yes' : 'No'}</p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wide font-medium mb-1">Purpose</p>
                      <p className="text-sm font-semibold text-blue-700">{info.purpose}</p>
                    </div>
                  </div>

                  <div className="mb-4">
                    <p className="text-xs text-slate-400 uppercase tracking-wide font-medium mb-2">Data Quality Score</p>
                    <QualityBar score={ds.quality_score} />
                  </div>

                  <div className="mb-6 p-3 bg-slate-50 rounded-lg border border-slate-100">
                    <p className="text-xs text-slate-500"><span className="font-semibold text-slate-700">Use Cases:</span> {info.useCase}</p>
                  </div>

                  <div className="flex flex-wrap gap-3">
                    {!isBankMarketing ? (
                      <>
                        {!isReplaying ? (
                          <button
                            onClick={() => startReplay(ds.id, 5)}
                            disabled={ds.status !== 'AVAILABLE'}
                            className="flex items-center gap-2 px-4 py-2 bg-blue-900 text-white text-sm font-semibold rounded-lg hover:bg-blue-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                          >
                            <Play className="w-4 h-4" /> Start Replay (5x Speed)
                          </button>
                        ) : (
                          <button
                            onClick={() => stopReplay(ds.id)}
                            className="flex items-center gap-2 px-4 py-2 bg-rose-600 text-white text-sm font-semibold rounded-lg hover:bg-rose-700 transition-colors"
                          >
                            <XCircle className="w-4 h-4" /> Stop Replay
                          </button>
                        )}
                        <Link
                          href="/transactions"
                          className="flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 text-slate-700 text-sm font-semibold rounded-lg hover:bg-slate-50 transition-colors shadow-sm"
                        >
                          <Eye className="w-4 h-4" /> View Transaction Stream
                        </Link>
                      </>
                    ) : (
                      <Link
                        href="/ethical-ai"
                        className="flex items-center gap-2 px-4 py-2 bg-purple-700 text-white text-sm font-semibold rounded-lg hover:bg-purple-800 transition-colors"
                      >
                        <BarChart2 className="w-4 h-4" /> Fairness Analysis →
                      </Link>
                    )}
                  </div>

                  {isReplaying && (
                    <div className="mt-4 flex items-center gap-2 text-sm text-emerald-700 bg-emerald-50 px-3 py-2 rounded-lg border border-emerald-200">
                      <RefreshCw className="w-4 h-4 animate-spin" />
                      Replay active — transactions are flowing through the Risk Engine. Check{' '}
                      <Link href="/transactions" className="underline font-semibold">Transactions</Link> and{' '}
                      <Link href="/" className="underline font-semibold">Dashboard</Link> for live updates.
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
