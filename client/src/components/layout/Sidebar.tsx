import { useState } from 'react'
import { NavLink } from 'react-router-dom'
import { Home, User, LogIn, LogOut, Moon, Sun, PenSquare, Inbox } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { useThemeStore } from '@/store/theme'
import { logout } from '@/api/auth'
import { AuthModal } from '@/components/auth/AuthModal'
import { Composer } from '@/components/composer/Composer'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'

function NavItem({ to, icon: Icon, label }: { to: string; icon: React.ElementType; label: string }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        cn(
          'flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-colors w-full',
          isActive
            ? 'bg-brand-50 dark:bg-brand-900/30 text-brand-700 dark:text-brand-400'
            : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-slate-100',
        )
      }
    >
      <Icon className="h-5 w-5 shrink-0" />
      <span>{label}</span>
    </NavLink>
  )
}

export function Sidebar() {
  const user = useAuthStore((s) => s.user)
  const setUser = useAuthStore((s) => s.setUser)
  const { theme, toggle } = useThemeStore()
  const [authOpen, setAuthOpen] = useState(false)
  const [composerOpen, setComposerOpen] = useState(false)

  async function handleLogout() {
    await logout()
    setUser(null)
    toast.success('Signed out')
  }

  return (
    <aside className="hidden lg:flex flex-col gap-1 w-56 xl:w-64 px-2 py-4 sticky top-0 h-screen overflow-y-auto">
      <div className="px-4 mb-4">
        <span className="text-xl font-bold text-brand-600">go-pub</span>
      </div>

      <NavItem to="/" icon={Home} label="Home" />
      {user && <NavItem to={`/users/${user.username}`} icon={User} label="Profile" />}
      {user && <NavItem to={`/users/${user.username}/inbox`} icon={Inbox} label="Inbox" />}

      <div className="flex-1" />

      {user ? (
        <div className="space-y-1">
          <Button
            variant="default"
            className="w-full gap-2"
            onClick={() => setComposerOpen(true)}
          >
            <PenSquare className="h-4 w-4" />
            New post
          </Button>
          <div className="flex items-center gap-2 px-2 py-1.5 rounded-xl">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-slate-900 dark:text-slate-100 truncate">@{user.username}</p>
            </div>
            <button
              type="button"
              onClick={toggle}
              className="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </button>
            <button
              type="button"
              onClick={handleLogout}
              className="p-1.5 rounded-lg text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          <Button variant="default" className="w-full gap-2" onClick={() => setAuthOpen(true)}>
            <LogIn className="h-4 w-4" />
            Sign in
          </Button>
          <button
            type="button"
            onClick={toggle}
            className="flex items-center gap-2 px-3 py-2 w-full rounded-xl text-sm text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          >
            {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            {theme === 'dark' ? 'Light mode' : 'Dark mode'}
          </button>
        </div>
      )}

      <AuthModal open={authOpen} onOpenChange={setAuthOpen} />
      <Composer open={composerOpen} onOpenChange={setComposerOpen} />
    </aside>
  )
}
