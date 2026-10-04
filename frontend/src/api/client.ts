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
