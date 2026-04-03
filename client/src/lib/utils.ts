import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatRelativeTime(date: string | undefined): string {
  if (!date) return ''
  const d = new Date(date)
  const now = Date.now()
  const diff = now - d.getTime()
  const seconds = Math.floor(diff / 1000)
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d`
  return d.toLocaleDateString()
}

export function getDisplayName(actor: { name?: string; preferredUsername?: string; id: string }): string {
  return actor.name || actor.preferredUsername || actor.id
}

export function getHandle(actor: { preferredUsername?: string; id: string }): string {
  if (actor.preferredUsername) {
    try {
      const url = new URL(actor.id)
      return `@${actor.preferredUsername}@${url.hostname}`
    } catch {
      return `@${actor.preferredUsername}`
    }
  }
  return actor.id
}

export function isLocal(url: string): boolean {
  const base = window.location.protocol + '//' + window.location.hostname
  return url.startsWith(base)
}

export function isValidUrl(str: string): boolean {
  try {
    new URL(str)
    return true
  } catch {
    return false
  }
}
