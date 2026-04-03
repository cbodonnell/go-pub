import { useQuery } from '@tanstack/react-query'
import { fetchObject } from '@/api/objects'

export function useObject(iri: string | undefined) {
  return useQuery({
    queryKey: ['object', iri],
    queryFn: () => fetchObject(iri!),
    enabled: !!iri,
    staleTime: 2 * 60 * 1000,
  })
}
