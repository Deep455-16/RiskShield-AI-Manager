'use client'

import { useEffect, useState } from 'react'
import { Vote, CheckCircle, XCircle, Clock } from 'lucide-react'

export default function ApprovalsPage() {
  const [workflows, setWorkflows] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/approvals/workflows').then(r => r.json()).then(d => setWorkflows(d.data || [])).finally(() => setLoading(false))
  }, [])

  const vote = async (id: string, decision: string) => {
    await fetch(`/api/v1/approvals/${id}/vote`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ decision }),
    })
    // Refresh
    fetch('/api/v1/approvals/workflows').then(r => r.json()).then(d => setWorkflows(d.data || []))
  }

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center gap-3">
        <div className="bg-violet-700 p-2 rounded-lg"><Vote className="w-6 h-6 text-white" /></div>
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Approval Gates</h1>
          <p className="text-slate-500">Multi-stage approval workflows for critical AI governance decisions</p>
        </div>
      </div>

      {loading ? (
        <div className="text-center text-slate-500 py-20">Loading workflows...</div>
      ) : workflows.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-12 text-center shadow-sm">
          <Vote className="w-12 h-12 text-slate-300 mx-auto mb-4" />
          <h2 className="text-lg font-bold text-slate-900 mb-2">No Pending Approvals</h2>
          <p className="text-slate-500">When AI systems or policy changes require review, approval requests will appear here.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {(workflows as any[]).map((wf: any) => (
            <div key={wf.id} className={`bg-white rounded-xl border shadow-sm overflow-hidden ${wf.status === 'pending' ? 'border-amber-200' : wf.status === 'approved' ? 'border-emerald-200' : 'border-rose-200'}`}>
              <div className="px-6 py-4 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-3 mb-1">
                    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold uppercase
                      ${wf.status === 'pending' ? 'bg-amber-100 text-amber-700' : wf.status === 'approved' ? 'bg-emerald-100 text-emerald-700' : 'bg-rose-100 text-rose-700'}`}>
                      {wf.status === 'pending' ? <Clock className="w-3 h-3" /> : wf.status === 'approved' ? <CheckCircle className="w-3 h-3" /> : <XCircle className="w-3 h-3" />}
                      {wf.status}
                    </span>
                    <h3 className="font-bold text-slate-900 capitalize">{wf.resource_type} Review</h3>
                  </div>
                  <p className="text-sm text-slate-500">Resource: <span className="font-mono text-slate-700">{wf.resource_id}</span></p>
                  <p className="text-xs text-slate-400 mt-1">Min. {wf.min_approvals} approval(s) required · Created {new Date(wf.created_at).toLocaleDateString()}</p>
                </div>
                {wf.status === 'pending' && (
                  <div className="flex gap-3">
                    <button
                      onClick={() => vote(wf.id, 'approve')}
                      className="flex items-center gap-1.5 px-4 py-2 bg-emerald-600 text-white font-semibold rounded-lg hover:bg-emerald-700 text-sm transition-colors"
                    >
                      <CheckCircle className="w-4 h-4" /> Approve
                    </button>
                    <button
                      onClick={() => vote(wf.id, 'reject')}
                      className="flex items-center gap-1.5 px-4 py-2 bg-rose-600 text-white font-semibold rounded-lg hover:bg-rose-700 text-sm transition-colors"
                    >
                      <XCircle className="w-4 h-4" /> Reject
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
