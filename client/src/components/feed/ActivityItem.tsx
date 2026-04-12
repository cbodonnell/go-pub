import { useActivity } from '@/hooks/useActivity'
import { useExpandObject, useExpandActivity, useExpandActor } from '@/hooks/useExpand'
import { PostCard } from './PostCard'
import { ActorCard } from '@/components/profile/ActorCard'
import { IriLink } from '@/components/common/IriLink'
import { PostCardSkeleton } from '@/components/common/Skeleton'
import { Heart, Repeat2, UserPlus, UserCheck, UserX, Undo2 } from 'lucide-react'
import { useActor } from '@/hooks/useActor'
import type { APActivity, APObject, APActor } from '@/types/activitypub'

interface ActivityItemProps {
  iri?: string
  activity?: APActivity
}

// ── Small labelled header shown above boosted/liked/followed items ────────────

function BoostLabel({ actorIri, icon: Icon, label }: { actorIri: string; icon: React.ElementType; label: string }) {
  const { data: actor } = useActor(actorIri)
  return (
    <div className="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 px-4 pt-2 -mb-1">
      <Icon className="h-3.5 w-3.5" />
      {actor ? (
        <IriLink iri={actor.id} type="actor" className="hover:underline font-medium text-slate-600 dark:text-slate-300">
          {actor.name || actor.preferredUsername}
        </IriLink>
      ) : (
        <IriLink iri={actorIri} type="actor" className="hover:underline font-medium text-slate-600 dark:text-slate-300">
          {actorIri}
        </IriLink>
      )}
      <span>{label}</span>
    </div>
  )
}

// ── Lazy-expand helpers ───────────────────────────────────────────────────────

/** Renders an AP object, fetching it by IRI if needed. */
export function ExpandedObject({ value }: { value: string | APObject | undefined }) {
  const { data: obj, isLoading } = useExpandObject(value)
  if (isLoading) return <PostCardSkeleton />
  if (!obj) return null
  return <PostCard object={obj} />
}

/** Renders an AP activity, fetching it by IRI if needed. */
function ExpandedActivity({ value }: { value: string | APActivity | undefined }) {
  const { data: activity, isLoading } = useExpandActivity(value)
  if (isLoading) return <PostCardSkeleton />
  if (!activity) return null
  return <ActivityItemInner activity={activity} />
}

/** Renders an AP actor card, fetching it by IRI if needed. */
function ExpandedActor({ value }: { value: string | APActor | undefined }) {
  const { data: actor } = useExpandActor(value)
  return <ActorCard iri={actor?.id ?? (typeof value === 'string' ? value : value?.id)} actor={actor} />
}

// ── Core renderer ─────────────────────────────────────────────────────────────

function ActivityItemInner({ activity }: { activity: APActivity }) {
  const actorIri = typeof activity.actor === 'string' ? activity.actor : activity.actor?.id

  if (activity.type === 'Create' || activity.type === 'Update') {
    return <ExpandedObject value={activity.object as string | APObject | undefined} />
  }

  if (activity.type === 'Announce') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={Repeat2} label="boosted" />}
        <ExpandedObject value={activity.object as string | APObject | undefined} />
      </div>
    )
  }

  if (activity.type === 'Like') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={Heart} label="liked" />}
        <ExpandedObject value={activity.object as string | APObject | undefined} />
      </div>
    )
  }

  if (activity.type === 'Follow') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={UserPlus} label="followed" />}
        <ExpandedActor value={activity.object as string | APActor | undefined} />
      </div>
    )
  }

  if (activity.type === 'Accept') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={UserCheck} label="accepted" />}
        <ExpandedActivity value={activity.object as string | APActivity | undefined} />
      </div>
    )
  }

  if (activity.type === 'Reject') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={UserX} label="rejected" />}
        <ExpandedActivity value={activity.object as string | APActivity | undefined} />
      </div>
    )
  }

  if (activity.type === 'Undo') {
    return (
      <div>
        {actorIri && <BoostLabel actorIri={actorIri} icon={Undo2} label="undid" />}
        <ExpandedActivity value={activity.object as string | APActivity | undefined} />
      </div>
    )
  }

  // Generic fallback: expand whatever object is attached, or show a labelled link
  const obj = activity.object
  if (obj) {
    const isObjectType = typeof obj === 'object' && 'content' in obj
    return (
      <div>
        {actorIri && (
          <div className="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 px-4 pt-2 -mb-1">
            <span className="font-medium">{activity.type}</span>
          </div>
        )}
        {isObjectType
          ? <ExpandedObject value={obj as APObject} />
          : <ExpandedActivity value={obj as string | APActivity} />}
      </div>
    )
  }

  return (
    <div className="px-4 py-3 text-sm text-slate-500 dark:text-slate-400">
      <IriLink iri={activity.id} type="activity" className="hover:underline text-brand-600">
        {activity.type} activity
      </IriLink>
    </div>
  )
}

// ── Public component ──────────────────────────────────────────────────────────

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
