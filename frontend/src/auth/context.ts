import { createContext, useContext } from 'react'
import type { SignupRequest, User } from '../api/types'

export interface AuthContextValue {
  user: User | null
  login: (username: string, password: string) => Promise<void>
  signup: (input: SignupRequest) => Promise<void>
  logout: () => void
}

