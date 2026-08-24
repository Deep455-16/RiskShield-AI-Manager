import './globals.css'
import type { Metadata } from 'next'
import NavSidebar from '@/components/NavSidebar'

export const metadata: Metadata = {
  title: 'RiskShield AI',
  description: 'Detect AI risk. Explain it. Control it. Prevent it.',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-slate-50 text-slate-900 flex min-h-screen">
        <NavSidebar />
        <main className="flex-1 min-w-0 flex flex-col">
          {children}
        </main>
      </body>
    </html>
  )
}
