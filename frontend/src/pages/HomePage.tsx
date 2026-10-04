import { useEffect, useState } from 'react'
import { getHealth } from '../api/endpoints'
import { useAuth } from '../auth/context'

type ServerStatus = 'checking' | 'online' | 'offline'

export function HomePage() {
  const { user, logout } = useAuth()
  const [status, setStatus] = useState<ServerStatus>('checking')

  useEffect(() => {
    getHealth()
      .then((res) => setStatus(res.status === 'ok' ? 'online' : 'offline'))
      .catch(() => setStatus('offline'))
  }, [])

  return (
    <main className="home">
      <header className="topbar">
        <span className="brand">codeFlow</span>
        <span className={`status status-${status}`}>API {status}</span>
        <button type="button" className="secondary" onClick={logout}>
          Log out
        </button>
      </header>

      <section className="card">
        <h2>Welcome, {user?.name}</h2>
        <p>
          Signed in as <strong>{user?.username}</strong> ({user?.email}).
        </p>
        <p className="muted">Boards are coming next. Once the backend has board endpoints, they'll show up here.</p>
      </section>
    </main>
  )
}
