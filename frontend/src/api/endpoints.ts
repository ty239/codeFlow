import { request } from './client'
import type { HealthResponse, LoginRequest, LoginResponse, SignupRequest, User } from './types'

export function getHealth(): Promise<HealthResponse> {
  return request('GET', '/health')
}

