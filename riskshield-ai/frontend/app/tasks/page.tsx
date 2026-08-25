'use client'

import { useEffect, useState } from 'react'
import { ClipboardList, PlusCircle, Clock, CheckCircle, AlertCircle } from 'lucide-react'

export default function TasksPage() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [newTask, setNewTask] = useState({ title: '', priority: 'medium' })

  const fetchTasks = () => {
    setLoading(true)
    fetch('/api/v1/tasks').then(r => r.json()).then(d => setTasks(d.data || [])).finally(() => setLoading(false))
  }

  useEffect(() => { fetchTasks() }, [])

  const addTask = async () => {
    if (!newTask.title.trim()) return
    await fetch('/api/v1/tasks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newTask),
    })
    setNewTask({ title: '', priority: 'medium' })
    setShowAdd(false)
    fetchTasks()
  }

  const updateStatus = async (id: string, status: string) => {
    await fetch(`/api/v1/tasks/${id}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    })
    fetchTasks()
  }

  const priorityColor: Record<string, string> = {
    critical: 'bg-rose-100 text-rose-700 border-rose-200',
    high: 'bg-orange-100 text-orange-700 border-orange-200',
    medium: 'bg-amber-100 text-amber-700 border-amber-200',
    low: 'bg-slate-100 text-slate-600 border-slate-200',
  }

  const statusColor: Record<string, string> = {
    open: 'bg-blue-100 text-blue-700',
    in_progress: 'bg-indigo-100 text-indigo-700',
    blocked: 'bg-rose-100 text-rose-700',
    done: 'bg-emerald-100 text-emerald-700',
    cancelled: 'bg-slate-100 text-slate-500',
  }

  return (
    <div className="min-h-screen bg-white p-8 pb-20">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-blue-700 p-2 rounded-lg"><ClipboardList className="w-6 h-6 text-white" /></div>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Task Management</h1>
            <p className="text-slate-500">Track governance tasks, sorted by priority and due date</p>
          </div>
        </div>
        <button onClick={() => setShowAdd(!showAdd)} className="flex items-center gap-2 px-4 py-2 bg-blue-700 text-white font-semibold rounded-lg hover:bg-blue-800 transition-colors">
          <PlusCircle className="w-4 h-4" /> New Task
        </button>
      </div>

      {showAdd && (
        <div className="mb-6 bg-blue-50 border border-blue-200 rounded-xl p-6">
          <h3 className="font-bold text-slate-900 mb-4">Create New Task</h3>
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Task title…"
              value={newTask.title}
              onChange={e => setNewTask(n => ({ ...n, title: e.target.value }))}
              className="flex-1 border border-slate-200 rounded-lg px-4 py-2 text-sm text-slate-900 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <select
              value={newTask.priority}
              onChange={e => setNewTask(n => ({ ...n, priority: e.target.value }))}
              className="border border-slate-200 rounded-lg px-4 py-2 text-sm text-slate-700 bg-white focus:outline-none"
            >
              <option value="critical">Critical</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
            <button onClick={addTask} className="px-4 py-2 bg-blue-700 text-white font-semibold rounded-lg hover:bg-blue-800">Create</button>
          </div>
        </div>
      )}

      {loading ? (
        <div className="text-center text-slate-500 py-20">Loading tasks...</div>
      ) : (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200">
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Priority</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Task</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Status</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold">Due Date</th>
                <th className="px-6 py-4 text-xs uppercase tracking-wider text-slate-500 font-semibold text-right">Update</th>
              </tr>
            </thead>
            <tbody>
              {tasks.length === 0 ? (
                <tr><td colSpan={5} className="px-6 py-12 text-center text-slate-500">No tasks yet. Create one to get started.</td></tr>
              ) : (
                (tasks as any[]).map((t: any) => (
                  <tr key={t.id} className={`border-b border-slate-100 hover:bg-slate-50 ${t.overdue ? 'bg-rose-50/40' : ''}`}>
                    <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded text-xs font-bold uppercase border ${priorityColor[t.priority] || 'bg-slate-100 text-slate-600 border-slate-200'}`}>
                        {t.priority}
                      </span>
                    </td>
                    <td className="px-6 py-4 font-semibold text-slate-900">
                      {t.title}
                      {t.overdue && <span className="ml-2 text-xs text-rose-600 font-bold">⚠ OVERDUE</span>}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-bold ${statusColor[t.status] || 'bg-slate-100 text-slate-600'}`}>{t.status}</span>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600">
                      {t.due_date ? new Date(t.due_date).toLocaleDateString() : '—'}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <select
                        value={t.status}
                        onChange={e => updateStatus(t.id, e.target.value)}
                        className="text-xs border border-slate-200 rounded-lg px-2 py-1 text-slate-700 bg-white"
                      >
                        <option value="open">Open</option>
                        <option value="in_progress">In Progress</option>
                        <option value="blocked">Blocked</option>
                        <option value="done">Done</option>
                        <option value="cancelled">Cancelled</option>
                      </select>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
