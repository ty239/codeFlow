import { useEffect, useMemo, useState, type ReactNode } from 'react'
import * as api from '../api/endpoints'
import type { SignupRequest, User } from '../api/types'
import { AuthContext, type AuthContextValue } from './context'
import { clearSession, getSession, onSessionChange, setSession } from './session'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => getSession()?.user ?? null)

  // Stay in sync when the session changes elsewhere, e.g. the API client
  // clearing it after a 401.
  useEffect(() => onSessionChange((session) => setUser(session?.user ?? null)), [])

  const value = useMemo<AuthContextValue>(() => {
    async function login(username: string, password: string) {
