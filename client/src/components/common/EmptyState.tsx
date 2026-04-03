import { Inbox } from 'lucide-react'

interface EmptyStateProps {
  message?: string
  icon?: React.ReactNode
}

export function EmptyState({ message = 'Nothing here yet', icon }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 gap-3 text-slate-400 dark:text-slate-500">
      {icon ?? <Inbox className="h-8 w-8" />}
      <p className="text-sm">{message}</p>
    </div>
  )
}
