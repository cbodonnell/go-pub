import { proxyCacheClient, apHeaders } from './client'
import type { APActor } from '@/types/activitypub'

export async function fetchActor(iri: string): Promise<APActor> {
  const res = await proxyCacheClient.get<APActor>(iri, { headers: apHeaders })
  return res.data
}
