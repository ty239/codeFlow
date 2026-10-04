import { createContext, useContext } from 'react'
import type { SignupRequest, User } from '../api/types'

export interface AuthContextValue {
  user: User | null
  login: (username: string, password: string) => Promise<void>
  signup: (input: SignupRequest) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)

export function useAuth(): AuthContextValue {
  const value = useContext(AuthContext)
  if (!value) {
    throw new Error('useAuth must be used inside <AuthProvider>')
  }
  return value
}
