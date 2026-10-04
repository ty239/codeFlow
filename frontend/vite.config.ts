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
