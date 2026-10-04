# codeFlow frontend

React + TypeScript + Vite client for the codeFlow Go API.

## Running locally

```sh
# 1. Start the backend (from the repo root; needs Postgres running)
go run ./cmd/server

# 2. Start the frontend
cd frontend
npm install
npm run dev
```

The dev server forwards `/api/*` to the backend at `http://localhost:8080`
(override with `BACKEND_URL`), so the browser only talks to one origin.

## Scripts

| Command           | What it does                         |
| ----------------- | ------------------------------------ |
| `npm run dev`     | Dev server with hot reload           |
| `npm run build`   | Type-check and build into `dist/`    |
| `npm run preview` | Serve the production build locally   |
| `npm run lint`    | Lint with oxlint                     |

## Layout

```
src/
├── api/
│   ├── client.ts      Single fetch wrapper all API calls go through
│   ├── config.ts      API base URL; enforces HTTPS in production
│   ├── endpoints.ts   One function per backend endpoint
│   └── types.ts       TypeScript shapes of the backend's JSON
├── auth/
│   ├── session.ts       Stores the JWT and current user
│   ├── context.ts       useAuth() hook
│   └── AuthProvider.tsx Login / signup / logout state
└── pages/             AuthPage (login + signup), HomePage
```

## Security

- **Encrypted in transit:** production builds refuse any API URL that isn't
  `https://` (localhost is exempt for `npm run preview`). TLS is what encrypts
  passwords and tokens; the app doesn't add its own payload encryption on top,
  because that wouldn't protect anything HTTPS doesn't already.
- **Token handling:** the JWT lives in memory and `sessionStorage` (per-tab,
  cleared when the tab closes). Expired tokens are discarded, and any 401 from
  the API logs you out.
- **Requests:** every call goes through `api/client.ts`, which sends no cookies,
  refuses redirects, disables caching and times out after 15 seconds.
- **Content Security Policy:** production builds add a CSP that only allows
  scripts and styles from this site and network calls to this site and the API.
- **Secrets:** `.env*` files are gitignored. Only `VITE_*` variables reach the
  browser, so never put secrets in them.

For production, the backend should also be served over HTTPS and send
`Strict-Transport-Security`. A stronger next step is moving the token into an
`HttpOnly` cookie set by the backend so JavaScript can't read it at all.
