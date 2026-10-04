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
