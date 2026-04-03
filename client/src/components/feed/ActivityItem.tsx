import { useActivity } from '@/hooks/useActivity'
import { useObject } from '@/hooks/useObject'
import { useActor } from '@/hooks/useActor'
import { PostCard } from './PostCard'
import { ActorCard } from '@/components/profile/ActorCard'
import { IriLink } from '@/components/common/IriLink'
import { PostCardSkeleton } from '@/components/common/Skeleton'
import { Heart, Repeat2, UserPlus } from 'lucide-react'
import type { APActivity, APObject, APActor } from '@/types/activitypub'

interface ActivityItemProps {
  iri?: string
  activity?: APActivity
}

function BoostLabel({ actorIri, icon: Icon, label }: { actorIri: string; icon: React.ElementType; label: string }) {
  const { data: actor } = useActor(actorIri)
  return (
    <div className="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 px-4 pt-2 -mb-1">
      <Icon className="h-3.5 w-3.5" />
      {actor && (
        <IriLink iri={actor.id} type="actor" className="hover:underline font-medium text-slate-600 dark:text-slate-300">
          {actor.name || actor.preferredUsername}
        </IriLink>
      )}
      <span>{label}</span>
    </div>
  )
}

function ActivityItemInner({ activity }: { activity: APActivity }) {
  const actorIri = typeof activity.actor === 'string' ? activity.actor : activity.actor?.id

  if (activity.type === 'Create' || activity.type === 'Update') {
    const obj = activity.object
    if (!obj) return null
    if (typeof obj === 'string') {
      return <ObjectByIri iri={obj} />
    }
    return <PostCard object={obj as APObject} />
  }

  if (activity.type === 'Announce') {
    const obj = activity.object
    if (!obj) return null
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={Repeat2} label="boosted" />}
        {typeof obj === 'string' ? <ObjectByIri iri={obj} /> : <PostCard object={obj as APObject} />}
      </div>
    )
  }

  if (activity.type === 'Like') {
    const obj = activity.object
    if (!obj) return null
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={Heart} label="liked" />}
        {typeof obj === 'string' ? <ObjectByIri iri={obj} /> : <PostCard object={obj as APObject} />}
      </div>
    )
  }

  if (activity.type === 'Follow') {
    const obj = activity.object
    if (!obj) return null
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={UserPlus} label="followed" />}
        <ActorCard iri={typeof obj === 'string' ? obj : (obj as APActor).id} />
      </div>
    )
  }

  // Default: link to the activity
  return (
    <div className="px-4 py-3 text-sm text-slate-500">
      <IriLink iri={activity.id} type="activity" className="hover:underline text-brand-600">
        {activity.type} activity
      </IriLink>
    </div>
  )
}

function ObjectByIri({ iri }: { iri: string }) {
  const { data: obj, isLoading } = useObject(iri)
  if (isLoading) return <PostCardSkeleton />
  if (!obj) return null
  return <PostCard object={obj} />
}

function ActivityItemFetched({ iri }: { iri: string }) {
  const { data: activity, isLoading } = useActivity(iri)
  if (isLoading) return <PostCardSkeleton />
  if (!activity) return null
  return <ActivityItemInner activity={activity} />
}

export function ActivityItem({ iri, activity }: ActivityItemProps) {
  if (activity) return <ActivityItemInner activity={activity} />
  if (iri) return <ActivityItemFetched iri={iri} />
  return null
}
