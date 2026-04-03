import { proxyClient, apHeaders } from './client'
import type { APCollection, APCollectionPage } from '@/types/activitypub'

export async function fetchCollection(iri: string, withCredentials = false): Promise<APCollection> {
  const res = await proxyClient.get<APCollection>(iri, {
    headers: apHeaders,
    withCredentials,
  })
  return res.data
}

export async function fetchCollectionPage(iri: string, withCredentials = false): Promise<APCollectionPage> {
  const res = await proxyClient.get<APCollectionPage>(iri, {
    headers: apHeaders,
    withCredentials,
  })
  return res.data
}
