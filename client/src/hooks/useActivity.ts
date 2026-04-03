import { useQuery } from '@tanstack/react-query'
import { fetchActivity } from '@/api/activities'

export function useActivity(iri: string | undefined) {
  return useQuery({
    queryKey: ['activity', iri],
    queryFn: () => fetchActivity(iri!),
    enabled: !!iri,
    staleTime: 2 * 60 * 1000,
  })
}
