import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv, type Plugin } from 'vite'

// Content Security Policy for production builds. It only allows scripts and
// styles from this site and network calls to this site and the API, which
// blocks injected scripts from running or sending data anywhere else.
// It's build-only because the dev server relies on inline scripts for hot reload.
function contentSecurityPolicy(apiUrl: string | undefined): Plugin {
  const connectSrc = ["'self'"]
  if (apiUrl && /^https?:\/\//.test(apiUrl)) {
    connectSrc.push(new URL(apiUrl).origin)
  }

  const policy = [
    "default-src 'self'",
    "script-src 'self'",
    "style-src 'self'",
    "img-src 'self' data:",
    `connect-src ${connectSrc.join(' ')}`,
    "object-src 'none'",
    "base-uri 'self'",
    "form-action 'self'",
  ].join('; ')

  return {
    name: 'content-security-policy',
    apply: 'build',
    transformIndexHtml(html) {
      return html.replace(
        '<meta charset="UTF-8" />',
        `<meta charset="UTF-8" />\n    <meta http-equiv="Content-Security-Policy" content="${policy}" />`,
      )
    },
  }
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  // In development the browser calls /api on the Vite server, which forwards
  // to the Go backend. Same origin means no CORS and no mixed content.
  const apiProxy = {
    '/api': {
      target: env.BACKEND_URL ?? 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path: string) => path.replace(/^\/api/, ''),
    },
  }

  return {
    plugins: [react(), contentSecurityPolicy(env.VITE_API_URL)],
    server: { proxy: apiProxy },
    preview: { proxy: apiProxy },
  }
})
