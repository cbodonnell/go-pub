import { useParams, useSearchParams } from 'react-router-dom'
import { useObject } from '@/hooks/useObject'
import { PostCard } from '@/components/feed/PostCard'
import { FeedTimeline } from '@/components/feed/FeedTimeline'
import { PostCardSkeleton } from '@/components/common/Skeleton'
import { ErrorState } from '@/components/common/ErrorState'
import { MessageCircle } from 'lucide-react'

export function PostPage() {
  const { id } = useParams<{ id: string }>()
  const [search] = useSearchParams()
  const remoteIri = search.get('remote')

  const iri = remoteIri ?? (id ? `/objects/${id}` : undefined)

  const { data: object, isLoading, isError, refetch } = useObject(iri)

  if (isLoading) return <PostCardSkeleton />
  if (isError || !object) return <ErrorState onRetry={() => void refetch()} />

  return (
    <div>
      <div className="sticky top-0 z-10 bg-white/80 dark:bg-slate-950/80 backdrop-blur border-b border-slate-200 dark:border-slate-800 px-4 py-3">
        <h1 className="text-base font-semibold text-slate-900 dark:text-slate-100">Post</h1>
      </div>

      <PostCard object={object} />

      {object.replies && (
        <div className="border-t border-slate-200 dark:border-slate-800">
          <div className="flex items-center gap-2 px-4 py-3 text-sm font-medium text-slate-500 dark:text-slate-400">
            <MessageCircle className="h-4 w-4" />
            Replies
          </div>
          <FeedTimeline
            iri={typeof object.replies === 'string' ? object.replies : object.replies.id}
            itemType="Object"
          />
        </div>
      )}
    </div>
  )
}
