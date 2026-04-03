import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { login, register } from '@/api/auth'
import { useAuthStore } from '@/store/auth'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'

interface AuthModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

type Mode = 'login' | 'register'

export function AuthModal({ open, onOpenChange }: AuthModalProps) {
  const [mode, setMode] = useState<Mode>('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [isPending, setIsPending] = useState(false)
  const setUser = useAuthStore((s) => s.setUser)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!username.trim() || !password.trim()) return
    setIsPending(true)
    try {
      const user = mode === 'login' ? await login(username, password) : await register(username, password)
      setUser(user)
      toast.success(mode === 'login' ? 'Welcome back!' : 'Account created!')
      onOpenChange(false)
      setUsername('')
      setPassword('')
    } catch {
      toast.error(mode === 'login' ? 'Invalid credentials' : 'Registration failed')
    } finally {
      setIsPending(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>{mode === 'login' ? 'Sign in' : 'Create account'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-1">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full rounded-lg border border-slate-200 dark:border-slate-700 bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
              autoComplete="username"
              required
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-1">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full rounded-lg border border-slate-200 dark:border-slate-700 bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
              autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
              required
            />
          </div>
          <Button type="submit" className="w-full" disabled={isPending}>
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {mode === 'login' ? 'Sign in' : 'Create account'}
          </Button>
        </form>
        <div className="text-center text-xs text-slate-500 mt-2">
          {mode === 'login' ? (
            <>
              No account?{' '}
              <button type="button" className="text-brand-600 hover:underline font-medium" onClick={() => setMode('register')}>
                Register
              </button>
            </>
          ) : (
            <>
              Have an account?{' '}
              <button type="button" className="text-brand-600 hover:underline font-medium" onClick={() => setMode('login')}>
                Sign in
              </button>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
