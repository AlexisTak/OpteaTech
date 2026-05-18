# OpteaTech Complete Implementation Design

**Date:** 2026-05-18  
**Author:** Claude (Brainstorming Skill)  
**Status:** Approved

## Executive Summary

Complete the OpteaTech project by:
1. Setting up monorepo structure with pnpm workspaces
2. Correcting start-dev.sh to remove non-existent NestJS backend reference
3. Creating a complete Refine v4 admin dashboard with Ant Design
4. Initializing git repository and connecting to GitHub remote
5. Setting up docker-compose for PostgreSQL and Redis

**Approach:** Incremental implementation in 3 phases (Infrastructure → Basic Dashboard → Complete Dashboard)

**Timeline:** ~7-9 hours total

## Requirements

### Decisions Made

1. **Dashboard Type:** Refine v4 autonomous dashboard consuming existing Go API endpoints
2. **Monorepo Structure:** Plan A - apps/web, apps/admin, apps/desktop, api/ at root
3. **start-dev.sh:** Remove NestJS backend reference, keep Go API + Web + Desktop
4. **API Architecture:** Keep existing Go admin endpoints, Refine consumes them directly
5. **Git:** Initialize repository now, connect to git@github.com:AlexisTak/OpteaTech.git
6. **Docker:** Create docker-compose with PostgreSQL 16 + Redis 7
7. **Dashboard Resources:** Complete - Projects, Services, Testimonials, Messages, Client Requests + Dashboard stats
8. **UI Library:** Ant Design (@refinedev/antd)

## Architecture

### Global Structure

```
OpteaTech/
├── package.json              # Root workspace config
├── pnpm-workspace.yaml       # pnpm workspaces
├── docker-compose.yml        # PostgreSQL + Redis
├── start-dev.sh             # Launch Go API + Web + Desktop
├── .gitignore               # Git exclusions
├── apps/
│   ├── web/                 # Next.js (frontend renamed)
│   ├── admin/               # Refine dashboard (new)
│   └── desktop/             # Tauri app (existing)
├── api/                     # Go Fiber API (existing, stays at root)
└── docs/
    └── superpowers/
        └── specs/           # Design docs
```

### Communication Flow

- **apps/web** → Go API `http://localhost:3001/api`
- **apps/admin** → Go API `http://localhost:3001/api/admin`  
- **apps/desktop** → Go API (same endpoints)
- All apps consume Go API as single source of truth
- No BFF/NestJS layer - Go API handles everything

### Tech Stack

**Monorepo:**
- pnpm workspaces
- Concurrent script management

**Admin Dashboard:**
- Refine v4 (`@refinedev/core`, `@refinedev/react-router-v6`)
- Ant Design UI (`@refinedev/antd`, `antd`)
- REST Data Provider (`@refinedev/simple-rest`)
- Custom Auth Provider (JWT with Go API)
- React Router v6
- TypeScript + Vite

**Infrastructure:**
- Docker Compose: PostgreSQL 16 + Redis 7
- Git + GitHub remote

## Phase 1: Infrastructure (1-2h)

### Monorepo Setup

**Root package.json:**
```json
{
  "name": "optea-tech",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "concurrently \"pnpm dev:api\" \"pnpm dev:web\" \"pnpm dev:admin\"",
    "dev:api": "cd api && go run ./cmd/server",
    "dev:web": "pnpm --filter @optea/web dev",
    "dev:admin": "pnpm --filter @optea/admin dev",
    "build": "pnpm build:web && pnpm build:admin",
    "build:web": "pnpm --filter @optea/web build",
    "build:admin": "pnpm --filter @optea/admin build"
  },
  "devDependencies": {
    "concurrently": "^8.2.2"
  },
  "packageManager": "pnpm@9.0.0"
}
```

**pnpm-workspace.yaml:**
```yaml
packages:
  - "apps/*"
```

**Actions:**
1. Create root package.json and pnpm-workspace.yaml
2. Rename `frontend/` → `apps/web/`
3. Update `apps/web/package.json` name to `@optea/web`
4. Keep `apps/desktop/` as-is
5. Install root dependencies: `pnpm install`

### Git Repository

**Initialize:**
```bash
git init
git remote add origin git@github.com:AlexisTak/OpteaTech.git
```

**.gitignore:**
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
pnpm-debug.log*

# Editor
.DS_Store
*.pem
.vscode/settings.json
.idea/

# Go
api/bin/
api/vendor/
*.exe
*.test
*.out

# Tauri
apps/desktop/src-tauri/target/
```

**First commit:**
```bash
git add .
git commit -m "chore: initialize monorepo structure"
```

### Docker Compose

**docker-compose.yml:**
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

**.env (not committed):**
```env
POSTGRES_PASSWORD=optea_dev_password
```

### start-dev.sh Correction

**Updated script:**
```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_PORT="${API_PORT:-3001}"
WEB_PORT="${WEB_PORT:-3000}"
ADMIN_PORT="${ADMIN_PORT:-3002}"

echo "Starting OpteaTech services..."
echo "- Go API:    http://127.0.0.1:$API_PORT"
echo "- Frontend:  http://127.0.0.1:$WEB_PORT"
echo "- Admin:     http://127.0.0.1:$ADMIN_PORT"
echo ""

PIDS=()

cleanup() {
  echo "Stopping all services..."
  for pid in "${PIDS[@]:-}"; do
    kill "$pid" >/dev/null 2>&1 || true
  done
}

trap cleanup EXIT INT TERM

start_service() {
  local name="$1"
  local command="$2"
  (bash -lc "$command" 2>&1 | sed -u "s/^/[$name] /") &
  PIDS+=("$!")
}

start_service "api" "cd '$ROOT_DIR/api' && PORT='$API_PORT' go run ./cmd/server"
start_service "web" "cd '$ROOT_DIR/apps/web' && npm run dev -- --port '$WEB_PORT'"
start_service "admin" "cd '$ROOT_DIR/apps/admin' && npm run dev -- --port '$ADMIN_PORT'"

wait
```

**Why:** Removes non-existent `backend/` (NestJS BFF) reference, corrects paths for monorepo structure.

### Phase 1 Deliverables

✅ Monorepo structure functional with pnpm workspaces  
✅ Git repository initialized and connected to GitHub remote  
✅ Docker containers startable (`docker-compose up -d`)  
✅ start-dev.sh launches Go API + Web + Admin (admin doesn't exist yet, will be created in Phase 2)

### Phase 1 Validation

```bash
# Test monorepo
pnpm -r list  # Shows apps/web, apps/desktop

# Test docker
docker-compose up -d
docker ps  # Shows postgres + redis running

# Test start script (will fail on admin until Phase 2)
./start-dev.sh

# Test git
git log  # Shows initial commit
git remote -v  # Shows GitHub remote
```

## Phase 2: Dashboard Refine Basique (2-3h)

### Initialize Refine App

**Create `apps/admin/`:**
```bash
cd apps
npm create refine-app@latest admin -- --preset refine-vite
cd admin
```

**Update `apps/admin/package.json`:**
```json
{
  "name": "@optea/admin",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite --port 3002",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@refinedev/core": "^4.54.0",
    "@refinedev/react-router-v6": "^4.6.0",
    "@refinedev/simple-rest": "^5.0.8",
    "@refinedev/antd": "^5.43.0",
    "antd": "^5.21.0",
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "react-router-dom": "^6.26.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.12",
    "@types/react-dom": "^18.3.1",
    "@vitejs/plugin-react": "^4.3.3",
    "typescript": "^5.6.3",
    "vite": "^5.4.11"
  }
}
```

### Auth Provider Implementation

**File: `apps/admin/src/providers/authProvider.ts`**

```typescript
import { AuthProvider } from "@refinedev/core";

const API_URL = "http://localhost:3001/api";

export const authProvider: AuthProvider = {
  login: async ({ email, password }) => {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      return { success: false, error: { message: "Login failed", name: "Invalid credentials" } };
    }

    const data = await response.json();
    localStorage.setItem("auth", JSON.stringify(data));
    return { success: true, redirectTo: "/" };
  },

  logout: async () => {
    localStorage.removeItem("auth");
    return { success: true, redirectTo: "/login" };
  },

  check: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      return { authenticated: true };
    }
    return { authenticated: false, redirectTo: "/login" };
  },

  getPermissions: async () => null,

  getIdentity: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      const { user } = JSON.parse(auth);
      return user;
    }
    return null;
  },

  onError: async (error) => {
    if (error.statusCode === 401) {
      return { logout: true };
    }
    return { error };
  },
};
```

### Data Provider with JWT

**File: `apps/admin/src/providers/dataProvider.ts`**

```typescript
import dataProviderSimpleRest from "@refinedev/simple-rest";

const API_URL = "http://localhost:3001/api/admin";

const simpleRestProvider = dataProviderSimpleRest(API_URL);

export const dataProvider = {
  ...simpleRestProvider,
  custom: async ({ url, method, payload, headers }) => {
    const auth = localStorage.getItem("auth");
    const token = auth ? JSON.parse(auth).access_token : null;

    const response = await fetch(`${API_URL}${url}`, {
      method,
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...headers,
      },
      body: payload ? JSON.stringify(payload) : undefined,
    });

    return { data: await response.json() };
  },
};

// Axios interceptor for adding JWT to all requests
import axios from "axios";

axios.interceptors.request.use((config) => {
  const auth = localStorage.getItem("auth");
  if (auth) {
    const { access_token } = JSON.parse(auth);
    config.headers.Authorization = `Bearer ${access_token}`;
  }
  return config;
});
```

### Projects Resource (CRUD)

**File: `apps/admin/src/pages/projects/list.tsx`**

```typescript
import { List, useTable } from "@refinedev/antd";
import { Table, Space } from "antd";
import { EditButton, ShowButton, DeleteButton } from "@refinedev/antd";

export const ProjectList = () => {
  const { tableProps } = useTable({ resource: "projects" });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="title" title="Title" />
        <Table.Column dataIndex="slug" title="Slug" />
        <Table.Column dataIndex="category" title="Category" />
        <Table.Column dataIndex="status" title="Status" />
        <Table.Column
          title="Actions"
          render={(_, record) => (
            <Space>
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
              <DeleteButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
```

**File: `apps/admin/src/pages/projects/create.tsx`**

```typescript
import { Create, useForm } from "@refinedev/antd";
import { Form, Input, Select } from "antd";

export const ProjectCreate = () => {
  const { formProps, saveButtonProps } = useForm({ resource: "projects" });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item label="Title" name="title" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item label="Slug" name="slug" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item label="Category" name="category" rules={[{ required: true }]}>
          <Select>
            <Select.Option value="web">Web</Select.Option>
            <Select.Option value="logiciel">Logiciel</Select.Option>
            <Select.Option value="ia">IA</Select.Option>
            <Select.Option value="conseil">Conseil</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item label="Short Description" name="short_description">
          <Input.TextArea />
        </Form.Item>
        <Form.Item label="Full Description" name="full_description">
          <Input.TextArea rows={6} />
        </Form.Item>
      </Form>
    </Create>
  );
};
```

**File: `apps/admin/src/pages/projects/edit.tsx`** - Similar to create, with edit mode

**File: `apps/admin/src/pages/projects/show.tsx`** - Display project details

### Layout & Routing

**File: `apps/admin/src/App.tsx`**

```typescript
import { Refine } from "@refinedev/core";
import { RefineThemes, ThemedLayoutV2 } from "@refinedev/antd";
import routerBindings from "@refinedev/react-router-v6";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ConfigProvider } from "antd";

import { authProvider } from "./providers/authProvider";
import { dataProvider } from "./providers/dataProvider";
import { ProjectList, ProjectCreate, ProjectEdit, ProjectShow } from "./pages/projects";

function App() {
  return (
    <BrowserRouter>
      <ConfigProvider theme={RefineThemes.Blue}>
        <Refine
          authProvider={authProvider}
          dataProvider={dataProvider}
          routerProvider={routerBindings}
          resources={[
            {
              name: "projects",
              list: "/projects",
              create: "/projects/create",
              edit: "/projects/edit/:id",
              show: "/projects/show/:id",
              meta: { icon: "📁", label: "Projects" },
            },
          ]}
        >
          <Routes>
            <Route element={<ThemedLayoutV2 />}>
              <Route path="/projects" element={<ProjectList />} />
              <Route path="/projects/create" element={<ProjectCreate />} />
              <Route path="/projects/edit/:id" element={<ProjectEdit />} />
              <Route path="/projects/show/:id" element={<ProjectShow />} />
            </Route>
          </Routes>
        </Refine>
      </ConfigProvider>
    </BrowserRouter>
  );
}

export default App;
```

### Phase 2 Deliverables

✅ Refine app initialized in `apps/admin/`  
✅ Auth Provider connected to Go API `/api/auth/login`  
✅ Data Provider with JWT token injection  
✅ Projects resource with complete CRUD  
✅ Dashboard accessible at `http://localhost:3002`

### Phase 2 Validation

```bash
# Start all services
./start-dev.sh

# Visit admin dashboard
open http://localhost:3002

# Test auth flow
1. Login with Go API credentials
2. Should redirect to dashboard
3. Navigate to Projects
4. Test: List, Create, Edit, Delete operations
5. Verify data persists in PostgreSQL (via Go API)
```

## Phase 3: Dashboard Complet (3-4h)

### Additional Resources

**Services Resource:**
- List: Table with name, slug, is_active toggle
- Create/Edit: Form (name, slug, description, long_description, icon, color, features array, starting_price, order_index)
- Show: Display service details with features list

**Testimonials Resource:**
- List: Table with preview cards (client info, rating stars)
- Create/Edit: Form (client_name, client_role, client_company, client_avatar_url, content, rating 1-5, project_id relation, is_active, order_index)
- Show: Display testimonial with star rating visualization

**Messages Resource:**
- List: Read-only table with filters (read/unread status)
- Actions: Mark as read, Delete
- Show: Display full message with metadata (IP, user agent, timestamp)

**Client Requests Resource (Complex):**
- List: Table with status badges (pending, in_progress, completed, cancelled)
- Show: Detailed view with timeline of events
- Actions:
  - Update status
  - Update progress percentage
  - Create/update milestones
  - Add messages
  - Upload deliverables
  - Set/update quote
  - Manage access tokens (revoke, regenerate)

### Dashboard Page

**File: `apps/admin/src/pages/dashboard/index.tsx`**

```typescript
import { Row, Col, Card, Statistic } from "antd";
import { useCustom } from "@refinedev/core";

export const Dashboard = () => {
  const { data: stats } = useCustom({
    url: "/api/admin/dashboard",
    method: "get",
  });

  return (
    <div>
      <h1>Dashboard</h1>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic title="Total Projects" value={stats?.data?.total_projects || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="Unread Messages" value={stats?.data?.unread_messages || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="Active Requests" value={stats?.data?.active_requests || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="Testimonials" value={stats?.data?.total_testimonials || 0} />
          </Card>
        </Col>
      </Row>
    </div>
  );
};
```

### Final Navigation

**Updated `apps/admin/src/App.tsx` resources:**

```typescript
resources={[
  {
    name: "dashboard",
    list: "/",
    meta: { icon: "📊", label: "Dashboard" },
  },
  {
    name: "projects",
    list: "/projects",
    create: "/projects/create",
    edit: "/projects/edit/:id",
    show: "/projects/show/:id",
    meta: { icon: "📁", label: "Projects" },
  },
  {
    name: "services",
    list: "/services",
    create: "/services/create",
    edit: "/services/edit/:id",
    show: "/services/show/:id",
    meta: { icon: "🛠️", label: "Services" },
  },
  {
    name: "testimonials",
    list: "/testimonials",
    create: "/testimonials/create",
    edit: "/testimonials/edit/:id",
    show: "/testimonials/show/:id",
    meta: { icon: "💬", label: "Testimonials" },
  },
  {
    name: "messages",
    list: "/messages",
    show: "/messages/show/:id",
    meta: { icon: "✉️", label: "Messages" },
  },
  {
    name: "requests",
    list: "/requests",
    show: "/requests/show/:id",
    edit: "/requests/edit/:id",
    meta: { icon: "📋", label: "Client Requests" },
  },
]}
```

### Phase 3 Deliverables

✅ All resources (Projects, Services, Testimonials, Messages, Requests) with CRUD  
✅ Dashboard page with real-time statistics  
✅ Complete sidebar navigation  
✅ Production-ready admin panel

### Phase 3 Validation

```bash
# Test all resources
1. Navigate to each menu item
2. Test CRUD operations for each resource
3. Verify dashboard stats update correctly
4. Test complex operations (Requests: status update, milestones, messages)
5. Verify all data persists via Go API
```

## Data Flow & Error Handling

### Request Flow

```
User Action (Refine)
  ↓
Data Provider intercepts
  ↓
Auth Provider adds JWT header: "Authorization: Bearer {token}"
  ↓
HTTP Request → Go API endpoint
  ↓
Go Middleware validates JWT
  ↓
Repository Layer processes request
  ↓
PostgreSQL read/write
  ↓
JSON Response ← Go API
  ↓
Refine updates UI (optimistic updates)
```

### Error Handling Strategy

**401 Unauthorized (Token expired/invalid):**
- Auto-logout
- Redirect to login page
- Clear localStorage

**400 Bad Request (Validation errors):**
- API returns: `{ error: "message", details: {...} }`
- Refine displays in form fields (Ant Design Form.Item errors)

**500 Server Error:**
- Display Ant Design notification (red)
- Log to console in dev mode
- Do not crash UI

**Network Errors:**
- Retry once with exponential backoff
- If retry fails: notification "Connection lost"
- Allow manual retry

### CORS Configuration

**Why:** Go API must allow requests from admin dashboard origin.

**In API Go (`api/internal/middleware/cors.go`):**
```go
cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000", // Web frontend
        "http://localhost:3002", // Admin dashboard
    },
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
    AllowCredentials: true,
})
```

### Token Management

**Simple approach (Phase 2-3):**
- Store JWT in localStorage
- Add to all requests via axios interceptor
- On 401: logout immediately

**Optional enhancement (future):**
- Implement refresh token flow
- Interceptor detects 401, attempts `/api/auth/refresh` before logout
- More complex but better UX

## Testing & Validation Criteria

### Phase 1 Success Criteria

- [ ] `pnpm -r list` shows: `@optea/web`, `apps/desktop`
- [ ] `docker-compose up -d` starts PostgreSQL + Redis without errors
- [ ] `docker ps` shows both containers running
- [ ] `./start-dev.sh` launches Go API (port 3001) + Web (port 3000)
- [ ] `git log` shows initial commit with monorepo structure
- [ ] `git remote -v` shows GitHub remote `git@github.com:AlexisTak/OpteaTech.git`

### Phase 2 Success Criteria

- [ ] `http://localhost:3002` displays Refine login page
- [ ] Login with Go API credentials successfully redirects to dashboard
- [ ] Projects list displays data from Go API
- [ ] Create new project → visible in list and persisted in PostgreSQL
- [ ] Edit existing project → changes saved
- [ ] Delete project → removed from list
- [ ] Logout → redirects to login, clears auth token

### Phase 3 Success Criteria

- [ ] All resources accessible via sidebar navigation
- [ ] Services CRUD operations functional
- [ ] Testimonials CRUD with rating stars display correctly
- [ ] Messages list shows read/unread status
- [ ] Mark message as read updates status
- [ ] Client Requests list displays with status badges
- [ ] Request detail page shows timeline and all data
- [ ] Dashboard stats display correct numbers from API
- [ ] All operations persist data via Go API

### Global Success Criteria

- [ ] Monorepo structure clean and organized
- [ ] No regression on `apps/web` frontend functionality
- [ ] Docker containers run reliably
- [ ] start-dev.sh launches all services without errors
- [ ] Admin dashboard 100% functional for all CRUD operations
- [ ] Git history clean with logical commits per phase
- [ ] All endpoints communicate correctly with JWT auth
- [ ] Error handling graceful (no UI crashes)

## Risk Mitigation

### Potential Issues

**Risk 1: API Go endpoints mismatch with Refine expectations**
- **Mitigation:** Map API responses in Data Provider custom methods
- **Example:** Go returns `{data: [...], total: 10}`, ensure Data Provider returns same structure

**Risk 2: CORS errors blocking requests**
- **Mitigation:** Configure Go API CORS middleware to allow `localhost:3002`
- **Test early:** Use browser dev tools Network tab to verify headers

**Risk 3: JWT token expiry during long admin sessions**
- **Mitigation:** Implement token refresh flow (optional) or set longer expiry for admin tokens
- **Fallback:** User re-login (acceptable for admin dashboard)

**Risk 4: Monorepo structure breaks existing frontend**
- **Mitigation:** Test `apps/web` thoroughly after renaming from `frontend/`
- **Validation:** Ensure `next.config.ts` paths still resolve correctly

**Risk 5: Docker containers fail to start on some systems**
- **Mitigation:** Document port conflicts (5432, 6379), provide troubleshooting steps
- **Alternative:** Users can run PostgreSQL/Redis natively if Docker issues persist

## Future Enhancements

**Not in scope, but could be added later:**

1. **File Upload for Projects:**
   - Add image upload for project cover images
   - Store in S3 or local filesystem via Go API endpoint

2. **Rich Text Editor for Descriptions:**
   - Replace textarea with CKEditor or TinyMCE
   - Allow formatted content in project/service descriptions

3. **Real-time Notifications:**
   - WebSocket connection to Go API
   - Push notifications for new messages/requests

4. **Advanced Dashboard Analytics:**
   - Charts with Chart.js or Recharts
   - Time-series data for projects created, messages received

5. **Role-Based Access Control:**
   - Admin vs Editor roles
   - Restrict delete operations to admins only

6. **Audit Log:**
   - Track all admin actions (who changed what, when)
   - Display in dedicated audit page

## Appendix

### Key Files Summary

**Phase 1:**
- `package.json` - Root workspace config
- `pnpm-workspace.yaml` - Workspace definition
- `docker-compose.yml` - PostgreSQL + Redis
- `.gitignore` - Git exclusions
- `start-dev.sh` - Dev script (corrected)

**Phase 2:**
- `apps/admin/package.json` - Admin dashboard dependencies
- `apps/admin/src/providers/authProvider.ts` - Auth logic
- `apps/admin/src/providers/dataProvider.ts` - Data fetching with JWT
- `apps/admin/src/pages/projects/*` - Projects CRUD pages
- `apps/admin/src/App.tsx` - Main Refine setup

**Phase 3:**
- `apps/admin/src/pages/services/*` - Services CRUD
- `apps/admin/src/pages/testimonials/*` - Testimonials CRUD
- `apps/admin/src/pages/messages/*` - Messages read-only
- `apps/admin/src/pages/requests/*` - Client Requests complex CRUD
- `apps/admin/src/pages/dashboard/index.tsx` - Dashboard stats

### Tech Stack Versions

- **Node.js:** >=18.0.0
- **pnpm:** 9.0.0
- **Go:** 1.21+
- **PostgreSQL:** 16
- **Redis:** 7
- **Refine:** ^4.54.0
- **Ant Design:** ^5.21.0
- **React:** ^18.3.1
- **TypeScript:** ^5.6.3
- **Vite:** ^5.4.11

### Useful Commands

```bash
# Install all workspaces
pnpm install

# Start all services
./start-dev.sh

# Start specific services
pnpm dev:web
pnpm dev:admin
pnpm dev:api  # or: cd api && go run ./cmd/server

# Build all
pnpm build

# Docker commands
docker-compose up -d      # Start containers
docker-compose down       # Stop containers
docker-compose logs -f    # View logs

# Git commands
git status
git add .
git commit -m "message"
git push origin main
```

---

**End of Design Document**
