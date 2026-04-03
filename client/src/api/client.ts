import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { isLocal } from '@/lib/utils'

const AP_ACCEPT = 'application/activity+json'
const CACHE_TTL = 2 * 60 * 1000 // 2 minutes
const PROXY_BASE = '/proxy'

// --- Simple in-memory cache ---
interface CacheEntry {
  data: unknown
  expiresAt: number
}
const memCache = new Map<string, CacheEntry>()

function cacheGet(key: string): unknown | null {
  const entry = memCache.get(key)
  if (!entry) return null
  if (Date.now() > entry.expiresAt) {
    memCache.delete(key)
    return null
  }
  return entry.data
}

function cacheSet(key: string, data: unknown, ttl: number) {
  memCache.set(key, { data, expiresAt: Date.now() + ttl })
}

// --- Request deduplication ---
class HttpClient {
  private queue: Map<string, Promise<AxiosResponse>> = new Map()
  client: AxiosInstance

  constructor(axiosInstance: AxiosInstance) {
    this.client = axiosInstance
  }

  get<T = unknown>(url: string, config?: object): Promise<AxiosResponse<T>> {
    if (!this.queue.has(url)) {
      const promise = this.client.get<T>(url, config).finally(() => this.queue.delete(url))
      this.queue.set(url, promise as Promise<AxiosResponse>)
    }
    return this.queue.get(url) as Promise<AxiosResponse<T>>
  }

  post<T = unknown>(url: string, data?: unknown, config?: object) {
    return this.client.post<T>(url, data, config)
  }
}

// --- Proxy interceptor ---
function makeProxyRequestInterceptor(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
  if (config.url && !isLocal(config.url)) {
    config.url = `${PROXY_BASE}/?url=${encodeURIComponent(config.url)}`
    config.withCredentials = true
  }
  return config
}

// --- Cache interceptors (for proxyCacheClient) ---
function makeCacheRequestInterceptor(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig | Promise<never> {
  if ((config.method === 'get' || config.method === 'GET') && config.url) {
    const cached = cacheGet(config.url)
    if (cached !== null) {
      const controller = new AbortController()
      controller.abort()
      const err = Object.assign(new Error('cache_hit'), { isCacheHit: true, data: cached, config })
      return Promise.reject(err)
    }
  }
  return config
}

function makeCacheResponseInterceptor(response: AxiosResponse, ttl: number): AxiosResponse {
  if ((response.config.method === 'get' || response.config.method === 'GET') && response.config.url) {
    cacheSet(response.config.url, response.data, ttl)
  }
  return response
}

function cacheErrorInterceptor(error: unknown): Promise<AxiosResponse> | Promise<never> {
  if (error && typeof error === 'object' && 'isCacheHit' in error) {
    const e = error as unknown as { data: unknown; config: InternalAxiosRequestConfig }
    return Promise.resolve({ data: e.data, status: 200, statusText: 'OK (cached)', headers: {}, config: e.config } as AxiosResponse)
  }
  return Promise.reject(error)
}

// --- proxyClient: proxy remote requests, no caching ---
const _proxyAxios = axios.create({ withCredentials: true })
_proxyAxios.interceptors.request.use(makeProxyRequestInterceptor)
export const proxyClient = new HttpClient(_proxyAxios)

// --- proxyCacheClient: proxy + 2-min cache ---
const _proxyCacheAxios = axios.create({ withCredentials: true })
_proxyCacheAxios.interceptors.request.use((config) => {
  config = makeProxyRequestInterceptor(config)
  return makeCacheRequestInterceptor(config) as InternalAxiosRequestConfig
})
_proxyCacheAxios.interceptors.response.use(
  (res) => makeCacheResponseInterceptor(res, CACHE_TTL),
  cacheErrorInterceptor,
)
export const proxyCacheClient = new HttpClient(_proxyCacheAxios)

// --- authClient: no proxy, credentials always, for auth endpoints ---
export const authClient = axios.create({ withCredentials: true })

// Default AP headers helper
export const apHeaders = { Accept: AP_ACCEPT }
