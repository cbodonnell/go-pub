import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { BottomNav } from './BottomNav'
import { Toaster } from 'sonner'

export function AppLayout() {
  return (
    <div className="min-h-screen flex max-w-5xl mx-auto">
      <Sidebar />
      <main className="flex-1 min-w-0 border-x border-slate-200 dark:border-slate-800 pb-20 lg:pb-0">
        <Outlet />
      </main>
      {/* Right sidebar placeholder — can add trending/discover later */}
      <div className="hidden xl:block w-80" />
      <BottomNav />
      <Toaster richColors position="bottom-right" />
    </div>
  )
}
