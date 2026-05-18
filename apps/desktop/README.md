# optea.tech Desktop Admin (Tauri 2)

Desktop admin app scaffold for optea.tech.

## Stack

- Tauri 2 (Rust backend)
- React 19 + Vite + TypeScript strict
- Axios + React Query + Zustand
- Go API direct (`VITE_GO_API_URL`)

## Start

1. Install deps:

```bash
cd apps/desktop
pnpm install
```

2. Run frontend only:

```bash
pnpm dev
```

3. Run Tauri desktop:

```bash
pnpm tauri:dev
```

## Environment

Copy `.env.example` to `.env` and set:

- `VITE_GO_API_URL=http://localhost:3001` (dev)
- `VITE_GO_API_URL=https://api.opteatech.fr` (prod)

## Notes

- Admin auth is Go-direct (`/api/auth/*`)
- Admin routes are Go-direct (`/api/admin/*`)
- Tokens are persisted through Tauri store commands (`store_tokens`, `get_tokens`, `clear_tokens`)
