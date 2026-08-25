'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  Shield, TrendingUp, Activity, Monitor,
  AlertTriangle, FileText, CheckCircle, Lock, Play, Settings,
  Database, Scale, Eye, Building2, ClipboardList, Vote
} from 'lucide-react'

const links = [
  { href: '/',             label: 'Risk Dashboard',       icon: Shield },
  { href: '/transactions', label: 'Transactions & Fraud', icon: Activity },
  { href: '/risk',         label: 'Threat Analysis',      icon: TrendingUp },
  { href: '/incidents',    label: 'Incident & CAPA',      icon: AlertTriangle },
  { href: '/agents',       label: 'Agent Governance',     icon: Monitor },
  { href: '/shadow-ai',    label: 'Shadow AI Discovery',  icon: Eye },
  { href: '/systems',      label: 'AI System Registry',   icon: Monitor },
  { href: '/compliance',   label: 'Compliance Register',  icon: CheckCircle },
  { href: '/policies',     label: 'Policy Management',    icon: FileText },
  { href: '/vendors',      label: 'Vendor Risk',          icon: Building2 },
  { href: '/tasks',        label: 'Task Management',      icon: ClipboardList },
  { href: '/approvals',    label: 'Approval Gates',       icon: Vote },
  { href: '/audit',        label: 'Audit Trail',          icon: Lock },
  { href: '/datahub',      label: 'Data Hub',             icon: Database },
  { href: '/ethical-ai',   label: 'Ethical AI',           icon: Scale },
  { href: '/simulator',    label: 'Risk Simulator',       icon: Play },
]

export default function NavSidebar() {
  const pathname = usePathname()
  if (pathname === '/login') return null

  return (
    <aside className="w-64 shrink-0 bg-white border-r border-slate-200 sticky top-0 h-screen flex flex-col shadow-sm">
      {/* Logo */}
      <div className="p-6 border-b border-slate-100 flex items-center gap-3">
        <div className="bg-blue-900 p-2 rounded-lg">
          <Building2 className="w-6 h-6 text-white shrink-0" />
        </div>
        <div>
          <div className="font-extrabold text-base text-slate-900 tracking-tight">TrustBank AI</div>
          <div className="text-xs text-slate-500 font-medium uppercase tracking-wider mt-0.5">Risk Manager</div>
        </div>
      </div>

      {/* Nav links */}
      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3 px-3">Core Modules</div>
        {links.map(({ href, label, icon: Icon }) => {
          const active = pathname === href
          return (
            <Link
              key={href}
              href={href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-md text-sm font-medium transition-all ${
                active
                  ? 'bg-blue-50 text-blue-700 shadow-sm border border-blue-100'
                  : 'text-slate-600 hover:text-blue-700 hover:bg-slate-50 border border-transparent'
              }`}
            >
              <Icon className={`w-4 h-4 shrink-0 ${active ? 'text-blue-600' : 'text-slate-400'}`} />
              {label}
            </Link>
          )
        })}
      </nav>

      {/* Footer */}
      <div className="p-5 border-t border-slate-100 bg-slate-50">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 font-bold text-sm">
            DB
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-sm font-medium text-slate-900 truncate">Demo Branch</div>
            <div className="text-xs text-slate-500 truncate">System Administrator</div>
          </div>
        </div>
      </div>
    </aside>
  )
}
