import { useState } from 'react'
import { ActorAvatar } from './ActorAvatar'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { useFollow, useFollowStatus } from '@/hooks/useInteract'
import { getDisplayName, getHandle } from '@/lib/utils'
import { Loader2 } from 'lucide-react'
import type { APActor } from '@/types/activitypub'

interface ProfileHeaderProps {
  actor: APActor
}

function FollowButton({ actor }: { actor: APActor }) {
  const user = useAuthStore((s) => s.user)
  const isOwnProfile = user?.username && actor.preferredUsername === user.username
  const { data: followStatus, isLoading: statusLoading } = useFollowStatus(isOwnProfile ? undefined : actor.id)
  const followMutation = useFollow(actor.id)

  if (!user || isOwnProfile) return null

  return (
    <Button
      variant={followStatus?.exists ? 'outline' : 'default'}
      size="sm"
      disabled={statusLoading || followMutation.isPending}
      onClick={() => followMutation.mutate(followStatus?.exists ?? false)}
    >
      {followMutation.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
      {followStatus?.exists ? 'Following' : 'Follow'}
    </Button>
  )
}

export function ProfileHeader({ actor }: ProfileHeaderProps) {
  const [bioExpanded, setBioExpanded] = useState(false)

  const coverGradient = 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 50%, #2563eb 100%)'

  return (
    <div>
      {/* Cover */}
      <div className="h-32 md:h-40" style={{ background: coverGradient }} />

      {/* Avatar + actions row */}
      <div className="px-4 flex items-end justify-between -mt-10 md:-mt-12 mb-3">
        <ActorAvatar
          actor={actor}
          size="xl"
          className="border-4 border-white dark:border-slate-900 ring-2 ring-slate-200 dark:ring-slate-700"
        />
        <div className="flex gap-2 mt-12">
          <FollowButton actor={actor} />
        </div>
      </div>

      {/* Info */}
      <div className="px-4 pb-4">
        <h1 className="text-xl font-bold text-slate-900 dark:text-slate-100">{getDisplayName(actor)}</h1>
        <p className="text-sm text-slate-500 dark:text-slate-400">{getHandle(actor)}</p>

        {actor.summary && (
          <div className="mt-2">
            <div
              className={`text-sm text-slate-700 dark:text-slate-300 [&_a]:text-brand-600 [&_a]:hover:underline ${!bioExpanded ? 'line-clamp-3' : ''}`}
              dangerouslySetInnerHTML={{ __html: actor.summary }}
            />
            {actor.summary.length > 200 && (
              <button
                type="button"
                className="text-xs text-brand-600 hover:underline mt-0.5"
                onClick={() => setBioExpanded(!bioExpanded)}
              >
                {bioExpanded ? 'Show less' : 'Show more'}
              </button>
            )}
          </div>
        )}

        {/* Stats */}
        <div className="flex gap-4 mt-3 text-sm">
          {actor.following && (
            <a href={actor.following} className="hover:underline">
              <span className="font-semibold text-slate-900 dark:text-slate-100">—</span>{' '}
              <span className="text-slate-500">Following</span>
            </a>
          )}
          {actor.followers && (
            <a href={actor.followers} className="hover:underline">
              <span className="font-semibold text-slate-900 dark:text-slate-100">—</span>{' '}
              <span className="text-slate-500">Followers</span>
            </a>
          )}
        </div>
      </div>
    </div>
  )
}
