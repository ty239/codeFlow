// Where the API lives. Defaults to /api, which the Vite dev server proxies to
// the Go backend. Set VITE_API_URL for production (e.g. https://api.example.com).
const configuredUrl: string = import.meta.env.VITE_API_URL ?? '/api'

function isLocalhost(hostname: string): boolean {
  return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]'
}

// Production builds refuse to talk to the API over plain HTTP, so passwords and
// tokens are always encrypted in transit. Localhost is exempt for `npm run preview`.
function resolveApiBaseUrl(): string {
