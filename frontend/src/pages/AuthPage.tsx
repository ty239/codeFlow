import { useState, type FormEvent } from 'react'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/context'

type Mode = 'login' | 'signup'

// Mirrors the backend's rules so users get instant feedback. The server still
// enforces them; these checks are only for convenience.
function validatePassword(password: string): string | null {
  if (password.length < 8) return 'Password must be at least 8 characters.'
  if (new TextEncoder().encode(password).length > 72) return 'Password must be at most 72 bytes.'
  return null
}

function errorMessage(err: unknown): string {
  return err instanceof ApiError ? err.message : 'Something went wrong. Please try again.'
}

export function AuthPage() {
  const { login, signup } = useAuth()
  const [mode, setMode] = useState<Mode>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const isSignup = mode === 'signup'

  function switchMode() {
    setMode(isSignup ? 'login' : 'signup')
    setError(null)
    setPassword('')
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)

    if (isSignup) {
      const problem = validatePassword(password)
      if (problem) {
        setError(problem)
        return
      }
    }

    setSubmitting(true)
    try {
      if (isSignup) {
        await signup({ username, email, name, password })
      } else {
        await login(username, password)
      }
    } catch (err) {
      setError(errorMessage(err))
      setSubmitting(false)
    }
  }

  return (
    <main className="auth">
      <h1 className="brand">codeFlow</h1>
      <p className="tagline">A whiteboard for data structures and algorithms.</p>

      <form className="card" onSubmit={handleSubmit} noValidate>
        <h2>{isSignup ? 'Create an account' : 'Log in'}</h2>

        <label>
          Username
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            autoComplete="username"
            maxLength={50}
            required
          />
        </label>

        {isSignup && (
          <>
            <label>
              Email
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
                maxLength={255}
                required
              />
            </label>
            <label>
              Name
              <input
                value={name}
                onChange={(e) => setName(e.target.value)}
                autoComplete="name"
                maxLength={255}
                required
              />
            </label>
          </>
        )}

        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete={isSignup ? 'new-password' : 'current-password'}
            required
          />
        </label>

        {error && (
          <p className="error" role="alert">
            {error}
          </p>
        )}

        <button type="submit" disabled={submitting}>
          {submitting ? 'Please wait…' : isSignup ? 'Sign up' : 'Log in'}
        </button>

        <p className="switch">
          {isSignup ? 'Already have an account?' : 'New to codeFlow?'}{' '}
          <button type="button" className="link" onClick={switchMode}>
            {isSignup ? 'Log in' : 'Create an account'}
          </button>
        </p>
      </form>
    </main>
  )
}
