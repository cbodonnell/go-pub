import { authClient } from './client'
import type { AuthUser } from '@/types/activitypub'

const AUTH_BASE = '/auth'

export async function fetchCurrentUser(): Promise<AuthUser | null> {
  try {
    const res = await authClient.get<AuthUser>(`${AUTH_BASE}/`)
    return res.data
  } catch {
    return null
  }
}

export async function login(username: string, password: string): Promise<AuthUser> {
  const params = new URLSearchParams({ username, password })
  const res = await authClient.post<AuthUser>(`${AUTH_BASE}/login`, params, {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  })
  return res.data
}

export async function logout(): Promise<void> {
  await authClient.get(`${AUTH_BASE}/logout`)
}

export async function register(username: string, password: string): Promise<AuthUser> {
  const params = new URLSearchParams({ username, password })
  const res = await authClient.post<AuthUser>(`${AUTH_BASE}/register`, params, {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  })
  return res.data
}
