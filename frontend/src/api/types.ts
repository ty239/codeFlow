// Shapes of the Go backend's JSON. Keep in sync with internal/api/users.go.

export interface User {
  id: string
  username: string
  email: string
  name: string
  created_at: string
}

export interface SignupRequest {
  username: string
  email: string
  name: string
  password: string
