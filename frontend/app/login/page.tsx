'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Shield } from 'lucide-react'

export default function Login() {
  const [email, setEmail] = useState('admin@riskshield.demo')
  const [password, setPassword] = useState('DemoAdmin123!')
  const [error, setError] = useState('')
  const router = useRouter()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
    if (res.ok) {
      const data = await res.json()
      localStorage.setItem('access_token', data.access_token)
      router.push('/')
    } else {
      setError('Invalid credentials')
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-slate-900 rounded-xl border border-slate-800 p-8">
        <div className="flex items-center justify-center gap-2 mb-8">
          <Shield className="w-10 h-10 text-emerald-500" />
          <h1 className="text-2xl font-bold">RiskShield AI</h1>
        </div>
        <p className="text-center text-slate-400 mb-8">Detect AI risk. Explain it. Control it. Prevent it.</p>
        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-sm text-slate-400 mb-1">Email</label>
            <input type="email" value={email} onChange={e => setEmail(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-emerald-500" />
          </div>
          <div>
            <label className="block text-sm text-slate-400 mb-1">Password</label>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-emerald-500" />
          </div>
          {error && <p className="text-red-400 text-sm">{error}</p>}
          <button type="submit" className="w-full bg-emerald-600 hover:bg-emerald-500 text-white font-medium py-2 rounded-lg transition-colors">
            Sign In
          </button>
        </form>
        <p className="text-center text-slate-500 text-xs mt-6">Demo: admin@riskshield.demo / DemoAdmin123!</p>
      </div>
    </div>
  )
}
