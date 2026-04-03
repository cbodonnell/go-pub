import { useQuery } from '@tanstack/react-query'
import { fetchActor } from '@/api/actors'

export function useActor(iri: string | undefined) {
  return useQuery({
    queryKey: ['actor', iri],
    queryFn: () => fetchActor(iri!),
    enabled: !!iri,
    staleTime: 2 * 60 * 1000,
  })
}
