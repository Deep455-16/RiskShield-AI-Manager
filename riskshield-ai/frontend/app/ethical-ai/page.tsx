'use client'

import { useEffect, useState } from 'react'
import { Scale, CheckCircle, AlertTriangle, BarChart2, Users, Target, Shield, Info } from 'lucide-react'

interface FairnessMetrics {
  selection_rate: number
  true_positive_rate: number
  false_positive_rate: number
  demographic_parity_difference: number
  equal_opportunity_difference: number
  status: string
}

function MetricCard({ label, value, good, format = 'pct', info }: {
  label: string
  value: number
  good?: boolean
  format?: 'pct' | 'num'
  info?: string
}) {
  const display = format === 'pct' ? `${(value * 100).toFixed(1)}%` : value.toFixed(3)
  return (
    <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
      <div className="flex items-start justify-between mb-3">
        <p className="text-xs text-slate-400 uppercase tracking-wider font-semibold">{label}</p>
        {good !== undefined && (
          <span className={`px-2 py-0.5 rounded text-xs font-bold ${good ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'}`}>
            {good ? 'PASS' : 'REVIEW'}
          </span>
        )}
      </div>
      <p className="text-3xl font-bold text-slate-900">{display}</p>
      {info && <p className="text-xs text-slate-500 mt-1">{info}</p>}
    </div>
  )
}

export default function EthicalAIPage() {
  const [metrics, setMetrics] = useState<FairnessMetrics | null>(null)
  const [loading, setLoading] = useState(false)
  const [ran, setRan] = useState(false)

  const runEvaluation = () => {
    setLoading(true)
    fetch('/api/v1/ethical-ai/evaluate')
      .then(r => r.json())
      .then(d => { setMetrics(d.data); setRan(true) })
      .finally(() => setLoading(false))
  }

  const dpd = metrics?.demographic_parity_difference ?? 0
  const eod = metrics?.equal_opportunity_difference ?? 0

  return (
    <div className="min-h-screen bg-slate-50 p-8 pb-20">
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <div className="bg-purple-700 p-2 rounded-lg">
            <Scale className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Ethical AI — Fairness Evaluation</h1>
            <p className="text-slate-500">Based on the UCI Bank Marketing Dataset (features.csv + targets.csv)</p>
          </div>
        </div>
        <div className="mt-4 p-4 bg-purple-50 border border-purple-100 rounded-lg text-sm text-purple-900">
          <p className="font-semibold mb-1 flex items-center gap-2"><Info className="w-4 h-4" /> About this module</p>
          <p>The Bank Marketing dataset is <strong>not</strong> used for transaction fraud detection. It is exclusively used here to evaluate AI fairness metrics — measuring whether our AI-driven decisions achieve demographic parity, equal opportunity, and equalized odds across protected attribute groups (age, marital status, job type).</p>
        </div>
      </div>

      {!ran && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-12 text-center">
          <Scale className="w-12 h-12 text-purple-400 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-slate-900 mb-2">Ready to Evaluate</h2>
          <p className="text-slate-500 mb-6 max-w-md mx-auto">
            Clicking the button below will load the local Bank Marketing dataset, run the AI model against it,
            and compute fairness metrics across demographic groups.
          </p>
          <button
            onClick={runEvaluation}
            disabled={loading}
            className="px-6 py-3 bg-purple-700 text-white font-semibold rounded-lg hover:bg-purple-800 disabled:opacity-50 transition-colors flex items-center gap-2 mx-auto"
          >
            <BarChart2 className="w-5 h-5" />
            {loading ? 'Running Evaluation…' : 'Run Fairness Evaluation'}
          </button>
        </div>
      )}

      {metrics && (
        <div className="space-y-8">
          {/* Overall Status */}
          <div className={`flex items-center gap-4 p-5 rounded-xl border ${metrics.status === 'compliant' ? 'bg-emerald-50 border-emerald-200' : 'bg-rose-50 border-rose-200'}`}>
            {metrics.status === 'compliant' ? (
              <CheckCircle className="w-8 h-8 text-emerald-600" />
            ) : (
              <AlertTriangle className="w-8 h-8 text-rose-600" />
            )}
            <div>
              <p className="font-bold text-lg text-slate-900">Overall Status: {metrics.status === 'compliant' ? 'COMPLIANT' : 'FAIRNESS ALERT'}</p>
              <p className="text-sm text-slate-600">
                {metrics.status === 'compliant'
                  ? 'All fairness metrics are within acceptable thresholds. The AI system treats all demographic groups equitably.'
                  : 'One or more fairness thresholds have been exceeded. Review the metrics below and investigate the affected groups.'}
              </p>
            </div>
          </div>

          {/* Metrics Grid */}
          <div>
            <h2 className="text-lg font-bold text-slate-900 mb-4">Performance Metrics</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <MetricCard
                label="Selection Rate"
                value={metrics.selection_rate}
                info="% of all applicants approved by the AI model"
              />
              <MetricCard
                label="True Positive Rate (Recall)"
                value={metrics.true_positive_rate}
                good={metrics.true_positive_rate > 0.75}
                info="% of genuinely good applicants correctly approved"
              />
              <MetricCard
                label="False Positive Rate"
                value={metrics.false_positive_rate}
                good={metrics.false_positive_rate < 0.2}
                info="% of bad applicants incorrectly approved — aim for low"
              />
            </div>
          </div>

          <div>
            <h2 className="text-lg font-bold text-slate-900 mb-4">Fairness Metrics (Across Demographic Groups)</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <MetricCard
                label="Demographic Parity Difference"
                value={dpd}
                format="num"
                good={dpd < 0.1}
                info="Difference in selection rates across groups. < 0.1 is acceptable."
              />
              <MetricCard
                label="Equal Opportunity Difference"
                value={eod}
                format="num"
                good={eod < 0.1}
                info="Difference in TPR across groups. < 0.1 is acceptable."
              />
            </div>
          </div>

          <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
              <Shield className="w-5 h-5 text-purple-600" /> What Does This Mean?
            </h2>
            <div className="space-y-3 text-sm text-slate-600">
              <p><strong className="text-slate-800">Demographic Parity Difference ({dpd.toFixed(3)}):</strong> Our AI approves loan/subscription applications at nearly identical rates across all demographic groups. A value below 0.1 means no group is being systematically disadvantaged.</p>
              <p><strong className="text-slate-800">Equal Opportunity Difference ({eod.toFixed(3)}):</strong> For customers who genuinely qualify, the AI identifies them at a consistent rate regardless of their age, marital status, or occupation.</p>
              <p><strong className="text-slate-800">Monitoring Frequency:</strong> This evaluation should be run after every model update, quarterly at minimum, and whenever new demographic segments are onboarded.</p>
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={runEvaluation}
              className="flex items-center gap-2 px-4 py-2 bg-purple-700 text-white text-sm font-semibold rounded-lg hover:bg-purple-800 transition-colors"
            >
              <BarChart2 className="w-4 h-4" /> Re-run Evaluation
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
