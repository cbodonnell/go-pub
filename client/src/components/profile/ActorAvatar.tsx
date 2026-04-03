import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import type { APActor } from '@/types/activitypub'
import { cn } from '@/lib/utils'

interface ActorAvatarProps {
  actor: APActor
  size?: 'sm' | 'md' | 'lg' | 'xl'
  className?: string
}

const sizeClasses = {
  sm: 'h-8 w-8',
  md: 'h-10 w-10',
  lg: 'h-14 w-14',
  xl: 'h-20 w-20',
}

function getIconUrl(actor: APActor): string | undefined {
  if (!actor.icon) return undefined
  if (typeof actor.icon === 'string') return actor.icon
  if (actor.icon.type === 'Image') return actor.icon.url
  return undefined
}

function getInitials(actor: APActor): string {
  const name = actor.name || actor.preferredUsername || actor.id
  return name.slice(0, 2).toUpperCase()
}

export function ActorAvatar({ actor, size = 'md', className }: ActorAvatarProps) {
  const iconUrl = getIconUrl(actor)
  return (
    <Avatar className={cn(sizeClasses[size], className)}>
      {iconUrl && <AvatarImage src={iconUrl} alt={actor.name || actor.preferredUsername} />}
      <AvatarFallback>{getInitials(actor)}</AvatarFallback>
    </Avatar>
  )
}
