import { useInfiniteQuery } from '@tanstack/react-query'
import { fetchCollection, fetchCollectionPage } from '@/api/collections'
import type { APActivity, APActor, APObject } from '@/types/activitypub'

type CollectionItem = string | APObject | APActivity | APActor

export function useCollection(iri: string | undefined, withCredentials = false) {
  return useInfiniteQuery({
    queryKey: ['collection', iri, withCredentials],
    queryFn: async ({ pageParam }: { pageParam: string | null }) => {
      if (pageParam) {
        return fetchCollectionPage(pageParam, withCredentials)
      }
      // First load: fetch collection, then its first page
      const collection = await fetchCollection(iri!, withCredentials)
      if (!collection.first) {
        return { orderedItems: [] as CollectionItem[], next: undefined, totalItems: collection.totalItems }
      }
      const firstIri = typeof collection.first === 'string' ? collection.first : collection.first.id
      const page = await fetchCollectionPage(firstIri, withCredentials)
      return { ...page, totalItems: collection.totalItems }
    },
    initialPageParam: null as string | null,
    getNextPageParam: (lastPage) => {
      if ('next' in lastPage && lastPage.next) return lastPage.next
      return undefined
    },
    enabled: !!iri,
  })
}
