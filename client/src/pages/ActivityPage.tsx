import { useParams, useSearchParams } from 'react-router-dom'
import { useActivity } from '@/hooks/useActivity'
import { ActivityItem } from '@/components/feed/ActivityItem'
import { PostCardSkeleton } from '@/components/common/Skeleton'
import { ErrorState } from '@/components/common/ErrorState'

export function ActivityPage() {
  const { id } = useParams<{ id: string }>()
  const [search] = useSearchParams()
  const remoteIri = search.get('remote')

  const iri = remoteIri ?? (id ? `/activities/${id}` : undefined)

  const { data: activity, isLoading, isError, refetch } = useActivity(iri)

  if (isLoading) return <PostCardSkeleton />
  if (isError || !activity) return <ErrorState onRetry={() => void refetch()} />

  return (
    <div>
      <div className="sticky top-0 z-10 bg-white/80 dark:bg-slate-950/80 backdrop-blur border-b border-slate-200 dark:border-slate-800 px-4 py-3">
        <h1 className="text-base font-semibold text-slate-900 dark:text-slate-100">Activity</h1>
      </div>
      <ActivityItem activity={activity} />
    </div>
  )
}
