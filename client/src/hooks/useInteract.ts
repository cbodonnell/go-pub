import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { postActivity, checkActivity } from '@/api/activities'
import { useAuthStore } from '@/store/auth'
import { toast } from 'sonner'

function buildOutboxIri(username: string) {
  return `/users/${username}/outbox`
}

function buildCheckIri(username: string) {
  return `/users/${username}/check`
}

export function useLikeStatus(objectIri: string | undefined) {
  const user = useAuthStore((s) => s.user)
  return useQuery({
    queryKey: ['check', 'Like', objectIri, user?.username],
    queryFn: () => checkActivity(buildCheckIri(user!.username), 'Like', objectIri!),
    enabled: !!user && !!objectIri,
    staleTime: 30 * 1000,
  })
}

export function useFollowStatus(actorIri: string | undefined) {
  const user = useAuthStore((s) => s.user)
  return useQuery({
    queryKey: ['check', 'Follow', actorIri, user?.username],
    queryFn: () => checkActivity(buildCheckIri(user!.username), 'Follow', actorIri!),
    enabled: !!user && !!actorIri,
    staleTime: 30 * 1000,
  })
}

export function useAnnounceStatus(objectIri: string | undefined) {
  const user = useAuthStore((s) => s.user)
  return useQuery({
    queryKey: ['check', 'Announce', objectIri, user?.username],
    queryFn: () => checkActivity(buildCheckIri(user!.username), 'Announce', objectIri!),
    enabled: !!user && !!objectIri,
    staleTime: 30 * 1000,
  })
}

export function useLike(objectIri: string) {
  const user = useAuthStore((s) => s.user)
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (liked: boolean) => {
      const outbox = buildOutboxIri(user!.username)
      if (liked) {
        return postActivity(outbox, {
          '@context': 'https://www.w3.org/ns/activitystreams',
          type: 'Undo',
          id: '',
          object: objectIri,
        })
      }
      return postActivity(outbox, {
        '@context': 'https://www.w3.org/ns/activitystreams',
        type: 'Like',
        id: '',
        object: objectIri,
      })
    },
    onSuccess: (_data, liked) => {
      qc.invalidateQueries({ queryKey: ['check', 'Like', objectIri] })
      toast.success(liked ? 'Unliked' : 'Liked')
    },
    onError: () => toast.error('Action failed'),
  })
}

export function useFollow(actorIri: string) {
  const user = useAuthStore((s) => s.user)
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (following: boolean) => {
      const outbox = buildOutboxIri(user!.username)
      if (following) {
        return postActivity(outbox, {
          '@context': 'https://www.w3.org/ns/activitystreams',
          type: 'Undo',
          id: '',
          object: actorIri,
        })
      }
      return postActivity(outbox, {
        '@context': 'https://www.w3.org/ns/activitystreams',
        type: 'Follow',
        id: '',
        object: actorIri,
      })
    },
    onSuccess: (_data, following) => {
      qc.invalidateQueries({ queryKey: ['check', 'Follow', actorIri] })
      toast.success(following ? 'Unfollowed' : 'Following')
    },
    onError: () => toast.error('Action failed'),
  })
}

export function useAnnounce(objectIri: string) {
  const user = useAuthStore((s) => s.user)
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (announced: boolean) => {
      const outbox = buildOutboxIri(user!.username)
      if (announced) {
        return postActivity(outbox, {
          '@context': 'https://www.w3.org/ns/activitystreams',
          type: 'Undo',
          id: '',
          object: objectIri,
        })
      }
      return postActivity(outbox, {
        '@context': 'https://www.w3.org/ns/activitystreams',
        type: 'Announce',
        id: '',
        object: objectIri,
      })
    },
    onSuccess: (_data, announced) => {
      qc.invalidateQueries({ queryKey: ['check', 'Announce', objectIri] })
      toast.success(announced ? 'Unboosted' : 'Boosted')
    },
    onError: () => toast.error('Action failed'),
  })
}

export function useCreateNote() {
  const user = useAuthStore((s) => s.user)
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ content, inReplyTo }: { content: string; inReplyTo?: string }) => {
      const outbox = buildOutboxIri(user!.username)
      const object: Record<string, unknown> = {
        type: 'Note',
        content,
        to: ['https://www.w3.org/ns/activitystreams#Public'],
        cc: [`/users/${user!.username}/followers`],
      }
      if (inReplyTo) object.inReplyTo = inReplyTo
      return postActivity(outbox, {
        '@context': 'https://www.w3.org/ns/activitystreams',
        type: 'Create',
        id: '',
        object: object as never,
      })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['collection'] })
      toast.success('Posted!')
    },
    onError: () => toast.error('Failed to post'),
  })
}
