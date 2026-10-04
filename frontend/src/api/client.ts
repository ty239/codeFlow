import { clearSession, getSession } from '../auth/session'
import { API_BASE_URL, REQUEST_TIMEOUT_MS } from './config'

// An error response from the API, or a network failure (status 0).
export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

type Method = 'GET' | 'POST' | 'PUT' | 'DELETE'

// Every API call goes through here so security settings are applied in one place.
export async function request<T>(method: Method, path: string, body?: unknown): Promise<T> {
  const headers: Record<string, string> = { Accept: 'application/json' }
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }

  const session = getSession()
  if (session) {
    headers.Authorization = `Bearer ${session.token}`
  }

  let response: Response
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      // Auth uses the Authorization header, never cookies, so cookie-based
      // cross-site request forgery doesn't apply.
      credentials: 'omit',
      // Don't let the browser cache responses that may contain user data.
      cache: 'no-store',
      // Refuse redirects so a request can't be bounced to another host or to plain HTTP.
      redirect: 'error',
      signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
    })
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === 'TimeoutError'
    throw new ApiError(0, timedOut ? 'The server took too long to respond.' : 'Could not reach the server.')
  }

  const data: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    // The server rejected our token (expired or revoked), so log out locally.
    if (response.status === 401 && session) {
      clearSession()
    }
    const message =
      data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
        ? data.error
        : `Request failed (${response.status}).`
    throw new ApiError(response.status, message)
  }

  return data as T
}
