import { Link } from 'react-router-dom'
import { isLocal } from '@/lib/utils'

interface IriLinkProps {
  iri: string
  type: 'actor' | 'object' | 'activity' | 'collection'
  collectionItemType?: string
  className?: string
  children: React.ReactNode
}

function iriToLocalRoute(iri: string): string | null {
  try {
    const url = new URL(iri)
    const localBase = window.location.origin
    if (!iri.startsWith(localBase)) return null
    return url.pathname + url.search
  } catch {
    return null
  }
}

function remoteRoute(iri: string, type: IriLinkProps['type'], collectionItemType?: string): string {
  const base = `/${type}?remote=${encodeURIComponent(iri)}`
  return collectionItemType ? `/${type}?type=${collectionItemType}&remote=${encodeURIComponent(iri)}` : base
}

export function IriLink({ iri, type, collectionItemType, className, children }: IriLinkProps) {
  if (!iri) return <>{children}</>

  const localRoute = iriToLocalRoute(iri)

  if (isLocal(iri) && localRoute) {
    return (
      <Link to={localRoute} className={className}>
        {children}
      </Link>
    )
  }

  return (
    <Link to={remoteRoute(iri, type, collectionItemType)} className={className}>
      {children}
    </Link>
  )
}
