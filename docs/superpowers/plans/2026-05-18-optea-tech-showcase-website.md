# optea·tech Showcase Website Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a production-ready showcase website for optea·tech agency with Next.js 15 frontend, Go Fiber API backend, and Refine admin panel.

**Architecture:** Monorepo with three main applications: Next.js 15 (public site), Go Fiber v3 (REST API), Refine v4 (admin dashboard). PostgreSQL database with Redis for rate limiting.

**Tech Stack:** Next.js 15 App Router, TypeScript, Tailwind CSS, shadcn/ui, Go 1.21, Fiber v3, pgx, Refine v4, PostgreSQL, Redis, Resend, JWT auth.

---

## Phase 1: Project Foundation & Infrastructure

### Task 1: Initialize Monorepo Structure

**Files:**
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `.gitignore`
- Create: `README.md`
- Create: `docker-compose.yml`

- [ ] **Step 1: Create root package.json**

```json
{
  "name": "optea-tech",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "concurrently \"pnpm dev:web\" \"pnpm dev:admin\"",
    "dev:web": "pnpm --filter @optea/web dev",
    "dev:admin": "pnpm --filter @optea/admin dev",
    "build": "pnpm build:web && pnpm build:admin",
    "build:web": "pnpm --filter @optea/web build",
    "build:admin": "pnpm --filter @optea/admin build",
    "lint": "pnpm lint:web && pnpm lint:admin",
    "lint:web": "pnpm --filter @optea/web lint",
    "lint:admin": "pnpm --filter @optea/admin lint",
    "typecheck": "pnpm typecheck:web && pnpm typecheck:admin",
    "typecheck:web": "pnpm --filter @optea/web typecheck",
    "typecheck:admin": "pnpm --filter @optea/admin typecheck"
  },
  "devDependencies": {
    "concurrently": "^8.2.2"
  },
  "packageManager": "pnpm@9.0.0"
}
```

- [ ] **Step 2: Create pnpm-workspace.yaml**

```yaml
packages:
  - "apps/*"
  - "packages/*"
```

- [ ] **Step 3: Create .gitignore**

```gitignore
# Dependencies
node_modules/
.pnpm-store/

# Build outputs
dist/
build/
.next/
out/

# Environment
.env
.env.local
.env.*.local

# Logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# Editor
.DS_Store
*.pem
.vscode/
.idea/
*.swp
*.swo
*~

# Testing
coverage/
.turbo

# Go
api/cmd/server/bin/
api/vendor/
*.exe
*.test
*.out
```

- [ ] **Step 4: Create README.md**

```markdown
# optea·tech

Agence web tech — Sites web, logiciels sur mesure & solutions IA.

## Stack

- **Frontend:** Next.js 15 (App Router) + TypeScript + Tailwind CSS
- **Backend:** Go 1.21 + Fiber v3 + pgx
- **Admin:** Refine v4 + @refinedev/antd
- **Database:** PostgreSQL
- **Cache:** Redis

## Development

```bash
# Install dependencies
pnpm install

# Start all dev servers (web + admin)
pnpm dev

# Start only Next.js dev server
pnpm dev:web

# Start only Refine admin
pnpm dev:admin

# Type check
pnpm typecheck

# Lint
pnpm lint
```

## Project Structure

```
optea-tech/
├── apps/
│   ├── web/          # Next.js 15 - Site vitrine
│   └── admin/        # Refine v4 - Dashboard admin
├── api/              # Go Fiber v3 - Backend API
└── docker-compose.yml
```
```

- [ ] **Step 5: Create docker-compose.yml**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: optea-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: optea_tech
      POSTGRES_USER: optea
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-optea_dev_password}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./api/migrations:/docker-entrypoint-initdb.d

  redis:
    image: redis:7-alpine
    container_name: optea-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

- [ ] **Step 6: Commit**

```bash
git add package.json pnpm-workspace.yaml .gitignore README.md docker-compose.yml
git commit -m "chore: initialize monorepo structure"
```

---

### Task 2: Database Schema & Migrations

**Files:**
- Create: `api/migrations/001_init.sql`
- Create: `api/internal/models/project.go`
- Create: `api/internal/models/service.go`
- Create: `api/internal/models/message.go`
- Create: `api/internal/models/testimonial.go`
- Create: `api/internal/models/user.go`

- [ ] **Step 1: Create database migration**

```sql
-- api/migrations/001_init.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Projets / Portfolio
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(200) NOT NULL,
  slug VARCHAR(200) UNIQUE NOT NULL,
  short_description TEXT,
  full_description TEXT,
  cover_image_url TEXT,
  images JSONB DEFAULT '[]'::jsonb,
  tags TEXT[] DEFAULT '{}',
  category VARCHAR(50) CHECK (category IN ('web', 'logiciel', 'ia', 'conseil')),
  client_name VARCHAR(100),
  project_url TEXT,
  github_url TEXT,
  status VARCHAR(20) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'archived')),
  featured BOOLEAN DEFAULT false,
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_category ON projects(category);
CREATE INDEX idx_projects_featured ON projects(featured);
CREATE INDEX idx_projects_slug ON projects(slug);

-- Services proposés
CREATE TABLE services (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  long_description TEXT,
  icon VARCHAR(50),
  color VARCHAR(20),
  features TEXT[] DEFAULT '{}',
  starting_price INTEGER,
  order_index INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_services_active ON services(is_active);
CREATE INDEX idx_services_order ON services(order_index);

-- Messages de contact
CREATE TABLE contact_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  email VARCHAR(200) NOT NULL,
  company VARCHAR(100),
  service_interest VARCHAR(100),
  budget_range VARCHAR(50),
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT false,
  is_replied BOOLEAN DEFAULT false,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_read ON contact_messages(is_read);
CREATE INDEX idx_messages_created ON contact_messages(created_at DESC);

-- Témoignages clients
CREATE TABLE testimonials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_name VARCHAR(100) NOT NULL,
  client_role VARCHAR(100),
  client_company VARCHAR(100),
  client_avatar_url TEXT,
  content TEXT NOT NULL,
  rating INTEGER DEFAULT 5 CHECK (rating BETWEEN 1 AND 5),
  project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
  is_active BOOLEAN DEFAULT true,
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_testimonials_active ON testimonials(is_active);
CREATE INDEX idx_testimonials_order ON testimonials(order_index);

-- Utilisateurs admin
CREATE TABLE admin_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(200) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  name VARCHAR(100),
  role VARCHAR(20) DEFAULT 'admin' CHECK (role IN ('admin', 'editor')),
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON admin_users(email);

-- Insert default admin user (password: admin123 - CHANGE IN PROD)
-- Hash generated with bcrypt cost 12
INSERT INTO admin_users (email, password_hash, name, role) 
VALUES ('admin@optea.tech', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS3MebAJu', 'Admin', 'admin');
```

- [ ] **Step 2: Create Go models - project.go**

```go
// api/internal/models/project.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectCategory string

const (
	ProjectCategoryWeb      ProjectCategory = "web"
	ProjectCategoryLogiciel ProjectCategory = "logiciel"
	ProjectCategoryIA       ProjectCategory = "ia"
	ProjectCategoryConseil  ProjectCategory = "conseil"
)

type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusPublished ProjectStatus = "published"
	ProjectStatusArchived  ProjectStatus = "archived"
)

type Project struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	Title            string          `json:"title" db:"title"`
	Slug             string          `json:"slug" db:"slug"`
	ShortDescription *string         `json:"short_description,omitempty" db:"short_description"`
	FullDescription  *string         `json:"full_description,omitempty" db:"full_description"`
	CoverImageURL    *string         `json:"cover_image_url,omitempty" db:"cover_image_url"`
	Images           []string        `json:"images" db:"images"`
	Tags             []string        `json:"tags" db:"tags"`
	Category         ProjectCategory `json:"category" db:"category"`
	ClientName       *string         `json:"client_name,omitempty" db:"client_name"`
	ProjectURL       *string         `json:"project_url,omitempty" db:"project_url"`
	GitHubURL        *string         `json:"github_url,omitempty" db:"github_url"`
	Status           ProjectStatus   `json:"status" db:"status"`
	Featured         bool            `json:"featured" db:"featured"`
	OrderIndex       int             `json:"order_index" db:"order_index"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateProjectInput struct {
	Title            string          `json:"title" validate:"required,min=1,max=200"`
	Slug             string          `json:"slug" validate:"required,min=1,max=200"`
	ShortDescription *string         `json:"short_description,omitempty"`
	FullDescription  *string         `json:"full_description,omitempty"`
	CoverImageURL    *string         `json:"cover_image_url,omitempty"`
	Images           []string        `json:"images,omitempty"`
	Tags             []string        `json:"tags,omitempty"`
	Category         ProjectCategory `json:"category" validate:"required,oneof=web logiciel ia conseil"`
	ClientName       *string         `json:"client_name,omitempty"`
	ProjectURL       *string         `json:"project_url,omitempty"`
	GitHubURL        *string         `json:"github_url,omitempty"`
	Status           ProjectStatus   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Featured         bool            `json:"featured"`
	OrderIndex       int             `json:"order_index"`
}

type UpdateProjectInput struct {
	Title            *string          `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Slug             *string          `json:"slug,omitempty" validate:"omitempty,min=1,max=200"`
	ShortDescription *string          `json:"short_description,omitempty"`
	FullDescription  *string          `json:"full_description,omitempty"`
	CoverImageURL    *string          `json:"cover_image_url,omitempty"`
	Images           *[]string        `json:"images,omitempty"`
	Tags             *[]string        `json:"tags,omitempty"`
	Category         *ProjectCategory `json:"category,omitempty" validate:"omitempty,oneof=web logiciel ia conseil"`
	ClientName       *string          `json:"client_name,omitempty"`
	ProjectURL       *string          `json:"project_url,omitempty"`
	GitHubURL        *string          `json:"github_url,omitempty"`
	Status           *ProjectStatus   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Featured         *bool            `json:"featured,omitempty"`
	OrderIndex       *int             `json:"order_index,omitempty"`
}
```

- [ ] **Step 3: Create Go models - service.go**

```go
// api/internal/models/service.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Slug            string    `json:"slug" db:"slug"`
	Description     *string   `json:"description,omitempty" db:"description"`
	LongDescription *string   `json:"long_description,omitempty" db:"long_description"`
	Icon            *string   `json:"icon,omitempty" db:"icon"`
	Color           *string   `json:"color,omitempty" db:"color"`
	Features        []string  `json:"features" db:"features"`
	StartingPrice   *int      `json:"starting_price,omitempty" db:"starting_price"`
	OrderIndex      int       `json:"order_index" db:"order_index"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type CreateServiceInput struct {
	Name            string   `json:"name" validate:"required,min=1,max=100"`
	Slug            string   `json:"slug" validate:"required,min=1,max=100"`
	Description     *string  `json:"description,omitempty"`
	LongDescription *string  `json:"long_description,omitempty"`
	Icon            *string  `json:"icon,omitempty"`
	Color           *string  `json:"color,omitempty"`
	Features        []string `json:"features,omitempty"`
	StartingPrice   *int     `json:"starting_price,omitempty"`
	OrderIndex      int      `json:"order_index"`
	IsActive        bool     `json:"is_active"`
}

type UpdateServiceInput struct {
	Name            *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug            *string  `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	Description     *string  `json:"description,omitempty"`
	LongDescription *string  `json:"long_description,omitempty"`
	Icon            *string  `json:"icon,omitempty"`
	Color           *string  `json:"color,omitempty"`
	Features        *[]string `json:"features,omitempty"`
	StartingPrice   *int     `json:"starting_price,omitempty"`
	OrderIndex      *int     `json:"order_index,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}
```

- [ ] **Step 4: Create Go models - message.go**

```go
// api/internal/models/message.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type ContactMessage struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Company       *string   `json:"company,omitempty" db:"company"`
	ServiceInterest *string `json:"service_interest,omitempty" db:"service_interest"`
	BudgetRange   *string   `json:"budget_range,omitempty" db:"budget_range"`
	Message       string    `json:"message" db:"message"`
	IsRead        bool      `json:"is_read" db:"is_read"`
	IsReplied     bool      `json:"is_replied" db:"is_replied"`
	IPAddress     *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent     *string   `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type CreateMessageInput struct {
	Name          string `json:"name" validate:"required,min=1,max=100"`
	Email         string `json:"email" validate:"required,email"`
	Company       *string `json:"company,omitempty"`
	ServiceInterest *string `json:"service_interest,omitempty"`
	BudgetRange   *string `json:"budget_range,omitempty"`
	Message       string `json:"message" validate:"required,min=20"`
	IPAddress     *string `json:"ip_address,omitempty"`
	UserAgent     *string `json:"user_agent,omitempty"`
}
```

- [ ] **Step 5: Create Go models - testimonial.go**

```go
// api/internal/models/testimonial.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Testimonial struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ClientName     string     `json:"client_name" db:"client_name"`
	ClientRole     *string    `json:"client_role,omitempty" db:"client_role"`
	ClientCompany  *string    `json:"client_company,omitempty" db:"client_company"`
	ClientAvatarURL *string   `json:"client_avatar_url,omitempty" db:"client_avatar_url"`
	Content        string     `json:"content" db:"content"`
	Rating         int        `json:"rating" db:"rating"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	OrderIndex     int        `json:"order_index" db:"order_index"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

type CreateTestimonialInput struct {
	ClientName     string     `json:"client_name" validate:"required,min=1,max=100"`
	ClientRole     *string    `json:"client_role,omitempty"`
	ClientCompany  *string    `json:"client_company,omitempty"`
	ClientAvatarURL *string   `json:"client_avatar_url,omitempty"`
	Content        string     `json:"content" validate:"required,min=10"`
	Rating         int        `json:"rating" validate:"required,min=1,max=5"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	IsActive       bool       `json:"is_active"`
	OrderIndex     int        `json:"order_index"`
}

type UpdateTestimonialInput struct {
	ClientName     *string    `json:"client_name,omitempty" validate:"omitempty,min=1,max=100"`
	ClientRole     *string    `json:"client_role,omitempty"`
	ClientCompany  *string    `json:"client_company,omitempty"`
	ClientAvatarURL *string   `json:"client_avatar_url,omitempty"`
	Content        *string    `json:"content,omitempty" validate:"omitempty,min=10"`
	Rating         *int       `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	IsActive       *bool      `json:"is_active,omitempty"`
	OrderIndex     *int       `json:"order_index,omitempty"`
}
```

- [ ] **Step 6: Create Go models - user.go**

```go
// api/internal/models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleEditor UserRole = "editor"
)

type AdminUser struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"` // Never expose password hash
	Name          *string    `json:"name,omitempty" db:"name"`
	Role          UserRole   `json:"role" db:"role"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type CreateUserInput struct {
	Email        string   `json:"email" validate:"required,email"`
	Password     string   `json:"password" validate:"required,min=8"`
	Name         *string  `json:"name,omitempty"`
	Role         UserRole `json:"role,omitempty" validate:"omitempty,oneof=admin editor"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Name  *string   `json:"name,omitempty"`
		Role  UserRole  `json:"role"`
	} `json:"user"`
}
```

- [ ] **Step 7: Commit**

```bash
git add api/migrations/001_init.sql api/internal/models/*.go
git commit -m "feat: add database schema and Go models"
```

---

### Task 3: Go API Server Setup

**Files:**
- Create: `api/go.mod`
- Create: `api/go.sum`
- Create: `api/.env.example`
- Create: `api/internal/config/config.go`
- Create: `api/internal/database/postgres.go`
- Create: `api/internal/database/redis.go`
- Create: `api/cmd/server/main.go`
- Create: `api/Dockerfile`

- [ ] **Step 1: Initialize Go module**

```bash
cd api
go mod init github.com/optea/api
```

- [ ] **Step 2: Create go.mod with dependencies**

```go
// api/go.mod
module github.com/optea/api

go 1.21

require (
	github.com/gofiber/fiber/v3 v3.0.0-beta.1
	github.com/google/uuid v1.5.0
	github.com/jackc/pgx/v5 v5.5.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.3.1
	github.com/go-playground/validator/v10 v10.16.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/rs/zerolog v1.31.0
	github.com/resend/resend-go/v2 v2.0.0
	golang.org/x/crypto v0.17.0
)
```

- [ ] **Step 3: Create config.go**

```go
// api/internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Resend   ResendConfig
}

type ServerConfig struct {
	Port         string
	Env          string
	AllowedOrigins []string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret     string
	ExpiresIn  time.Duration
	RefreshTTL time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type ResendConfig struct {
	APIKey          string
	FromEmail       string
	ToEmail         string
}

func Load() *Config {
	jwtExpires, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "15"))
	refreshTTL, _ := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRES_IN", "7"))

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "3001"),
			Env:          getEnv("ENV", "development"),
			AllowedOrigins: splitEnv(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://optea:optea_dev_password@localhost:5432/optea_tech?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "dev-secret-change-in-prod"),
			ExpiresIn:  time.Duration(jwtExpires) * time.Minute,
			RefreshTTL: time.Duration(refreshTTL) * 24 * time.Hour,
		},
		CORS: CORSConfig{
			AllowedOrigins: splitEnv(getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")),
		},
		Resend: ResendConfig{
			APIKey:    getEnv("RESEND_API_KEY", ""),
			FromEmail: getEnv("CONTACT_EMAIL_FROM", "no-reply@optea.tech"),
			ToEmail:   getEnv("CONTACT_EMAIL_TO", "hello@optea.tech"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func splitEnv(s string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	for _, v := range splitString(s, ",") {
		trimmed := trimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i = start - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%s", c.Server.Port)
}
```

- [ ] **Step 4: Create postgres.go**

```go
// api/internal/database/postgres.go
package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, url string) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL")

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
```

- [ ] **Step 5: Create redis.go**

```go
// api/internal/database/redis.go
package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB(ctx context.Context, url string) (*RedisDB, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Connected to Redis")

	return &RedisDB{Client: client}, nil
}

func (db *RedisDB) Close() {
	db.Client.Close()
}
```

- [ ] **Step 6: Create main.go**

```go
// api/cmd/server/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/optea/api/internal/config"
	"github.com/optea/api/internal/database"
	"github.com/optea/api/internal/handlers"
	"github.com/optea/api/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connections
	pgDB, err := database.NewPostgresDB(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDB.Close()

	redisDB, err := database.NewRedisDB(ctx, cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisDB.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 4 * 1024 * 1024, // 4MB
	})

	// Global middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Public API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.CreateLoginHandler(pgDB, cfg))
	auth.Post("/refresh", handlers.CreateRefreshHandler(pgDB, cfg))
	auth.Post("/logout", middleware.JWTAuth(), handlers.CreateLogoutHandler())

	// Public routes
	public := api.Group("")
	public.Get("/projects", handlers.CreateGetProjectsHandler(pgDB))
	public.Get("/projects/:slug", handlers.CreateGetProjectBySlugHandler(pgDB))
	public.Get("/services", handlers.CreateGetServicesHandler(pgDB))
	public.Get("/testimonials", handlers.CreateGetTestimonialsHandler(pgDB))
	public.Post("/contact", middleware.RateLimitMiddleware(redisDB.Client, 5), handlers.CreateContactHandler(pgDB, cfg))

	// Admin routes (protected)
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth())

	admin.Get("/dashboard", handlers.CreateDashboardHandler(pgDB))

	// Projects admin
	admin.Get("/projects", handlers.CreateGetAllProjectsHandler(pgDB))
	admin.Post("/projects", handlers.CreateProjectHandler(pgDB))
	admin.Put("/projects/:id", handlers.UpdateProjectHandler(pgDB))
	admin.Delete("/projects/:id", handlers.DeleteProjectHandler(pgDB))

	// Services admin
	admin.Get("/services", handlers.CreateGetAllServicesHandler(pgDB))
	admin.Post("/services", handlers.CreateServiceHandler(pgDB))
	admin.Put("/services/:id", handlers.UpdateServiceHandler(pgDB))
	admin.Delete("/services/:id", handlers.DeleteServiceHandler(pgDB))

	// Messages admin
	admin.Get("/messages", handlers.CreateGetMessagesHandler(pgDB))
	admin.Put("/messages/:id/read", handlers.MarkMessageReadHandler(pgDB))
	admin.Delete("/messages/:id", handlers.DeleteMessageHandler(pgDB))

	// Testimonials admin
	admin.Get("/testimonials", handlers.CreateGetAllTestimonialsHandler(pgDB))
	admin.Post("/testimonials", handlers.CreateTestimonialHandler(pgDB))
	admin.Put("/testimonials/:id", handlers.UpdateTestimonialHandler(pgDB))
	admin.Delete("/testimonials/:id", handlers.DeleteTestimonialHandler(pgDB))

	// Start server
	go func() {
		log.Printf("Server starting on %s", cfg.Addr())
		if err := app.Listen(cfg.Addr()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
```

- [ ] **Step 7: Create .env.example**

```env
# Server
PORT=3001
ENV=development

# Database
DATABASE_URL=postgresql://optea:optea_dev_password@localhost:5432/optea_tech?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=change-this-secret-in-production
JWT_EXPIRES_IN=15
REFRESH_TOKEN_EXPIRES_IN=7

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001

# Resend
RESEND_API_KEY=
CONTACT_EMAIL_FROM=no-reply@optea.tech
CONTACT_EMAIL_TO=hello@optea.tech
```

- [ ] **Step 8: Create Dockerfile**

```dockerfile
# api/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/.env.example .env

EXPOSE 3001

CMD ["./server"]
```

- [ ] **Step 9: Commit**

```bash
git add api/go.mod api/go.sum api/.env.example api/internal/config/config.go api/internal/database/*.go api/cmd/server/main.go api/Dockerfile
git commit -m "feat: initialize Go Fiber API server"
```

---

### Task 4: API Handlers Implementation

**Files:**
- Create: `api/internal/handlers/projects.go`
- Create: `api/internal/handlers/services.go`
- Create: `api/internal/handlers/messages.go`
- Create: `api/internal/handlers/testimonials.go`
- Create: `api/internal/handlers/auth.go`
- Create: `api/internal/handlers/dashboard.go`
- Create: `api/internal/middleware/auth.go`
- Create: `api/internal/middleware/ratelimit.go`
- Create: `api/internal/repository/postgres/projects.go`
- Create: `api/internal/repository/postgres/services.go`
- Create: `api/internal/repository/postgres/messages.go`
- Create: `api/internal/repository/postgres/testimonials.go`
- Create: `api/internal/repository/postgres/users.go`

- [ ] **Step 1: Create middleware/auth.go**

```go
// api/internal/middleware/auth.go
package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

func getJWTSecret() []byte {
	secret := getEnv("JWT_SECRET", "dev-secret")
	return []byte(secret)
}

func getEnv(key, defaultValue string) string {
	if value := strings.TrimSpace(value); value != "" {
		return value
	}
	return defaultValue
}

func GenerateJWT(userID uuid.UUID, email, role string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJWTSecret()))
}
```

- [ ] **Step 2: Create middleware/ratelimit.go**

```go
// api/internal/middleware/ratelimit.go
package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(client *redis.Client, maxRequests int) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		key := "ratelimit:" + ip

		ctx := c.Context()
		current, err := client.Get(ctx, key).Int()
		if err == redis.Nil {
			client.Set(ctx, key, 1, time.Minute)
			return c.Next()
		}
		if err != nil {
			return c.Next()
		}

		if current >= maxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}

		client.Incr(ctx, key)
		return c.Next()
	}
}
```

- [ ] **Step 3: Create repository/postgres/projects.go**

```go
// api/internal/repository/postgres/projects.go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea/api/internal/models"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAll(ctx context.Context, status models.ProjectStatus) ([]models.Project, error) {
	query := `
		SELECT id, title, slug, short_description, full_description, cover_image_url, 
		       images, tags, category, client_name, project_url, github_url, 
		       status, featured, order_index, created_at, updated_at
		FROM projects
		WHERE status = $1
		ORDER BY order_index, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID, &p.Title, &p.Slug, &p.ShortDescription, &p.FullDescription,
			&p.CoverImageURL, &p.Images, &p.Tags, &p.Category, &p.ClientName,
			&p.ProjectURL, &p.GitHubURL, &p.Status, &p.Featured, &p.OrderIndex,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *ProjectRepository) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	query := `
		SELECT id, title, slug, short_description, full_description, cover_image_url, 
		       images, tags, category, client_name, project_url, github_url, 
		       status, featured, order_index, created_at, updated_at
		FROM projects
		WHERE slug = $1
	`

	var p models.Project
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.ShortDescription, &p.FullDescription,
		&p.CoverImageURL, &p.Images, &p.Tags, &p.Category, &p.ClientName,
		&p.ProjectURL, &p.GitHubURL, &p.Status, &p.Featured, &p.OrderIndex,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Create(ctx context.Context, input models.CreateProjectInput) (*models.Project, error) {
	id := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO projects (id, title, slug, short_description, full_description, 
		                      cover_image_url, images, tags, category, client_name, 
		                      project_url, github_url, status, featured, order_index, 
		                      created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, title, slug, short_description, full_description, cover_image_url, 
		          images, tags, category, client_name, project_url, github_url, 
		          status, featured, order_index, created_at, updated_at
	`

	var p models.Project
	err := r.db.QueryRow(ctx, query,
		id, input.Title, input.Slug, input.ShortDescription, input.FullDescription,
		input.CoverImageURL, pq.Array(input.Images), pq.Array(input.Tags), input.Category,
		input.ClientName, input.ProjectURL, input.GitHubURL, input.Status, input.Featured,
		input.OrderIndex, now, now,
	).Scan(
		&p.ID, &p.Title, &p.Slug, &p.ShortDescription, &p.FullDescription,
		&p.CoverImageURL, &p.Images, &p.Tags, &p.Category, &p.ClientName,
		&p.ProjectURL, &p.GitHubURL, &p.Status, &p.Featured, &p.OrderIndex,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id uuid.UUID, input models.UpdateProjectInput) (*models.Project, error) {
	// Implementation for update
	return nil, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}
```

[Continuing with remaining handlers...]

Due to length constraints, I'll create the remaining files in a more condensed format. Let me write the complete plan file:
