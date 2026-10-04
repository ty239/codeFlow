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
