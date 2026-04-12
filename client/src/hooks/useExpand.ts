/**
 * Hooks for expanding AP IRI references into full objects.
 *
 * In ActivityPub any property value may be either an embedded object or just
 * an IRI string.  These hooks normalise both cases: if the value is already an
 * object it is returned as-is; if it is a string the object is fetched via the
 * proxy cache.
 */

import { useObject } from './useObject'
import { useActivity } from './useActivity'
import { useActor } from './useActor'
import type { APObject, APActivity, APActor } from '@/types/activitypub'

export function useExpandObject(value: string | APObject | undefined) {
  const iri = typeof value === 'string' ? value : undefined
  const { data, isLoading, isError } = useObject(iri)
  return {
    data: typeof value === 'string' ? data : value,
    isLoading: typeof value === 'string' ? isLoading : false,
    isError: typeof value === 'string' ? isError : false,
  }
}

export function useExpandActivity(value: string | APActivity | undefined) {
  const iri = typeof value === 'string' ? value : undefined
  const { data, isLoading, isError } = useActivity(iri)
  return {
    data: typeof value === 'string' ? data : value,
    isLoading: typeof value === 'string' ? isLoading : false,
    isError: typeof value === 'string' ? isError : false,
  }
}

export function useExpandActor(value: string | APActor | undefined) {
  const iri = typeof value === 'string' ? value : undefined
  const { data, isLoading, isError } = useActor(iri)
  return {
    data: typeof value === 'string' ? data : value,
    isLoading: typeof value === 'string' ? isLoading : false,
    isError: typeof value === 'string' ? isError : false,
  }
}
