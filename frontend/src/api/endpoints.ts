import { request } from './client'
import type { HealthResponse, LoginRequest, LoginResponse, SignupRequest, User } from './types'

export function getHealth(): Promise<HealthResponse> {
  return request('GET', '/health')
}

export function signup(input: SignupRequest): Promise<User> {
  return request('POST', '/signup', input)
}

export function login(input: LoginRequest): Promise<LoginResponse> {
  return request('POST', '/login', input)
}
