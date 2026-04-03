import { useState } from 'react'
import { Heart, Repeat2, MessageCircle, MoreHorizontal } from 'lucide-react'
import { useActor } from '@/hooks/useActor'
import { ActorAvatar } from '@/components/profile/ActorAvatar'
import { IriLink } from '@/components/common/IriLink'
import { Button } from '@/components/ui/button'
import { useLike, useLikeStatus, useAnnounce, useAnnounceStatus } from '@/hooks/useInteract'
import { useAuthStore } from '@/store/auth'
import { getDisplayName, getHandle, formatRelativeTime, cn } from '@/lib/utils'
import type { APObject } from '@/types/activitypub'

interface PostCardProps {
  object: APObject
  onReply?: (objectIri: string) => void
  compact?: boolean
}

function PostContent({ content }: { content: string }) {
  return (
    <div
      className="text-sm text-slate-800 dark:text-slate-200 leading-relaxed [&_a]:text-brand-600 [&_a]:hover:underline [&_p]:mb-1"
      dangerouslySetInnerHTML={{ __html: content }}
    />
  )
}

function AttachmentGrid({ attachments }: { attachments: import('@/types/activitypub').APAttachment[] }) {
  if (!attachments.length) return null
  return (
    <div className={cn('mt-3 rounded-xl overflow-hidden grid gap-0.5', attachments.length === 1 ? 'grid-cols-1' : 'grid-cols-2')}>
      {attachments.map((a, i) => {
        if (a.mediaType?.startsWith('image/') || a.type === 'Image') {
          return (
            <a key={i} href={a.url} target="_blank" rel="noreferrer">
              <img
                src={a.url}
                alt={a.name ?? ''}
                className="w-full h-48 object-cover hover:opacity-90 transition-opacity"
              />
            </a>
          )
        }
        return null
      })}
    </div>
  )
}

function ActionBar({ object }: { object: APObject }) {
  const user = useAuthStore((s) => s.user)
  const { data: likeStatus } = useLikeStatus(user ? object.id : undefined)
  const { data: announceStatus } = useAnnounceStatus(user ? object.id : undefined)
  const likeMutation = useLike(object.id)
  const announceMutation = useAnnounce(object.id)

  return (
    <div className="flex items-center gap-1 mt-3 -ml-2">
      <Button
        variant="ghost"
        size="sm"
        className={cn('gap-1.5 text-slate-500 dark:text-slate-400 text-xs', likeStatus?.exists && 'text-red-500')}
        onClick={() => user && likeMutation.mutate(likeStatus?.exists ?? false)}
        disabled={!user || likeMutation.isPending}
      >
        <Heart className={cn('h-4 w-4', likeStatus?.exists && 'fill-current')} />
      </Button>
      <Button
        variant="ghost"
        size="sm"
        className={cn('gap-1.5 text-slate-500 dark:text-slate-400 text-xs', announceStatus?.exists && 'text-green-500')}
        onClick={() => user && announceMutation.mutate(announceStatus?.exists ?? false)}
        disabled={!user || announceMutation.isPending}
      >
        <Repeat2 className="h-4 w-4" />
      </Button>
      <Button variant="ghost" size="sm" className="gap-1.5 text-slate-500 dark:text-slate-400 text-xs">
        <MessageCircle className="h-4 w-4" />
      </Button>
      <div className="flex-1" />
      <IriLink iri={object.id} type="object">
        <Button variant="ghost" size="icon-sm" className="text-slate-400">
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      </IriLink>
    </div>
  )
}

export function PostCard({ object, compact }: PostCardProps) {
  const [expanded] = useState(false)

  const attributedTo = Array.isArray(object.attributedTo)
    ? (object.attributedTo[0] as string | undefined)
    : (object.attributedTo as string | undefined)

  const actorIri = typeof attributedTo === 'string' ? attributedTo : undefined
  const { data: actor } = useActor(actorIri)

  const content = typeof object.content === 'string' ? object.content : ''

  const rawAttachments = object.attachment
  const attachments = rawAttachments
    ? (Array.isArray(rawAttachments) ? rawAttachments : [rawAttachments])
    : []

  return (
    <article className="px-4 py-3 border-b border-slate-200 dark:border-slate-800 hover:bg-slate-50/50 dark:hover:bg-slate-900/50 transition-colors">
      <div className="flex gap-3">
        {actor ? (
          <IriLink iri={actor.id} type="actor" className="shrink-0">
            <ActorAvatar actor={actor} size="md" />
          </IriLink>
        ) : (
          <div className="h-10 w-10 rounded-full bg-slate-200 dark:bg-slate-700 shrink-0" />
        )}
        <div className="flex-1 min-w-0">
          <div className="flex items-baseline gap-1.5 flex-wrap">
            {actor && (
              <>
                <IriLink iri={actor.id} type="actor" className="hover:underline">
                  <span className="font-semibold text-sm text-slate-900 dark:text-slate-100">
                    {getDisplayName(actor)}
                  </span>
                </IriLink>
                <span className="text-xs text-slate-500 dark:text-slate-400 truncate">{getHandle(actor)}</span>
              </>
            )}
            <span className="text-xs text-slate-400 ml-auto shrink-0">
              {formatRelativeTime(object.published)}
            </span>
          </div>

          {object.inReplyTo && (
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Replying to{' '}
              <IriLink iri={typeof object.inReplyTo === 'string' ? object.inReplyTo : object.inReplyTo.id} type="object" className="text-brand-600 hover:underline">
                a post
              </IriLink>
            </p>
          )}

          {content && (
            <div className={cn('mt-1', !expanded && compact && 'line-clamp-5')}>
              <PostContent content={content} />
            </div>
          )}

          {attachments.length > 0 && <AttachmentGrid attachments={attachments as import('@/types/activitypub').APAttachment[]} />}

          {!compact && <ActionBar object={object} />}
        </div>
      </div>
    </article>
  )
}
