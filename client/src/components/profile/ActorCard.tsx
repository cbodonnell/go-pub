import { useActor } from '@/hooks/useActor'
import { ActorAvatar } from './ActorAvatar'
import { IriLink } from '@/components/common/IriLink'
import { getDisplayName, getHandle } from '@/lib/utils'
import { ActorCardSkeleton } from '@/components/common/Skeleton'

interface ActorCardProps {
  iri?: string
  actor?: import('@/types/activitypub').APActor
  showBio?: boolean
}

function ActorCardInner({ actor, showBio }: { actor: import('@/types/activitypub').APActor; showBio?: boolean }) {
  return (
    <div className="flex items-start gap-3 p-3 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
      <IriLink iri={actor.id} type="actor">
        <ActorAvatar actor={actor} />
      </IriLink>
      <div className="flex-1 min-w-0">
        <IriLink iri={actor.id} type="actor" className="hover:underline">
          <p className="font-semibold text-sm text-slate-900 dark:text-slate-100 truncate">{getDisplayName(actor)}</p>
        </IriLink>
        <p className="text-xs text-slate-500 dark:text-slate-400 truncate">{getHandle(actor)}</p>
        {showBio && actor.summary && (
          <p
            className="text-sm text-slate-600 dark:text-slate-300 mt-1 line-clamp-2"
            dangerouslySetInnerHTML={{ __html: actor.summary }}
          />
        )}
      </div>
    </div>
  )
}

export function ActorCard({ iri, actor: actorProp, showBio }: ActorCardProps) {
  const { data: fetchedActor, isLoading } = useActor(actorProp ? undefined : iri)
  const actor = actorProp ?? fetchedActor

  if (isLoading) return <ActorCardSkeleton />
  if (!actor) return null

  return <ActorCardInner actor={actor} showBio={showBio} />
}
