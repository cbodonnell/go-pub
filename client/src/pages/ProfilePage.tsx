import { useParams, useSearchParams } from 'react-router-dom'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { ProfileHeader } from '@/components/profile/ProfileHeader'
import { FeedTimeline } from '@/components/feed/FeedTimeline'
import { useActor } from '@/hooks/useActor'
import { ProfileHeaderSkeleton } from '@/components/common/Skeleton'
import { ErrorState } from '@/components/common/ErrorState'
import { FileText, Heart, Users, UserCheck } from 'lucide-react'

export function ProfilePage() {
  const { name } = useParams<{ name: string }>()
  const [search] = useSearchParams()
  const remoteIri = search.get('remote')

  const iri = remoteIri ?? (name ? `/users/${name}` : undefined)

  const { data: actor, isLoading, isError, refetch } = useActor(iri)

  if (isLoading) return <ProfileHeaderSkeleton />
  if (isError || !actor) return <ErrorState onRetry={() => void refetch()} />

  return (
    <div>
      <ProfileHeader actor={actor} />

      <Tabs defaultValue="posts">
        <div className="sticky top-0 z-10 bg-white/80 dark:bg-slate-950/80 backdrop-blur border-b border-slate-200 dark:border-slate-800">
          <TabsList className="px-4">
            <TabsTrigger value="posts">
              <FileText className="h-3.5 w-3.5" />
              Posts
            </TabsTrigger>
            <TabsTrigger value="liked">
              <Heart className="h-3.5 w-3.5" />
              Liked
            </TabsTrigger>
            <TabsTrigger value="followers">
              <Users className="h-3.5 w-3.5" />
              Followers
            </TabsTrigger>
            <TabsTrigger value="following">
              <UserCheck className="h-3.5 w-3.5" />
              Following
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="posts">
          <FeedTimeline iri={actor.outbox} itemType="Activity" />
        </TabsContent>
        <TabsContent value="liked">
          <FeedTimeline iri={actor.liked} itemType="Object" />
        </TabsContent>
        <TabsContent value="followers">
          <FeedTimeline iri={actor.followers} itemType="Actor" />
        </TabsContent>
        <TabsContent value="following">
          <FeedTimeline iri={actor.following} itemType="Actor" />
        </TabsContent>
      </Tabs>
    </div>
  )
}
