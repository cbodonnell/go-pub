import { useEffect, useRef } from 'react'
import { useCollection } from '@/hooks/useCollection'
import { ActivityItem } from './ActivityItem'
import { PostCard } from './PostCard'
import { ActorCard } from '@/components/profile/ActorCard'
import { PostCardSkeleton } from '@/components/common/Skeleton'
import { ErrorState } from '@/components/common/ErrorState'
import { EmptyState } from '@/components/common/EmptyState'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import type { CollectionItemType, APActivity, APObject, APActor } from '@/types/activitypub'
import { ACTIVITY_TYPES, ACTOR_TYPES, OBJECT_TYPES } from '@/types/activitypub'

interface FeedTimelineProps {
  iri: string | undefined
  itemType: CollectionItemType
  withCredentials?: boolean
}

function guessItemType(item: unknown): CollectionItemType {
  if (typeof item === 'object' && item !== null && 'type' in item) {
    const t = (item as { type: string }).type
    if (ACTIVITY_TYPES.includes(t as never)) return 'Activity'
    if (ACTOR_TYPES.includes(t as never)) return 'Actor'
    if (OBJECT_TYPES.includes(t as never)) return 'Object'
  }
  return 'Activity'
}

export function FeedTimeline({ iri, itemType, withCredentials }: FeedTimelineProps) {
  const {
    data,
    isLoading,
    isError,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useCollection(iri, withCredentials)

  const sentinelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const sentinel = sentinelRef.current
    if (!sentinel) return
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          void fetchNextPage()
        }
      },
      { threshold: 0.1 },
    )
    observer.observe(sentinel)
    return () => observer.disconnect()
  }, [hasNextPage, isFetchingNextPage, fetchNextPage])

  if (isLoading) {
    return (
      <div>
        {Array.from({ length: 5 }).map((_, i) => (
          <PostCardSkeleton key={i} />
        ))}
      </div>
    )
  }

  if (isError) {
    return <ErrorState onRetry={() => void refetch()} />
  }

  const allItems = data?.pages.flatMap((p) => {
    const page = p as { orderedItems?: unknown[]; items?: unknown[] }
    return page.orderedItems ?? page.items ?? []
  }) ?? []

  if (!allItems.length) {
    return <EmptyState />
  }

  return (
    <div>
      {allItems.map((item, i) => {
        const resolvedType = typeof item === 'object' ? guessItemType(item) : itemType

        if (typeof item === 'string') {
          if (itemType === 'Actor') return <ActorCard key={i} iri={item} />
          if (itemType === 'Object') return <PostCard key={i} object={{ id: item, type: 'Note' }} />
          return <ActivityItem key={i} iri={item} />
        }

        if (resolvedType === 'Actor') return <ActorCard key={i} actor={item as APActor} />
        if (resolvedType === 'Object') return <PostCard key={i} object={item as APObject} />
        return <ActivityItem key={i} activity={item as APActivity} />
      })}

      <div ref={sentinelRef} className="py-4 flex justify-center">
        {isFetchingNextPage && <Loader2 className="h-5 w-5 animate-spin text-slate-400" />}
        {!hasNextPage && allItems.length > 0 && (
          <p className="text-xs text-slate-400">You're all caught up</p>
        )}
      </div>

      {hasNextPage && !isFetchingNextPage && (
        <div className="flex justify-center pb-4">
          <Button variant="outline" size="sm" onClick={() => void fetchNextPage()}>
            Load more
          </Button>
        </div>
      )}
    </div>
  )
}
