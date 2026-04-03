import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useCreateNote } from '@/hooks/useInteract'
import { useAuthStore } from '@/store/auth'
import { Loader2 } from 'lucide-react'

const MAX_CHARS = 500

interface ComposerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  inReplyTo?: string
}

export function Composer({ open, onOpenChange, inReplyTo }: ComposerProps) {
  const [content, setContent] = useState('')
  const user = useAuthStore((s) => s.user)
  const { mutateAsync, isPending } = useCreateNote()

  const remaining = MAX_CHARS - content.length
  const isOverLimit = remaining < 0
  const isEmpty = content.trim().length === 0

  async function handleSubmit() {
    if (isEmpty || isOverLimit || isPending) return
    await mutateAsync({ content, inReplyTo })
    setContent('')
    onOpenChange(false)
  }

  if (!user) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{inReplyTo ? 'Reply' : 'New post'}</DialogTitle>
        </DialogHeader>
        <div className="space-y-3">
          <textarea
            className="w-full min-h-[120px] resize-none rounded-xl border border-slate-200 dark:border-slate-700 bg-transparent px-3 py-2 text-sm text-slate-900 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition"
            placeholder="What's on your mind?"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            disabled={isPending}
          />
          <div className="flex items-center justify-between">
            <span className={`text-xs ${isOverLimit ? 'text-red-500 font-medium' : 'text-slate-400'}`}>
              {remaining}
            </span>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={() => onOpenChange(false)} disabled={isPending}>
                Cancel
              </Button>
              <Button size="sm" onClick={handleSubmit} disabled={isEmpty || isOverLimit || isPending}>
                {isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
                Post
              </Button>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
