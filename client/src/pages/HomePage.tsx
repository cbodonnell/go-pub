import { useState } from 'react'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { FeedTimeline } from '@/components/feed/FeedTimeline'
import { useAuthStore } from '@/store/auth'
import { Rss, Inbox, Send } from 'lucide-react'

export function HomePage() {
  const user = useAuthStore((s) => s.user)
  const [tab, setTab] = useState('feed')

  if (!user) {
    return (
      <div className="p-8 text-center space-y-3">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">Welcome to go-pub</h2>
        <p className="text-slate-500 dark:text-slate-400 text-sm">Sign in to see your personalized feed.</p>
      </div>
    )
  }

  return (
    <div>
      <div className="sticky top-0 z-10 bg-white/80 dark:bg-slate-950/80 backdrop-blur border-b border-slate-200 dark:border-slate-800 px-4 pt-4">
        <h1 className="text-base font-semibold text-slate-900 dark:text-slate-100 mb-3">Home</h1>
        <Tabs value={tab} onValueChange={setTab}>
          <TabsList>
            <TabsTrigger value="feed">
              <Rss className="h-3.5 w-3.5" />
              Feed
            </TabsTrigger>
            <TabsTrigger value="inbox">
              <Inbox className="h-3.5 w-3.5" />
              Inbox
            </TabsTrigger>
            <TabsTrigger value="outbox">
              <Send className="h-3.5 w-3.5" />
              Posts
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {tab === 'feed' && (
        <FeedTimeline
          iri={`/users/${user.username}/feed`}
          itemType="Activity"
          withCredentials
        />
      )}
      {tab === 'inbox' && (
        <FeedTimeline
          iri={`/users/${user.username}/inbox`}
          itemType="Activity"
          withCredentials
        />
      )}
      {tab === 'outbox' && (
        <FeedTimeline
          iri={`/users/${user.username}/outbox`}
          itemType="Activity"
        />
      )}
    </div>
  )
}
