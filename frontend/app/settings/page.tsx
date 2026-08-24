export default function Settings() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-8">
      <h1 className="text-3xl font-bold mb-6">Settings</h1>
      <div className="bg-slate-900 rounded-lg border border-slate-800 p-6 max-w-2xl">
        <h2 className="text-lg font-semibold mb-4">Organization</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm text-slate-400 mb-1">Organization Name</label>
            <input type="text" defaultValue="Demo Fintech Corp" className="w-full bg-slate-800 border border-slate-700 rounded px-3 py-2 text-sm" readOnly />
          </div>
          <div>
            <label className="block text-sm text-slate-400 mb-1">Risk Weights</label>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Fraud</span><span>25%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Security</span><span>25%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Privacy</span><span>15%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Compliance</span><span>15%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Fairness</span><span>10%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Reliability</span><span>5%</span></div>
              <div className="flex justify-between bg-slate-800 rounded px-3 py-2"><span>Operational</span><span>5%</span></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}