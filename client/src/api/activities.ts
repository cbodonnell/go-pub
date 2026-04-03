import { proxyCacheClient, proxyClient, apHeaders } from './client'
import type { APActivity, CheckResponse } from '@/types/activitypub'

export async function fetchActivity(iri: string): Promise<APActivity> {
  const res = await proxyCacheClient.get<APActivity>(iri, { headers: apHeaders })
  return res.data
}

export async function postActivity(outboxIri: string, activity: Record<string, unknown>): Promise<APActivity> {
  const res = await proxyClient.post<APActivity>(outboxIri, activity, {
    headers: {
      'Content-Type': 'application/activity+json',
      Accept: 'application/activity+json',
    },
    withCredentials: true,
  })
  return res.data
}

export async function checkActivity(
  checkIri: string,
  activityType: string,
  objectIri: string,
): Promise<CheckResponse> {
  const url = `${checkIri}?activity=${activityType}&object=${encodeURIComponent(objectIri)}`
  const res = await proxyClient.get<CheckResponse>(url, {
    headers: apHeaders,
    withCredentials: true,
  })
  return res.data
}
