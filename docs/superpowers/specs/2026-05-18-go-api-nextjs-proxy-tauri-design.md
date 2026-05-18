# Go API + Next.js Proxy + Tauri — Architecture Design

**Date:** 2026-05-18  
**Status:** Draft  
**Author:** Claude (via brainstorming skill)

---

## Executive Summary

Migrate from NestJS BFF architecture to a streamlined stack where **Go API is the sole source of truth**. Next.js serves as a lightweight proxy for public routes only. Tauri desktop app calls Go directly.

### Key Decisions

1. **Go API does EVERYTHING** — JWT auth, rate limiting, emails (Resend), validation, CORS
2. **Next.js `/api/go/[...path]`** — Proxy uniquement les routes `/public` (masque l'URL Go du navigateur)
3. **Tauri** — Appelle Go directement (`api.opteatech.fr` en prod, IP-whitelisté)
4. **Aucune logique métier dans Next.js** — juste du forwarding

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTS                              │
│                                                             │
│   Next.js 15 (App Router)      Tauri 2 (Desktop Admin)     │
│   Site vitrine public          App interne seulement        │
│   /api/go/* → proxy léger      Axios → Go direct            │
└──────────────┬──────────────────────────┬───────────────────┘
               │ HTTP REST                │ HTTP REST
               │ (via /api/go)            │ (réseau interne Docker)
               ▼                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go API — Fiber v3                        │
│              Source of truth · Port 3001                    │
│                                                             │
│  Auth JWT · Token client · Rate limit · Email Resend        │
│  Validation · CORS · Logging · Swagger                      │
└──────────────┬──────────────────────────────────────────────┘
               │
       ┌───────┴───────┐
       ▼               ▼
  PostgreSQL         Redis
  (données)    (cache · rate limit · sessions)
```

---

## 1. Go API — Routes et Middleware

### 1.1 Middleware Globaux

**Fichier:** `api/internal/middleware/global.go` (nouveau)

```go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/fiber/v3/middleware/limiter"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/gofiber/fiber/v3/middleware/requestid"
)

// GlobalMiddleware applique tous les middleware globaux
func GlobalMiddleware(app *fiber.App, allowedOrigins []string) {
    app.Use(recover.New())
    app.Use(requestid.New())
    
    app.Use(logger.New(logger.Config{
        Format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
    }))

    // CORS — origines whitelist
    app.Use(cors.New(cors.Config{
        AllowOrigins: allowedOrigins,
        AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Forwarded-For", "X-Request-ID"},
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowCredentials: true,
        MaxAge: 86400,
    }))

    // Rate limit global : 200 req/min par IP
    app.Use(limiter.New(limiter.Config{
        Max:        200,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c fiber.Ctx) string {
            if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
                return "global:" + forwarded
            }
            return "global:" + c.IP()
        },
        LimitReached: func(c fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Trop de requêtes. Réessayez dans une minute.",
                "code":  "RATE_LIMIT_EXCEEDED",
            })
        },
    }))
}
```

### 1.2 CORS — Origines

**Fichier:** `api/internal/config/config.go` (modification)

Ajouter les origines Tauri :

```go
// Origines par défaut :
// - Next.js prod
// - Next.js dev
// - Tauri prod
// - Tauri dev
AllowedOrigins: []string{
    "https://opteatech.fr",
    "https://www.opteatech.fr",
    "http://localhost:3000",
    "tauri://localhost",
    "http://tauri.localhost",
}
```

### 1.3 Admin JWT Middleware (amélioré)

**Fichier:** `api/internal/middleware/auth.go` (modification)

Actuellement le middleware valide juste le token. On doit extraire et stocker les claims :

```go
type AdminClaims struct {
    UserID string `json:"sub"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func AdminJWT(secret string) fiber.Handler {
    return func(c fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Token admin requis.",
                "code":  "ADMIN_TOKEN_MISSING",
            })
        }

        rawToken := strings.TrimPrefix(authHeader, "Bearer ")

        claims := &AdminClaims{}
        token, err := jwt.ParseWithClaims(rawToken, claims, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Token invalide ou expiré.",
                "code":  "ADMIN_TOKEN_INVALID",
            })
        }

        if claims.Role != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Accès refusé.",
                "code":  "FORBIDDEN",
            })
        }

        // Stocker les claims dans le contexte pour les handlers
        c.Locals("admin_user_id", claims.UserID)
        c.Locals("admin_email",   claims.Email)
        return c.Next()
    }
}
```

### 1.4 Routes Complètes

**Fichier:** `api/cmd/server/main.go` (refonte)

Structure finale :

```go
func main() {
    cfg := config.Load()
    app := fiber.New()

    // Middleware globaux
    middleware.GlobalMiddleware(app, cfg.AllowedOrigins)

    // Initialisation DB
    var dbPool *pgxpool.Pool
    if cfg.DatabaseURL != "" {
        pool, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
        if err != nil {
            log.Printf("postgres disabled, fallback to in-memory: %v", err)
        } else {
            dbPool = pool
        }
    }
    defer func() { if dbPool != nil { dbPool.Close() } }()

    // Health check
    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok", "version": "1.0.0"})
    })

    api := app.Group("/api")

    // ══════════════════════════════════════════════════════════════
    // ROUTES PUBLIQUES (via proxy Next.js /api/go)
    // ══════════════════════════════════════════════════════════════
    public := api.Group("/public")

    // Formulaire contact/demande — rate limit strict 3 req/heure
    public.Post("/requests",
        limiter.New(limiter.Config{
            Max: 3, Expiration: 1 * time.Hour,
            KeyGenerator: func(c fiber.Ctx) string {
                ip := c.Get("X-Forwarded-For")
                if ip == "" { ip = c.IP() }
                return "req:" + ip
            },
        }),
        handlers.CreateRequestHandler,
    )

    // Espace client — token client opaque
    public.Get("/client/dashboard",   middleware.ClientAuth(repo), handlers.GetDashboard)
    public.Post("/client/messages",   middleware.ClientAuth(repo), handlers.SendMessage)
    public.Get("/client/deliverables/:id/download", middleware.ClientAuth(repo), handlers.DownloadDeliverable)
    public.Post("/client/quote/accept", middleware.ClientAuth(repo), handlers.AcceptQuote)

    // Renouveau token
    public.Post("/client/request-new-link",
        limiter.New(limiter.Config{
            Max: 2, Expiration: 1 * time.Hour,
            KeyGenerator: func(c fiber.Ctx) string { return "newlink:" + c.IP() },
        }),
        handlers.RequestNewToken,
    )

    // ══════════════════════════════════════════════════════════════
    // AUTH ADMIN (JWT Go)
    // ══════════════════════════════════════════════════════════════
    auth := api.Group("/auth")
    auth.Post("/login",   handlers.Login)
    auth.Post("/refresh", handlers.Refresh)
    auth.Post("/logout",  middleware.AdminJWT(cfg.JWTSecret), handlers.Logout)

    // ══════════════════════════════════════════════════════════════
    // ROUTES ADMIN (JWT requis — Tauri uniquement)
    // ══════════════════════════════════════════════════════════════
    admin := api.Group("/admin", middleware.AdminJWT(cfg.JWTSecret))

    // Dashboard
    admin.Get("/dashboard", handlers.GetDashboard)

    // Demandes
    admin.Get("/requests",                        handlers.ListRequests)
    admin.Get("/requests/:id",                    handlers.GetRequest)
    admin.Put("/requests/:id/status",             handlers.UpdateStatus)
    admin.Put("/requests/:id/progress",           handlers.UpdateProgress)
    admin.Post("/requests/:id/milestones",        handlers.CreateMilestone)
    admin.Put("/requests/:id/milestones/:mid",    handlers.UpdateMilestone)
    admin.Post("/requests/:id/messages",          handlers.SendMessage)
    admin.Post("/requests/:id/deliverables",      handlers.AddDeliverable)
    admin.Post("/requests/:id/quote",             handlers.SetQuote)
    admin.Post("/requests/:id/revoke-token",      handlers.RevokeToken)
    admin.Post("/requests/:id/regenerate-token",  handlers.RegenerateToken)

    // Membres / clients
    admin.Get("/members",        handlers.ListMembers)
    admin.Get("/members/:email", handlers.GetMember)

    // Contacts (CRM)
    admin.Get("/contacts",         handlers.ListContacts)
    admin.Get("/contacts/:id",     handlers.GetContact)
    admin.Post("/contacts",        handlers.CreateContact)
    admin.Put("/contacts/:id",     handlers.UpdateContact)
    admin.Delete("/contacts/:id",  handlers.DeleteContact)
    admin.Post("/contacts/:id/notes", handlers.AddContactNote)
}
```

---

## 2. Next.js — Proxy Routes `/api/go`

### 2.1 Proxy Générique

**Fichier:** `frontend/app/api/go/[...path]/route.ts` (nouveau)

```typescript
import { NextRequest, NextResponse } from 'next/server'

const GO_API_URL = process.env.GO_API_INTERNAL_URL ?? 'http://localhost:3001'

async function proxy(req: NextRequest, { params }: { params: { path: string[] } }) {
  const path    = params.path.join('/')
  const url     = new URL(req.url)
  const goUrl   = `${GO_API_URL}/api/public/${path}${url.search}`

  // Reconstruire le body (éviter de le lire deux fois)
  const body = req.method !== 'GET' && req.method !== 'HEAD'
    ? await req.text()
    : undefined

  const response = await fetch(goUrl, {
    method:  req.method,
    headers: {
      'Content-Type':    'application/json',
      'X-Forwarded-For': req.headers.get('x-forwarded-for') ?? req.ip ?? '',
      'X-Real-IP':       req.ip ?? '',
      // Ne jamais forwarder Authorization — les routes /public n'en ont pas besoin
    },
    body,
    redirect: 'manual',
  })

  const data = await response.text()

  return new NextResponse(data, {
    status: response.status,
    headers: {
      'Content-Type': response.headers.get('Content-Type') ?? 'application/json',
    },
  })
}

export const GET    = proxy
export const POST   = proxy
export const PUT    = proxy
export const PATCH  = proxy
export const DELETE = proxy
```

### 2.2 Client API Utility

**Fichier:** `frontend/lib/api/client.ts` (refonte)

```typescript
const GO_INTERNAL = process.env.GO_API_INTERNAL_URL ?? 'http://localhost:3001'
const GO_PUBLIC   = '/api/go' // proxy Next.js — utilisé côté client browser

// Pour les Server Components (SSR, SSG, Route Handlers) : appel direct
export async function goFetch<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const res = await fetch(`${GO_INTERNAL}/api/public${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    next: { revalidate: 60 },
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error ?? `Go API error ${res.status}`)
  }

  return res.json()
}

// Pour les Client Components : passe par le proxy Next.js
export const apiProxy = {
  get:    (path: string) =>
    fetch(`${GO_PUBLIC}${path}`).then(r => r.json()),
  post:   (path: string, body: unknown) =>
    fetch(`${GO_PUBLIC}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    }).then(r => r.json()),
}
```

### 2.3 Exemples d'Usage

```typescript
// ── Server Component (SSR direct vers Go) ──────────────────────
// app/projets/page.tsx
import { goFetch } from '@/lib/api/client'

export default async function ProjetsPage() {
  const { data } = await goFetch<{ data: Project[] }>('/projects?status=published')
  return <ProjectGrid projects={data} />
}

// ── Client Component (via proxy /api/go) ───────────────────────
// components/ContactForm.tsx
'use client'
import { apiProxy } from '@/lib/api/client'

async function submit(data: ContactFormData) {
  return apiProxy.post('/requests', data)
}

// ── Espace client (dashboard token) ───────────────────────────
// app/suivi/page.tsx — Server Component
import { goFetch } from '@/lib/api/client'
import { cookies } from 'next/headers'

export default async function SuiviPage({ searchParams }: { searchParams: { token?: string } }) {
  const token = searchParams.token
  if (!token) return <TokenExpiredPage />

  try {
    const { data } = await goFetch(`/client/dashboard?token=${token}`)
    return <ClientDashboard data={data} token={token} />
  } catch {
    return <TokenExpiredPage />
  }
}
```

---

## 3. Tauri — Axios vers Go Direct

### 3.1 Client Axios

**Fichier:** `apps/tauri/src/lib/api/client.ts` (nouveau)

```typescript
import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { invoke } from '@tauri-apps/api/core'
import { useAuthStore } from '../store/auth.store'

// En dev  : http://localhost:3001
// En prod : variable injectée au build Tauri (réseau interne Docker)
const GO_URL = import.meta.env.VITE_GO_API_URL ?? 'http://localhost:3001'

export const goClient = axios.create({
  baseURL: `${GO_URL}/api`,
  timeout: 12_000,
  headers: { 'Content-Type': 'application/json' },
})

// Inject JWT admin sur chaque requête
goClient.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  const { accessToken } = useAuthStore.getState()
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`
  }
  return config
})

// Auto-refresh JWT (Go signe, Go vérifie)
let isRefreshing = false
let queue: Array<{ resolve: (t: string) => void; reject: (e: unknown) => void }> = []

goClient.interceptors.response.use(
  (r) => r,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !original._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => queue.push({ resolve, reject }))
          .then((token) => {
            original.headers.Authorization = `Bearer ${token}`
            return goClient(original)
          })
      }

      original._retry = true
      isRefreshing = true

      try {
        const { refreshToken, setTokens, logout } = useAuthStore.getState()
        if (!refreshToken) throw new Error('no refresh token')

        const { data } = await axios.post(`${GO_URL}/api/auth/refresh`, {
          refresh_token: refreshToken,
        })

        const newAccess = data.access_token
        setTokens(newAccess, refreshToken)
        await invoke('store_tokens', {
          tokens: { access_token: newAccess, refresh_token: refreshToken },
        })

        queue.forEach(p => p.resolve(newAccess))
        queue = []

        original.headers.Authorization = `Bearer ${newAccess}`
        return goClient(original)

      } catch (e) {
        queue.forEach(p => p.reject(e))
        queue = []
        await useAuthStore.getState().logout()
        return Promise.reject(e)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)
```

### 3.2 Auth API

**Fichier:** `apps/tauri/src/lib/api/auth.api.ts` (nouveau)

```typescript
import axios from 'axios'
import { invoke } from '@tauri-apps/api/core'
import { useAuthStore } from '../store/auth.store'

const GO_URL = import.meta.env.VITE_GO_API_URL ?? 'http://localhost:3001'

export const authApi = {
  login: async (email: string, password: string) => {
    const { data } = await axios.post(`${GO_URL}/api/auth/login`, { email, password })
    return data
  },

  logout: async () => {
    try {
      const { accessToken } = useAuthStore.getState()
      await axios.post(
        `${GO_URL}/api/auth/logout`,
        {},
        { headers: { Authorization: `Bearer ${accessToken}` } }
      )
    } finally {
      await invoke('clear_tokens')
      useAuthStore.getState().logout()
    }
  },
}
```

### 3.3 Admin API

**Fichier:** `apps/tauri/src/lib/api/requests.api.ts` (nouveau)

```typescript
import { goClient } from './client'

export const adminApi = {
  // Dashboard
  getDashboard: () => goClient.get('/admin/dashboard'),

  // Demandes
  listRequests: (params: Record<string, unknown>) =>
    goClient.get('/admin/requests', { params }),

  getRequest: (id: string) =>
    goClient.get(`/admin/requests/${id}`),

  updateStatus: (id: string, status: string) =>
    goClient.put(`/admin/requests/${id}/status`, { status }),

  updateProgress: (id: string, progress: number) =>
    goClient.put(`/admin/requests/${id}/progress`, { progress }),

  // Jalons
  createMilestone: (requestId: string, data: unknown) =>
    goClient.post(`/admin/requests/${requestId}/milestones`, data),

  updateMilestone: (requestId: string, milestoneId: string, data: unknown) =>
    goClient.put(`/admin/requests/${requestId}/milestones/${milestoneId}`, data),

  // Messages
  getMessages: (requestId: string) =>
    goClient.get(`/admin/requests/${requestId}/messages`),

  sendMessage: (requestId: string, content: string, attachments?: string[]) =>
    goClient.post(`/admin/requests/${requestId}/messages`, { content, attachments }),

  // Devis
  setQuote: (requestId: string, data: unknown) =>
    goClient.post(`/admin/requests/${requestId}/quote`, data),

  // Livrables
  addDeliverable: (requestId: string, data: unknown) =>
    goClient.post(`/admin/requests/${requestId}/deliverables`, data),

  // Sécurité token client
  revokeToken: (requestId: string) =>
    goClient.post(`/admin/requests/${requestId}/revoke-token`),

  regenerateToken: (requestId: string) =>
    goClient.post(`/admin/requests/${requestId}/regenerate-token`),

  // Membres
  listMembers: (params?: Record<string, unknown>) =>
    goClient.get('/admin/members', { params }),

  getMember: (email: string) =>
    goClient.get(`/admin/members/${encodeURIComponent(email)}`),

  // Contacts
  listContacts: (params?: Record<string, unknown>) =>
    goClient.get('/admin/contacts', { params }),

  createContact: (data: unknown) =>
    goClient.post('/admin/contacts', data),

  updateContact: (id: string, data: unknown) =>
    goClient.put(`/admin/contacts/${id}`, data),

  deleteContact: (id: string) =>
    goClient.delete(`/admin/contacts/${id}`),

  addContactNote: (id: string, note: string) =>
    goClient.post(`/admin/contacts/${id}/notes`, { note }),
}
```

---

## 4. Docker & Réseau

### 4.1 docker-compose.yml

**Fichier:** `docker-compose.yml` (refonte)

```yaml
services:

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: opteatech
      POSTGRES_USER: opteatech
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U opteatech"]
      interval: 10s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s

  go-api:
    build:
      context: ./api
      dockerfile: Dockerfile
    environment:
      PORT: 3001
      DATABASE_URL: postgresql://opteatech:${POSTGRES_PASSWORD}@postgres:5432/opteatech?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: ${JWT_SECRET}
      JWT_REFRESH_SECRET: ${JWT_REFRESH_SECRET}
      RESEND_API_KEY: ${RESEND_API_KEY}
      FROM_EMAIL: ${FROM_EMAIL}
      ADMIN_EMAIL: ${ADMIN_EMAIL}
      BASE_URL: https://opteatech.fr
      ALLOWED_ORIGINS: https://opteatech.fr,http://localhost:3000,tauri://localhost,http://tauri.localhost
    depends_on:
      postgres: { condition: service_healthy }
      redis:    { condition: service_healthy }
    expose:
      - "3001"
    # En prod : PAS de ports exposés publiquement

  nextjs:
    build:
      context: ./frontend
    environment:
      NODE_ENV: production
      GO_API_INTERNAL_URL: http://go-api:3001
      NEXT_PUBLIC_GO_API_URL: https://opteatech.fr/api/go
    depends_on:
      - go-api
    ports:
      - "3000:3000"

  nginx:
    image: nginx:alpine
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - go-api
      - nextjs

volumes:
  postgres_data:
  redis_data:
```

### 4.2 nginx.conf

**Fichier:** `nginx.conf` (nouveau)

```nginx
# opteatech.fr → Next.js
# opteatech.fr/api/go/* → Next.js proxy → Go (interne)
# api.opteatech.fr/* → Go (Tauri prod, IP whitelist)

server {
  listen 443 ssl;
  server_name opteatech.fr www.opteatech.fr;

  location / {
    proxy_pass http://nextjs:3000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }
}

server {
  listen 443 ssl;
  server_name api.opteatech.fr;

  # Whitelist IP : uniquement machines Tauri
  allow <IP_FIXE>;
  deny all;

  location / {
    proxy_pass http://go-api:3001;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

---

## 5. Variables d'Environnement

```env
# === Go API ===
PORT=3001
ENV=production
BASE_URL=https://opteatech.fr

DATABASE_URL=postgresql://opteatech:xxxx@postgres:5432/opteatech
REDIS_URL=redis://redis:6379

JWT_SECRET=<min-64-chars-random>
JWT_EXPIRES_IN=15m
JWT_REFRESH_SECRET=<min-64-chars-different>
JWT_REFRESH_EXPIRES_IN=7d

RESEND_API_KEY=re_xxxx
FROM_EMAIL=no-reply@opteatech.fr
ADMIN_EMAIL=hello@opteatech.fr

ALLOWED_ORIGINS=https://opteatech.fr,http://localhost:3000,tauri://localhost,http://tauri.localhost


# === Next.js ===
GO_API_INTERNAL_URL=http://go-api:3001
NEXT_PUBLIC_GO_API_URL=https://opteatech.fr/api/go


# === Tauri (injected at build) ===
VITE_GO_API_URL=https://api.opteatech.fr  # prod (IP whitelistée)
# VITE_GO_API_URL=http://localhost:3001   # dev
```

---

## 6. Sécurité

### 6.1 Rate Limiting

- **Global:** 200 req/min par IP (toutes routes)
- **Contact:** 3 req/heure par IP
- **Nouveau token:** 2 req/heure par IP

### 6.2 CORS

Go whitelist strict :
- `https://opteatech.fr` (Next.js prod)
- `https://www.opteatech.fr` (Next.js prod www)
- `http://localhost:3000` (Next.js dev)
- `tauri://localhost` (Tauri prod)
- `http://tauri.localhost` (Tauri dev)

### 6.3 Token Client

- SHA-256 opaque, 90 jours
- Jamais en clair en DB
- Magic link : `?token=<raw>&request_id=<uuid>`

### 6.4 API Tauri

- `api.opteatech.fr` whitelisté par IP (Nginx ou Cloudflare WAF)
- Go non exposé publiquement en prod (réseau Docker interne)

---

## 7. Migration Steps

1. **Mettre à jour `start-dev.sh`** — retirer références NestJS BFF
2. **Ajouter middleware CORS Tauri** dans `config.go`
3. **Améliorer `AdminJWT`** — extraire et stocker les claims
4. **Créer proxy Next.js** — `/api/go/[...path]/route.ts`
5. **Refondre `lib/api/client.ts`** — distinguer SSR vs client
6. **Créer API client Tauri** — avec auto-refresh JWT
7. **Tester le flux complet** — login → refresh → admin calls

---

## 8. Testing Checklist

- [ ] Contact form → Go direct (rate limit 3/h)
- [ ] Client dashboard → token magique → Go direct
- [ ] Admin login → JWT signé → Tauri direct
- [ ] Auto-refresh JWT → 401 → nouveau token
- [ ] CORS → origines Tauri acceptées
- [ ] Rate limit → X-Forwarded-For forwardé par Next.js
- [ ] Logout → revoke refresh token en DB

---

## 9. Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `api/internal/middleware/global.go` | Créer | Middleware globaux (logger, CORS, rate limit) |
| `api/internal/config/config.go` | Modifier | Ajouter origines Tauri |
| `api/internal/middleware/auth.go` | Modifier | Extraire claims AdminJWT |
| `api/cmd/server/main.go` | Refondre | Nouvelles routes |
| `frontend/app/api/go/[...path]/route.ts` | Créer | Proxy Next.js |
| `frontend/lib/api/client.ts` | Refondre | goFetch SSR + apiProxy client |
| `apps/tauri/src/lib/api/client.ts` | Créer | Axios + interceptors |
| `apps/tauri/src/lib/api/auth.api.ts` | Créer | Login/logout Tauri |
| `apps/tauri/src/lib/api/requests.api.ts` | Créer | Admin API Tauri |
| `docker-compose.yml` | Modifier | Retirer NestJS, ajouter nginx |
| `nginx.conf` | Créer | Routing + IP whitelist |
| `start-dev.sh` | Modifier | Retirer BFF NestJS |
