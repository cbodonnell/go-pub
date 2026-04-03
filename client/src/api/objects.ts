import { proxyCacheClient, apHeaders } from './client'
import type { APObject } from '@/types/activitypub'

export async function fetchObject(iri: string): Promise<APObject> {
  const res = await proxyCacheClient.get<APObject>(iri, { headers: apHeaders })
  return res.data
}
