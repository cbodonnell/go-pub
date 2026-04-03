import { cn } from '@/lib/utils'

function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('animate-pulse rounded-md bg-slate-200 dark:bg-slate-700', className)}
      {...props}
    />
  )
}

export function PostCardSkeleton() {
  return (
    <div className="p-4 border-b border-slate-200 dark:border-slate-800 flex gap-3">
      <Skeleton className="h-10 w-10 rounded-full shrink-0" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-32" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
      </div>
    </div>
  )
}

export function ProfileHeaderSkeleton() {
  return (
    <div>
      <Skeleton className="h-32 w-full rounded-none" />
      <div className="px-4 pb-4">
        <Skeleton className="h-20 w-20 rounded-full -mt-10 border-4 border-white dark:border-slate-900" />
        <Skeleton className="h-6 w-40 mt-3" />
        <Skeleton className="h-4 w-24 mt-1" />
        <Skeleton className="h-4 w-full mt-3" />
      </div>
    </div>
  )
}

export function ActorCardSkeleton() {
  return (
    <div className="flex items-center gap-3 p-3">
      <Skeleton className="h-10 w-10 rounded-full" />
      <div className="flex-1 space-y-1.5">
        <Skeleton className="h-4 w-28" />
        <Skeleton className="h-3 w-20" />
      </div>
    </div>
  )
}

export { Skeleton }
