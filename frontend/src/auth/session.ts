import type { User } from '../api/types'

// The logged-in session: the JWT plus the user it belongs to.
//
// It's kept in memory and mirrored to sessionStorage so a page refresh doesn't
// log you out. sessionStorage is per-tab and cleared when the tab closes, which
// limits how long a token lingers. Like any storage readable by JavaScript, it
// is exposed if an attacker gets a script onto the page, which is why the
// production build ships a strict Content Security Policy.

export interface Session {
  token: string
  user: User
}

const STORAGE_KEY = 'codeflow.session'

type Listener = (session: Session | null) => void
const listeners = new Set<Listener>()

let current: Session | null = loadSession()

// Reads the JWT's exp claim. This is not verification (only the server can do
// that); it just lets us drop tokens we already know have expired.
function isExpired(token: string): boolean {
  try {
    const payload = token.split('.')[1]
    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    const { exp } = JSON.parse(json) as { exp?: number }
    return typeof exp !== 'number' || exp * 1000 <= Date.now()
  } catch {
    return true
  }
}

function loadSession(): Session | null {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const session = JSON.parse(raw) as Session
    if (!session.token || !session.user || isExpired(session.token)) {
      sessionStorage.removeItem(STORAGE_KEY)
      return null
    }
    return session
  } catch {
    return null
  }
}

export function getSession(): Session | null {
  if (current && isExpired(current.token)) {
    clearSession()
  }
  return current
}

export function setSession(session: Session): void {
  current = session
  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(session))
  } catch {
    // Storage can be unavailable (private mode, blocked); memory still works.
  }
  listeners.forEach((listener) => listener(current))
}

export function clearSession(): void {
  current = null
  try {
    sessionStorage.removeItem(STORAGE_KEY)
  } catch {
    // Nothing to clear if storage is unavailable.
  }
  listeners.forEach((listener) => listener(null))
}

export function onSessionChange(listener: Listener): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}
