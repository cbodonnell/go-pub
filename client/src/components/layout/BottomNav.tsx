import { useState } from 'react'
import { NavLink } from 'react-router-dom'
import { Home, User, PenSquare, Inbox, LogIn } from 'lucide-react'
import { useAuthStore } from '@/store/auth'
import { AuthModal } from '@/components/auth/AuthModal'
import { Composer } from '@/components/composer/Composer'
import { cn } from '@/lib/utils'

function NavItem({ to, icon: Icon, label }: { to: string; icon: React.ElementType; label: string }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        cn(
          'flex flex-col items-center gap-0.5 px-3 py-2 rounded-xl text-xs font-medium transition-colors flex-1',
          isActive
            ? 'text-brand-600 dark:text-brand-400'
            : 'text-slate-500 dark:text-slate-400',
        )
      }
    >
      <Icon className="h-5 w-5" />
      <span>{label}</span>
    </NavLink>
  )
}

export function BottomNav() {
  const user = useAuthStore((s) => s.user)
  const [authOpen, setAuthOpen] = useState(false)
  const [composerOpen, setComposerOpen] = useState(false)

  return (
    <>
      <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/90 dark:bg-slate-900/90 backdrop-blur border-t border-slate-200 dark:border-slate-700 flex items-center px-2 pb-safe">
        <NavItem to="/" icon={Home} label="Home" />
        {user && <NavItem to={`/users/${user.username}/inbox`} icon={Inbox} label="Inbox" />}

        {user ? (
          <button
            type="button"
            onClick={() => setComposerOpen(true)}
            className="flex flex-col items-center gap-0.5 px-3 py-2 flex-1 text-xs font-medium text-brand-600 dark:text-brand-400"
          >
            <PenSquare className="h-5 w-5" />
            <span>Post</span>
          </button>
        ) : (
          <button
            type="button"
            onClick={() => setAuthOpen(true)}
            className="flex flex-col items-center gap-0.5 px-3 py-2 flex-1 text-xs font-medium text-slate-500"
          >
            <LogIn className="h-5 w-5" />
            <span>Sign in</span>
          </button>
        )}

        {user && <NavItem to={`/users/${user.username}`} icon={User} label="Profile" />}
      </nav>

      <AuthModal open={authOpen} onOpenChange={setAuthOpen} />
      <Composer open={composerOpen} onOpenChange={setComposerOpen} />
    </>
  )
}
