# OpteaTech Monorepo + Admin Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform existing OpteaTech project into monorepo structure, initialize infrastructure (git, docker), and build complete Refine v4 admin dashboard with Ant Design consuming Go API endpoints.

**Architecture:** pnpm workspaces monorepo (apps/web, apps/admin, apps/desktop, api/), PostgreSQL + Redis via docker-compose, Refine dashboard with JWT auth consuming existing Go Fiber API admin endpoints.

**Tech Stack:** pnpm workspaces, Git, Docker Compose, Refine v4, Ant Design, React Router v6, TypeScript, Vite, Go 1.21, PostgreSQL 16, Redis 7.

---

## File Structure Overview

**Root level (new/modified):**
- Create: `package.json` - Root workspace configuration
- Create: `pnpm-workspace.yaml` - Workspace definition
- Create: `docker-compose.yml` - PostgreSQL + Redis services
- Create: `.gitignore` - Git exclusions
- Create: `.env` - Environment variables (not committed)
- Modify: `start-dev.sh` - Remove NestJS reference, add admin

**Monorepo structure changes:**
- Rename: `frontend/` → `apps/web/`
- Modify: `apps/web/package.json` - Update name to `@optea/web`
- Keep: `apps/desktop/` - No changes needed
- Keep: `api/` - Go API stays at root

**Admin dashboard (new):**
- Create: `apps/admin/` - Entire Refine application
- Create: `apps/admin/package.json` - Dependencies
- Create: `apps/admin/vite.config.ts` - Vite configuration
- Create: `apps/admin/tsconfig.json` - TypeScript config
- Create: `apps/admin/index.html` - Entry HTML
- Create: `apps/admin/src/main.tsx` - React entry point
- Create: `apps/admin/src/App.tsx` - Refine setup
- Create: `apps/admin/src/providers/authProvider.ts` - Auth logic
- Create: `apps/admin/src/providers/dataProvider.ts` - Data fetching
- Create: `apps/admin/src/pages/login/index.tsx` - Login page
- Create: `apps/admin/src/pages/dashboard/index.tsx` - Dashboard stats
- Create: `apps/admin/src/pages/projects/list.tsx` - Projects list
- Create: `apps/admin/src/pages/projects/create.tsx` - Create project
- Create: `apps/admin/src/pages/projects/edit.tsx` - Edit project
- Create: `apps/admin/src/pages/projects/show.tsx` - Show project
- Create: `apps/admin/src/pages/services/list.tsx` - Services list
- Create: `apps/admin/src/pages/services/create.tsx` - Create service
- Create: `apps/admin/src/pages/services/edit.tsx` - Edit service
- Create: `apps/admin/src/pages/services/show.tsx` - Show service
- Create: `apps/admin/src/pages/testimonials/list.tsx` - Testimonials list
- Create: `apps/admin/src/pages/testimonials/create.tsx` - Create testimonial
- Create: `apps/admin/src/pages/testimonials/edit.tsx` - Edit testimonial
- Create: `apps/admin/src/pages/testimonials/show.tsx` - Show testimonial
- Create: `apps/admin/src/pages/messages/list.tsx` - Messages list
- Create: `apps/admin/src/pages/messages/show.tsx` - Show message
- Create: `apps/admin/src/pages/requests/list.tsx` - Client requests list
- Create: `apps/admin/src/pages/requests/show.tsx` - Show request
- Create: `apps/admin/src/pages/requests/edit.tsx` - Edit request

---

## Phase 1: Infrastructure Setup

### Task 1: Create Root Workspace Configuration

**Files:**
- Create: `package.json`
- Create: `pnpm-workspace.yaml`

- [ ] **Step 1: Create root package.json**

```bash
cat > package.json << 'EOF'
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
    "build:admin": "pnpm --filter @optea/admin build",
    "lint": "pnpm lint:web && pnpm lint:admin",
    "lint:web": "pnpm --filter @optea/web lint",
    "lint:admin": "pnpm --filter @optea/admin lint"
  },
  "devDependencies": {
    "concurrently": "^8.2.2"
  },
  "packageManager": "pnpm@9.0.0"
}
EOF
```

- [ ] **Step 2: Create pnpm-workspace.yaml**

```bash
cat > pnpm-workspace.yaml << 'EOF'
packages:
  - "apps/*"
EOF
```

- [ ] **Step 3: Verify files created**

```bash
ls -la package.json pnpm-workspace.yaml
```

Expected: Both files exist

---

### Task 2: Restructure to Monorepo

**Files:**
- Rename: `frontend/` → `apps/web/`
- Modify: `apps/web/package.json`

- [ ] **Step 1: Create apps directory**

```bash
mkdir -p apps
```

- [ ] **Step 2: Move frontend to apps/web**

```bash
mv frontend apps/web
```

- [ ] **Step 3: Update apps/web/package.json name**

```bash
cd apps/web
sed -i 's/"name": "frontend"/"name": "@optea\/web"/' package.json
cd ../..
```

- [ ] **Step 4: Verify structure**

```bash
ls -la apps/
```

Expected: Shows `web/` and `desktop/` directories

- [ ] **Step 5: Test Next.js app still works**

```bash
cd apps/web && npm run dev
```

Expected: Server starts on port 3000 (Ctrl+C to stop)

---

### Task 3: Initialize Git Repository

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Create .gitignore**

```bash
cat > .gitignore << 'EOF'
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
pnpm-debug.log*

# Editor
.DS_Store
*.pem
.vscode/settings.json
.idea/
*.swp
*.swo

# Testing
coverage/
.turbo

# Go
api/bin/
api/vendor/
*.exe
*.test
*.out

# Tauri
apps/desktop/src-tauri/target/
apps/desktop/dist/
EOF
```

- [ ] **Step 2: Initialize git**

```bash
git init
```

Expected: "Initialized empty Git repository"

- [ ] **Step 3: Add remote**

```bash
git remote add origin git@github.com:AlexisTak/OpteaTech.git
```

- [ ] **Step 4: Verify remote**

```bash
git remote -v
```

Expected: Shows origin with GitHub URL

- [ ] **Step 5: Stage all files**

```bash
git add .
```

- [ ] **Step 6: Create initial commit**

```bash
git commit -m "chore: initialize monorepo structure with git"
```

Expected: Large commit with all existing files

---

### Task 4: Setup Docker Compose

**Files:**
- Create: `docker-compose.yml`
- Create: `.env`

- [ ] **Step 1: Create docker-compose.yml**

```bash
cat > docker-compose.yml << 'EOF'
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
EOF
```

- [ ] **Step 2: Create .env file**

```bash
cat > .env << 'EOF'
POSTGRES_PASSWORD=optea_dev_password
EOF
```

- [ ] **Step 3: Test Docker Compose**

```bash
docker-compose up -d
```

Expected: Downloads images and starts containers

- [ ] **Step 4: Verify containers running**

```bash
docker ps
```

Expected: Shows optea-postgres and optea-redis running

- [ ] **Step 5: Test PostgreSQL connection**

```bash
docker exec -it optea-postgres psql -U optea -d optea_tech -c '\dt'
```

Expected: Shows tables (if migrations ran) or "No relations found"

- [ ] **Step 6: Stop containers**

```bash
docker-compose down
```

- [ ] **Step 7: Commit docker files**

```bash
git add docker-compose.yml .gitignore
git commit -m "feat: add docker-compose for PostgreSQL and Redis"
```

Note: .env is not committed (in .gitignore)

---

### Task 5: Update start-dev.sh Script

**Files:**
- Modify: `start-dev.sh`

- [ ] **Step 1: Backup existing script**

```bash
cp start-dev.sh start-dev.sh.backup
```

- [ ] **Step 2: Replace start-dev.sh content**

```bash
cat > start-dev.sh << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_PORT="${API_PORT:-3001}"
WEB_PORT="${WEB_PORT:-3000}"
ADMIN_PORT="${ADMIN_PORT:-3002}"

echo ""
echo "Starting OpteaTech services..."
echo "- Go API:    http://127.0.0.1:$API_PORT"
echo "- Frontend:  http://127.0.0.1:$WEB_PORT"
echo "- Admin:     http://127.0.0.1:$ADMIN_PORT"
echo ""
echo "Press Ctrl+C to stop all services."
echo ""

PIDS=()

cleanup() {
  echo ""
  echo "Stopping all services..."
  for pid in "${PIDS[@]:-}"; do
    kill "$pid" >/dev/null 2>&1 || true
  done
  wait >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM

start_service() {
  local name="$1"
  local command="$2"
  (
    bash -lc "$command" 2>&1 | sed -u "s/^/[$name] /"
  ) &
  PIDS+=("$!")
}

start_service "api" "cd '$ROOT_DIR/api' && PORT='$API_PORT' go run ./cmd/server"
start_service "web" "cd '$ROOT_DIR/apps/web' && npm run dev -- --port '$WEB_PORT'"
start_service "admin" "cd '$ROOT_DIR/apps/admin' && npm run dev -- --port '$ADMIN_PORT'"

wait
EOF
```

- [ ] **Step 3: Make script executable**

```bash
chmod +x start-dev.sh
```

- [ ] **Step 4: Test script (will fail on admin until Phase 2)**

```bash
./start-dev.sh
```

Expected: API and Web start, Admin fails (not created yet). Press Ctrl+C.

- [ ] **Step 5: Commit updated script**

```bash
git add start-dev.sh
git commit -m "fix: update start-dev.sh for monorepo structure and remove NestJS"
```

---

### Task 6: Install Root Dependencies

**Files:**
- Modify: `pnpm-lock.yaml` (auto-generated)

- [ ] **Step 1: Install root dependencies**

```bash
pnpm install
```

Expected: Installs concurrently and creates pnpm-lock.yaml

- [ ] **Step 2: Verify installation**

```bash
pnpm list
```

Expected: Shows concurrently in root dependencies

- [ ] **Step 3: Commit lock file**

```bash
git add pnpm-lock.yaml
git commit -m "chore: install root workspace dependencies"
```

---

### Task 7: Phase 1 Validation

**Files:**
- None (validation only)

- [ ] **Step 1: Verify monorepo structure**

```bash
pnpm -r list
```

Expected: Lists @optea/web and apps/desktop workspaces

- [ ] **Step 2: Start Docker services**

```bash
docker-compose up -d
```

- [ ] **Step 3: Verify Docker containers**

```bash
docker ps | grep optea
```

Expected: Shows postgres and redis containers running

- [ ] **Step 4: Test API server**

```bash
cd api && PORT=3001 go run ./cmd/server &
sleep 3
curl http://localhost:3001/health
kill %1
cd ..
```

Expected: {"status":"ok","version":"1.0.0"}

- [ ] **Step 5: Create Phase 1 completion commit**

```bash
git add -A
git commit -m "chore: complete Phase 1 - infrastructure setup" --allow-empty
```

- [ ] **Step 6: Push to GitHub**

```bash
git branch -M main
git push -u origin main
```

Expected: Successfully pushes to GitHub

---

## Phase 2: Admin Dashboard Basique

### Task 8: Initialize Refine App Structure

**Files:**
- Create: `apps/admin/` directory
- Create: `apps/admin/package.json`
- Create: `apps/admin/vite.config.ts`
- Create: `apps/admin/tsconfig.json`
- Create: `apps/admin/index.html`

- [ ] **Step 1: Create admin directory**

```bash
mkdir -p apps/admin/src
```

- [ ] **Step 2: Create package.json**

```bash
cat > apps/admin/package.json << 'EOF'
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
    "axios": "^1.7.0",
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
EOF
```

- [ ] **Step 3: Create vite.config.ts**

```bash
cat > apps/admin/vite.config.ts << 'EOF'
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3002,
  },
});
EOF
```

- [ ] **Step 4: Create tsconfig.json**

```bash
cat > apps/admin/tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
EOF
```

- [ ] **Step 5: Create tsconfig.node.json**

```bash
cat > apps/admin/tsconfig.node.json << 'EOF'
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
EOF
```

- [ ] **Step 6: Create index.html**

```bash
cat > apps/admin/index.html << 'EOF'
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>OpteaTech Admin</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
EOF
```

- [ ] **Step 7: Install dependencies**

```bash
cd apps/admin && npm install
cd ../..
```

Expected: Installs all dependencies successfully

---

### Task 9: Create Auth Provider

**Files:**
- Create: `apps/admin/src/providers/authProvider.ts`

- [ ] **Step 1: Create providers directory**

```bash
mkdir -p apps/admin/src/providers
```

- [ ] **Step 2: Create authProvider.ts**

```bash
cat > apps/admin/src/providers/authProvider.ts << 'EOF'
import { AuthProvider } from "@refinedev/core";

const API_URL = "http://localhost:3001/api";

export const authProvider: AuthProvider = {
  login: async ({ email, password }) => {
    try {
      const response = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: "Login failed",
            name: "Invalid credentials",
          },
        };
      }

      const data = await response.json();
      localStorage.setItem("auth", JSON.stringify(data));

      return {
        success: true,
        redirectTo: "/",
      };
    } catch (error) {
      return {
        success: false,
        error: {
          message: "Network error",
          name: "Connection failed",
        },
      };
    }
  },

  logout: async () => {
    localStorage.removeItem("auth");
    return {
      success: true,
      redirectTo: "/login",
    };
  },

  check: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      return { authenticated: true };
    }
    return {
      authenticated: false,
      redirectTo: "/login",
    };
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
    if (error.statusCode === 401 || error.statusCode === 403) {
      return {
        logout: true,
        redirectTo: "/login",
      };
    }
    return { error };
  },
};
EOF
```

- [ ] **Step 3: Verify file created**

```bash
cat apps/admin/src/providers/authProvider.ts | head -20
```

Expected: Shows authProvider code

---

### Task 10: Create Data Provider with JWT

**Files:**
- Create: `apps/admin/src/providers/dataProvider.ts`

- [ ] **Step 1: Create dataProvider.ts**

```bash
cat > apps/admin/src/providers/dataProvider.ts << 'EOF'
import dataProviderSimpleRest from "@refinedev/simple-rest";
import axios from "axios";

const API_URL = "http://localhost:3001/api/admin";

// Configure axios interceptor to add JWT token
axios.interceptors.request.use((config) => {
  const auth = localStorage.getItem("auth");
  if (auth) {
    const { access_token } = JSON.parse(auth);
    config.headers.Authorization = `Bearer ${access_token}`;
  }
  return config;
});

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("auth");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export const dataProvider = dataProviderSimpleRest(API_URL, axios);
EOF
```

- [ ] **Step 2: Verify file created**

```bash
cat apps/admin/src/providers/dataProvider.ts | head -15
```

Expected: Shows dataProvider code with axios interceptor

---

### Task 11: Create Login Page

**Files:**
- Create: `apps/admin/src/pages/login/index.tsx`

- [ ] **Step 1: Create pages/login directory**

```bash
mkdir -p apps/admin/src/pages/login
```

- [ ] **Step 2: Create login page**

```bash
cat > apps/admin/src/pages/login/index.tsx << 'EOF'
import { useLogin } from "@refinedev/core";
import { Form, Input, Button, Card, Typography, Space } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";

const { Title } = Typography;

export const LoginPage = () => {
  const { mutate: login, isLoading } = useLogin();

  const onFinish = (values: { email: string; password: string }) => {
    login(values);
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        backgroundColor: "#f0f2f5",
      }}
    >
      <Card style={{ width: 400, boxShadow: "0 2px 8px rgba(0,0,0,0.1)" }}>
        <Space direction="vertical" size="large" style={{ width: "100%" }}>
          <Title level={2} style={{ textAlign: "center", margin: 0 }}>
            OpteaTech Admin
          </Title>
          <Form
            name="login"
            layout="vertical"
            onFinish={onFinish}
            autoComplete="off"
          >
            <Form.Item
              name="email"
              rules={[
                { required: true, message: "Please input your email!" },
                { type: "email", message: "Please enter a valid email!" },
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="Email"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: "Please input your password!" },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="Password"
                size="large"
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={isLoading}
                size="large"
                block
              >
                Log in
              </Button>
            </Form.Item>
          </Form>
        </Space>
      </Card>
    </div>
  );
};
EOF
```

- [ ] **Step 3: Verify file created**

```bash
ls -la apps/admin/src/pages/login/
```

Expected: Shows index.tsx

---

### Task 12: Create Projects List Page

**Files:**
- Create: `apps/admin/src/pages/projects/list.tsx`

- [ ] **Step 1: Create pages/projects directory**

```bash
mkdir -p apps/admin/src/pages/projects
```

- [ ] **Step 2: Create projects list page**

```bash
cat > apps/admin/src/pages/projects/list.tsx << 'EOF'
import { List, useTable, EditButton, ShowButton, DeleteButton } from "@refinedev/antd";
import { Table, Space, Tag } from "antd";

export const ProjectList = () => {
  const { tableProps } = useTable({
    resource: "projects",
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="title" title="Title" />
        <Table.Column dataIndex="slug" title="Slug" />
        <Table.Column
          dataIndex="category"
          title="Category"
          render={(value) => {
            const colors: Record<string, string> = {
              web: "blue",
              logiciel: "green",
              ia: "purple",
              conseil: "orange",
            };
            return <Tag color={colors[value] || "default"}>{value}</Tag>;
          }}
        />
        <Table.Column
          dataIndex="status"
          title="Status"
          render={(value) => {
            const colors: Record<string, string> = {
              draft: "default",
              published: "success",
              archived: "error",
            };
            return <Tag color={colors[value] || "default"}>{value}</Tag>;
          }}
        />
        <Table.Column
          dataIndex="featured"
          title="Featured"
          render={(value) => (value ? "Yes" : "No")}
        />
        <Table.Column
          title="Actions"
          render={(_, record: any) => (
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
EOF
```

- [ ] **Step 3: Verify file created**

```bash
cat apps/admin/src/pages/projects/list.tsx | head -20
```

Expected: Shows ProjectList component code

---

### Task 13: Create Projects Create Page

**Files:**
- Create: `apps/admin/src/pages/projects/create.tsx`

- [ ] **Step 1: Create projects create page**

```bash
cat > apps/admin/src/pages/projects/create.tsx << 'EOF'
import { Create, useForm } from "@refinedev/antd";
import { Form, Input, Select, Switch } from "antd";

const { TextArea } = Input;

export const ProjectCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "projects",
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Title"
          name="title"
          rules={[{ required: true, message: "Please enter title" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Slug"
          name="slug"
          rules={[{ required: true, message: "Please enter slug" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Category"
          name="category"
          rules={[{ required: true, message: "Please select category" }]}
        >
          <Select>
            <Select.Option value="web">Web</Select.Option>
            <Select.Option value="logiciel">Logiciel</Select.Option>
            <Select.Option value="ia">IA</Select.Option>
            <Select.Option value="conseil">Conseil</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Status"
          name="status"
          initialValue="draft"
        >
          <Select>
            <Select.Option value="draft">Draft</Select.Option>
            <Select.Option value="published">Published</Select.Option>
            <Select.Option value="archived">Archived</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="Short Description" name="short_description">
          <TextArea rows={3} />
        </Form.Item>

        <Form.Item label="Full Description" name="full_description">
          <TextArea rows={6} />
        </Form.Item>

        <Form.Item label="Cover Image URL" name="cover_image_url">
          <Input />
        </Form.Item>

        <Form.Item label="Project URL" name="project_url">
          <Input />
        </Form.Item>

        <Form.Item label="GitHub URL" name="github_url">
          <Input />
        </Form.Item>

        <Form.Item label="Client Name" name="client_name">
          <Input />
        </Form.Item>

        <Form.Item label="Featured" name="featured" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
EOF
```

- [ ] **Step 2: Verify file created**

```bash
ls -la apps/admin/src/pages/projects/
```

Expected: Shows list.tsx and create.tsx

---

### Task 14: Create Projects Edit Page

**Files:**
- Create: `apps/admin/src/pages/projects/edit.tsx`

- [ ] **Step 1: Create projects edit page**

```bash
cat > apps/admin/src/pages/projects/edit.tsx << 'EOF'
import { Edit, useForm } from "@refinedev/antd";
import { Form, Input, Select, Switch } from "antd";

const { TextArea } = Input;

export const ProjectEdit = () => {
  const { formProps, saveButtonProps, queryResult } = useForm({
    resource: "projects",
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Title"
          name="title"
          rules={[{ required: true, message: "Please enter title" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Slug"
          name="slug"
          rules={[{ required: true, message: "Please enter slug" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Category"
          name="category"
          rules={[{ required: true, message: "Please select category" }]}
        >
          <Select>
            <Select.Option value="web">Web</Select.Option>
            <Select.Option value="logiciel">Logiciel</Select.Option>
            <Select.Option value="ia">IA</Select.Option>
            <Select.Option value="conseil">Conseil</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="Status" name="status">
          <Select>
            <Select.Option value="draft">Draft</Select.Option>
            <Select.Option value="published">Published</Select.Option>
            <Select.Option value="archived">Archived</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="Short Description" name="short_description">
          <TextArea rows={3} />
        </Form.Item>

        <Form.Item label="Full Description" name="full_description">
          <TextArea rows={6} />
        </Form.Item>

        <Form.Item label="Cover Image URL" name="cover_image_url">
          <Input />
        </Form.Item>

        <Form.Item label="Project URL" name="project_url">
          <Input />
        </Form.Item>

        <Form.Item label="GitHub URL" name="github_url">
          <Input />
        </Form.Item>

        <Form.Item label="Client Name" name="client_name">
          <Input />
        </Form.Item>

        <Form.Item label="Featured" name="featured" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
EOF
```

- [ ] **Step 2: Verify file created**

```bash
ls -la apps/admin/src/pages/projects/ | wc -l
```

Expected: Shows 4 (. .. list.tsx create.tsx edit.tsx)

---

### Task 15: Create Projects Show Page

**Files:**
- Create: `apps/admin/src/pages/projects/show.tsx`

- [ ] **Step 1: Create projects show page**

```bash
cat > apps/admin/src/pages/projects/show.tsx << 'EOF'
import { Show } from "@refinedev/antd";
import { useShow } from "@refinedev/core";
import { Typography, Tag, Space } from "antd";

const { Title, Text } = Typography;

export const ProjectShow = () => {
  const { queryResult } = useShow({
    resource: "projects",
  });

  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Title</Title>
      <Text>{record?.title}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Slug
      </Title>
      <Text>{record?.slug}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Category
      </Title>
      <Tag>{record?.category}</Tag>

      <Title level={5} style={{ marginTop: 16 }}>
        Status
      </Title>
      <Tag>{record?.status}</Tag>

      <Title level={5} style={{ marginTop: 16 }}>
        Featured
      </Title>
      <Text>{record?.featured ? "Yes" : "No"}</Text>

      {record?.short_description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Short Description
          </Title>
          <Text>{record.short_description}</Text>
        </>
      )}

      {record?.full_description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Full Description
          </Title>
          <Text>{record.full_description}</Text>
        </>
      )}

      {record?.cover_image_url && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Cover Image
          </Title>
          <img
            src={record.cover_image_url}
            alt={record.title}
            style={{ maxWidth: "100%", maxHeight: 400 }}
          />
        </>
      )}

      {record?.client_name && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Client
          </Title>
          <Text>{record.client_name}</Text>
        </>
      )}

      {record?.project_url && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Project URL
          </Title>
          <a href={record.project_url} target="_blank" rel="noopener noreferrer">
            {record.project_url}
          </a>
        </>
      )}
    </Show>
  );
};
EOF
```

- [ ] **Step 2: Verify all project pages exist**

```bash
ls apps/admin/src/pages/projects/
```

Expected: Shows list.tsx, create.tsx, edit.tsx, show.tsx

---

### Task 16: Create Main App Component

**Files:**
- Create: `apps/admin/src/App.tsx`

- [ ] **Step 1: Create App.tsx**

```bash
cat > apps/admin/src/App.tsx << 'EOF'
import { Refine } from "@refinedev/core";
import { RefineThemes, ThemedLayoutV2, AuthPage } from "@refinedev/antd";
import routerBindings, {
  NavigateToResource,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router-v6";
import { BrowserRouter, Routes, Route, Outlet } from "react-router-dom";
import { ConfigProvider, App as AntdApp } from "antd";
import "@refinedev/antd/dist/reset.css";

import { authProvider } from "./providers/authProvider";
import { dataProvider } from "./providers/dataProvider";
import { LoginPage } from "./pages/login";
import { ProjectList } from "./pages/projects/list";
import { ProjectCreate } from "./pages/projects/create";
import { ProjectEdit } from "./pages/projects/edit";
import { ProjectShow } from "./pages/projects/show";

function App() {
  return (
    <BrowserRouter>
      <ConfigProvider theme={RefineThemes.Blue}>
        <AntdApp>
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
                meta: {
                  label: "Projects",
                },
              },
            ]}
            options={{
              syncWithLocation: true,
              warnWhenUnsavedChanges: true,
            }}
          >
            <Routes>
              <Route
                element={
                  <ThemedLayoutV2>
                    <Outlet />
                  </ThemedLayoutV2>
                }
              >
                <Route index element={<NavigateToResource resource="projects" />} />
                <Route path="/projects">
                  <Route index element={<ProjectList />} />
                  <Route path="create" element={<ProjectCreate />} />
                  <Route path="edit/:id" element={<ProjectEdit />} />
                  <Route path="show/:id" element={<ProjectShow />} />
                </Route>
              </Route>
              <Route path="/login" element={<LoginPage />} />
            </Routes>
            <UnsavedChangesNotifier />
            <DocumentTitleHandler />
          </Refine>
        </AntdApp>
      </ConfigProvider>
    </BrowserRouter>
  );
}

export default App;
EOF
```

- [ ] **Step 2: Verify file created**

```bash
cat apps/admin/src/App.tsx | head -30
```

Expected: Shows App component with Refine setup

---

### Task 17: Create Main Entry Point

**Files:**
- Create: `apps/admin/src/main.tsx`

- [ ] **Step 1: Create main.tsx**

```bash
cat > apps/admin/src/main.tsx << 'EOF'
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
EOF
```

- [ ] **Step 2: Verify file created**

```bash
cat apps/admin/src/main.tsx
```

Expected: Shows React entry point code

---

### Task 18: Test Admin Dashboard

**Files:**
- None (testing only)

- [ ] **Step 1: Start Docker services**

```bash
docker-compose up -d
```

- [ ] **Step 2: Start Go API in background**

```bash
cd api && PORT=3001 go run ./cmd/server &
API_PID=$!
cd ..
```

- [ ] **Step 3: Start admin dashboard**

```bash
cd apps/admin && npm run dev &
ADMIN_PID=$!
sleep 5
cd ../..
```

- [ ] **Step 4: Test admin dashboard loads**

```bash
curl -I http://localhost:3002
```

Expected: HTTP 200 OK

- [ ] **Step 5: Open browser (manual test)**

Open http://localhost:3002 in browser
Expected: See login page

- [ ] **Step 6: Stop services**

```bash
kill $API_PID $ADMIN_PID
```

---

### Task 19: Commit Phase 2

**Files:**
- All apps/admin/ files

- [ ] **Step 1: Stage all admin files**

```bash
git add apps/admin/
```

- [ ] **Step 2: Commit Phase 2**

```bash
git commit -m "feat: add Refine admin dashboard with Projects CRUD (Phase 2)"
```

- [ ] **Step 3: Push to GitHub**

```bash
git push origin main
```

Expected: Successfully pushes Phase 2

---

## Phase 3: Complete Admin Dashboard

### Task 20: Create Services Pages

**Files:**
- Create: `apps/admin/src/pages/services/list.tsx`
- Create: `apps/admin/src/pages/services/create.tsx`
- Create: `apps/admin/src/pages/services/edit.tsx`
- Create: `apps/admin/src/pages/services/show.tsx`

- [ ] **Step 1: Create services directory**

```bash
mkdir -p apps/admin/src/pages/services
```

- [ ] **Step 2: Create services list page**

```bash
cat > apps/admin/src/pages/services/list.tsx << 'EOF'
import { List, useTable, EditButton, ShowButton, DeleteButton } from "@refinedev/antd";
import { Table, Space, Tag, Switch } from "antd";

export const ServiceList = () => {
  const { tableProps } = useTable({
    resource: "services",
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="slug" title="Slug" />
        <Table.Column
          dataIndex="starting_price"
          title="Starting Price"
          render={(value) => (value ? `${value}€` : "-")}
        />
        <Table.Column
          dataIndex="is_active"
          title="Active"
          render={(value) => <Tag color={value ? "success" : "default"}>{value ? "Yes" : "No"}</Tag>}
        />
        <Table.Column
          dataIndex="order_index"
          title="Order"
        />
        <Table.Column
          title="Actions"
          render={(_, record: any) => (
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
EOF
```

- [ ] **Step 3: Create services create page**

```bash
cat > apps/admin/src/pages/services/create.tsx << 'EOF'
import { Create, useForm } from "@refinedev/antd";
import { Form, Input, InputNumber, Switch } from "antd";

const { TextArea } = Input;

export const ServiceCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "services",
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name="name"
          rules={[{ required: true, message: "Please enter name" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Slug"
          name="slug"
          rules={[{ required: true, message: "Please enter slug" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item label="Description" name="description">
          <TextArea rows={3} />
        </Form.Item>

        <Form.Item label="Long Description" name="long_description">
          <TextArea rows={6} />
        </Form.Item>

        <Form.Item label="Icon" name="icon">
          <Input placeholder="e.g., 🛠️" />
        </Form.Item>

        <Form.Item label="Color" name="color">
          <Input placeholder="e.g., #3b82f6" />
        </Form.Item>

        <Form.Item label="Starting Price (€)" name="starting_price">
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Order Index" name="order_index" initialValue={0}>
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Active" name="is_active" valuePropName="checked" initialValue={true}>
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
EOF
```

- [ ] **Step 4: Create services edit page**

```bash
cat > apps/admin/src/pages/services/edit.tsx << 'EOF'
import { Edit, useForm } from "@refinedev/antd";
import { Form, Input, InputNumber, Switch } from "antd";

const { TextArea } = Input;

export const ServiceEdit = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "services",
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name="name"
          rules={[{ required: true, message: "Please enter name" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Slug"
          name="slug"
          rules={[{ required: true, message: "Please enter slug" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item label="Description" name="description">
          <TextArea rows={3} />
        </Form.Item>

        <Form.Item label="Long Description" name="long_description">
          <TextArea rows={6} />
        </Form.Item>

        <Form.Item label="Icon" name="icon">
          <Input placeholder="e.g., 🛠️" />
        </Form.Item>

        <Form.Item label="Color" name="color">
          <Input placeholder="e.g., #3b82f6" />
        </Form.Item>

        <Form.Item label="Starting Price (€)" name="starting_price">
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Order Index" name="order_index">
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Active" name="is_active" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
EOF
```

- [ ] **Step 5: Create services show page**

```bash
cat > apps/admin/src/pages/services/show.tsx << 'EOF'
import { Show } from "@refinedev/antd";
import { useShow } from "@refinedev/core";
import { Typography, Tag } from "antd";

const { Title, Text } = Typography;

export const ServiceShow = () => {
  const { queryResult } = useShow({
    resource: "services",
  });

  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Name</Title>
      <Text>{record?.name}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Slug
      </Title>
      <Text>{record?.slug}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Icon
      </Title>
      <Text>{record?.icon || "-"}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Starting Price
      </Title>
      <Text>{record?.starting_price ? `${record.starting_price}€` : "-"}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Status
      </Title>
      <Tag color={record?.is_active ? "success" : "default"}>
        {record?.is_active ? "Active" : "Inactive"}
      </Tag>

      {record?.description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Description
          </Title>
          <Text>{record.description}</Text>
        </>
      )}

      {record?.long_description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Long Description
          </Title>
          <Text>{record.long_description}</Text>
        </>
      )}
    </Show>
  );
};
EOF
```

- [ ] **Step 6: Verify services pages created**

```bash
ls apps/admin/src/pages/services/
```

Expected: Shows list.tsx, create.tsx, edit.tsx, show.tsx

---

### Task 21: Create Testimonials Pages

**Files:**
- Create: `apps/admin/src/pages/testimonials/list.tsx`
- Create: `apps/admin/src/pages/testimonials/create.tsx`
- Create: `apps/admin/src/pages/testimonials/edit.tsx`
- Create: `apps/admin/src/pages/testimonials/show.tsx`

- [ ] **Step 1: Create testimonials directory**

```bash
mkdir -p apps/admin/src/pages/testimonials
```

- [ ] **Step 2: Create testimonials list page**

```bash
cat > apps/admin/src/pages/testimonials/list.tsx << 'EOF'
import { List, useTable, EditButton, ShowButton, DeleteButton } from "@refinedev/antd";
import { Table, Space, Tag, Rate } from "antd";

export const TestimonialList = () => {
  const { tableProps } = useTable({
    resource: "testimonials",
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="client_name" title="Client" />
        <Table.Column dataIndex="client_company" title="Company" />
        <Table.Column
          dataIndex="rating"
          title="Rating"
          render={(value) => <Rate disabled defaultValue={value} />}
        />
        <Table.Column
          dataIndex="is_active"
          title="Active"
          render={(value) => <Tag color={value ? "success" : "default"}>{value ? "Yes" : "No"}</Tag>}
        />
        <Table.Column
          title="Actions"
          render={(_, record: any) => (
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
EOF
```

- [ ] **Step 3: Create testimonials create page**

```bash
cat > apps/admin/src/pages/testimonials/create.tsx << 'EOF'
import { Create, useForm } from "@refinedev/antd";
import { Form, Input, Rate, Switch, InputNumber } from "antd";

const { TextArea } = Input;

export const TestimonialCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "testimonials",
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Client Name"
          name="client_name"
          rules={[{ required: true, message: "Please enter client name" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item label="Client Role" name="client_role">
          <Input />
        </Form.Item>

        <Form.Item label="Client Company" name="client_company">
          <Input />
        </Form.Item>

        <Form.Item label="Client Avatar URL" name="client_avatar_url">
          <Input />
        </Form.Item>

        <Form.Item
          label="Content"
          name="content"
          rules={[{ required: true, message: "Please enter content" }]}
        >
          <TextArea rows={4} />
        </Form.Item>

        <Form.Item label="Rating" name="rating" initialValue={5}>
          <Rate />
        </Form.Item>

        <Form.Item label="Order Index" name="order_index" initialValue={0}>
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Active" name="is_active" valuePropName="checked" initialValue={true}>
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
EOF
```

- [ ] **Step 4: Create testimonials edit page**

```bash
cat > apps/admin/src/pages/testimonials/edit.tsx << 'EOF'
import { Edit, useForm } from "@refinedev/antd";
import { Form, Input, Rate, Switch, InputNumber } from "antd";

const { TextArea } = Input;

export const TestimonialEdit = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "testimonials",
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Client Name"
          name="client_name"
          rules={[{ required: true, message: "Please enter client name" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item label="Client Role" name="client_role">
          <Input />
        </Form.Item>

        <Form.Item label="Client Company" name="client_company">
          <Input />
        </Form.Item>

        <Form.Item label="Client Avatar URL" name="client_avatar_url">
          <Input />
        </Form.Item>

        <Form.Item
          label="Content"
          name="content"
          rules={[{ required: true, message: "Please enter content" }]}
        >
          <TextArea rows={4} />
        </Form.Item>

        <Form.Item label="Rating" name="rating">
          <Rate />
        </Form.Item>

        <Form.Item label="Order Index" name="order_index">
          <InputNumber min={0} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item label="Active" name="is_active" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Edit>
  );
};
EOF
```

- [ ] **Step 5: Create testimonials show page**

```bash
cat > apps/admin/src/pages/testimonials/show.tsx << 'EOF'
import { Show } from "@refinedev/antd";
import { useShow } from "@refinedev/core";
import { Typography, Tag, Rate } from "antd";

const { Title, Text, Paragraph } = Typography;

export const TestimonialShow = () => {
  const { queryResult } = useShow({
    resource: "testimonials",
  });

  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Client Name</Title>
      <Text>{record?.client_name}</Text>

      {record?.client_role && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Role
          </Title>
          <Text>{record.client_role}</Text>
        </>
      )}

      {record?.client_company && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Company
          </Title>
          <Text>{record.client_company}</Text>
        </>
      )}

      <Title level={5} style={{ marginTop: 16 }}>
        Rating
      </Title>
      <Rate disabled defaultValue={record?.rating} />

      <Title level={5} style={{ marginTop: 16 }}>
        Content
      </Title>
      <Paragraph>{record?.content}</Paragraph>

      <Title level={5} style={{ marginTop: 16 }}>
        Status
      </Title>
      <Tag color={record?.is_active ? "success" : "default"}>
        {record?.is_active ? "Active" : "Inactive"}
      </Tag>
    </Show>
  );
};
EOF
```

- [ ] **Step 6: Verify testimonials pages created**

```bash
ls apps/admin/src/pages/testimonials/
```

Expected: Shows list.tsx, create.tsx, edit.tsx, show.tsx

---

### Task 22: Create Messages Pages

**Files:**
- Create: `apps/admin/src/pages/messages/list.tsx`
- Create: `apps/admin/src/pages/messages/show.tsx`

- [ ] **Step 1: Create messages directory**

```bash
mkdir -p apps/admin/src/pages/messages
```

- [ ] **Step 2: Create messages list page**

```bash
cat > apps/admin/src/pages/messages/list.tsx << 'EOF'
import { List, useTable, ShowButton, DeleteButton } from "@refinedev/antd";
import { Table, Space, Tag, Button } from "antd";
import { CheckOutlined } from "@ant-design/icons";
import { useUpdate } from "@refinedev/core";

export const MessageList = () => {
  const { tableProps } = useTable({
    resource: "messages",
    syncWithLocation: true,
  });

  const { mutate: updateMessage } = useUpdate();

  const handleMarkAsRead = (id: string) => {
    updateMessage({
      resource: "messages",
      id,
      values: { is_read: true },
      mutationMode: "optimistic",
    });
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="email" title="Email" />
        <Table.Column dataIndex="company" title="Company" />
        <Table.Column
          dataIndex="is_read"
          title="Status"
          render={(value) => (
            <Tag color={value ? "default" : "warning"}>
              {value ? "Read" : "Unread"}
            </Tag>
          )}
        />
        <Table.Column
          dataIndex="created_at"
          title="Date"
          render={(value) => new Date(value).toLocaleDateString()}
        />
        <Table.Column
          title="Actions"
          render={(_, record: any) => (
            <Space>
              {!record.is_read && (
                <Button
                  type="default"
                  size="small"
                  icon={<CheckOutlined />}
                  onClick={() => handleMarkAsRead(record.id)}
                >
                  Mark as Read
                </Button>
              )}
              <ShowButton hideText size="small" recordItemId={record.id} />
              <DeleteButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
EOF
```

- [ ] **Step 3: Create messages show page**

```bash
cat > apps/admin/src/pages/messages/show.tsx << 'EOF'
import { Show } from "@refinedev/antd";
import { useShow } from "@refinedev/core";
import { Typography, Tag, Descriptions } from "antd";

const { Title, Paragraph } = Typography;

export const MessageShow = () => {
  const { queryResult } = useShow({
    resource: "messages",
  });

  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Descriptions bordered column={1}>
        <Descriptions.Item label="Name">{record?.name}</Descriptions.Item>
        <Descriptions.Item label="Email">{record?.email}</Descriptions.Item>
        {record?.company && (
          <Descriptions.Item label="Company">{record.company}</Descriptions.Item>
        )}
        {record?.service_interest && (
          <Descriptions.Item label="Service Interest">
            {record.service_interest}
          </Descriptions.Item>
        )}
        {record?.budget_range && (
          <Descriptions.Item label="Budget Range">
            {record.budget_range}
          </Descriptions.Item>
        )}
        <Descriptions.Item label="Status">
          <Tag color={record?.is_read ? "default" : "warning"}>
            {record?.is_read ? "Read" : "Unread"}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Date">
          {record?.created_at && new Date(record.created_at).toLocaleString()}
        </Descriptions.Item>
        {record?.ip_address && (
          <Descriptions.Item label="IP Address">{record.ip_address}</Descriptions.Item>
        )}
      </Descriptions>

      <Title level={5} style={{ marginTop: 24 }}>
        Message
      </Title>
      <Paragraph>{record?.message}</Paragraph>
    </Show>
  );
};
EOF
```

- [ ] **Step 4: Verify messages pages created**

```bash
ls apps/admin/src/pages/messages/
```

Expected: Shows list.tsx, show.tsx

---

### Task 23: Create Dashboard Page

**Files:**
- Create: `apps/admin/src/pages/dashboard/index.tsx`

- [ ] **Step 1: Create dashboard directory**

```bash
mkdir -p apps/admin/src/pages/dashboard
```

- [ ] **Step 2: Create dashboard page**

```bash
cat > apps/admin/src/pages/dashboard/index.tsx << 'EOF'
import { Row, Col, Card, Statistic, List, Typography } from "antd";
import { useCustom, useList } from "@refinedev/core";
import {
  ProjectOutlined,
  MailOutlined,
  StarOutlined,
  ToolOutlined,
} from "@ant-design/icons";

const { Title } = Typography;

export const DashboardPage = () => {
  const { data: dashboardData } = useCustom({
    url: "http://localhost:3001/api/admin/dashboard",
    method: "get",
  });

  const { data: recentMessages } = useList({
    resource: "messages",
    pagination: { pageSize: 5 },
    sorters: [{ field: "created_at", order: "desc" }],
  });

  const stats = dashboardData?.data || {};

  return (
    <div style={{ padding: 24 }}>
      <Title level={2}>Dashboard</Title>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Projects"
              value={stats.total_projects || 0}
              prefix={<ProjectOutlined />}
            />
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Unread Messages"
              value={stats.unread_messages || 0}
              prefix={<MailOutlined />}
              valueStyle={{ color: "#cf1322" }}
            />
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Active Services"
              value={stats.active_services || 0}
              prefix={<ToolOutlined />}
            />
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Testimonials"
              value={stats.total_testimonials || 0}
              prefix={<StarOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="Recent Messages">
            <List
              dataSource={recentMessages?.data || []}
              renderItem={(item: any) => (
                <List.Item>
                  <List.Item.Meta
                    title={`${item.name} - ${item.email}`}
                    description={item.message?.substring(0, 100) + "..."}
                  />
                  <div>{new Date(item.created_at).toLocaleDateString()}</div>
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};
EOF
```

- [ ] **Step 3: Verify dashboard page created**

```bash
cat apps/admin/src/pages/dashboard/index.tsx | head -20
```

Expected: Shows DashboardPage component code

---

### Task 24: Update App.tsx with All Resources

**Files:**
- Modify: `apps/admin/src/App.tsx`

- [ ] **Step 1: Backup current App.tsx**

```bash
cp apps/admin/src/App.tsx apps/admin/src/App.tsx.backup
```

- [ ] **Step 2: Replace App.tsx with complete resources**

```bash
cat > apps/admin/src/App.tsx << 'EOF'
import { Refine } from "@refinedev/core";
import { RefineThemes, ThemedLayoutV2 } from "@refinedev/antd";
import routerBindings, {
  NavigateToResource,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router-v6";
import { BrowserRouter, Routes, Route, Outlet } from "react-router-dom";
import { ConfigProvider, App as AntdApp } from "antd";
import "@refinedev/antd/dist/reset.css";
import {
  DashboardOutlined,
  ProjectOutlined,
  ToolOutlined,
  StarOutlined,
  MailOutlined,
  FileTextOutlined,
} from "@ant-design/icons";

import { authProvider } from "./providers/authProvider";
import { dataProvider } from "./providers/dataProvider";
import { LoginPage } from "./pages/login";
import { DashboardPage } from "./pages/dashboard";
import { ProjectList, ProjectCreate, ProjectEdit, ProjectShow } from "./pages/projects";
import { ServiceList, ServiceCreate, ServiceEdit, ServiceShow } from "./pages/services";
import { TestimonialList, TestimonialCreate, TestimonialEdit, TestimonialShow } from "./pages/testimonials";
import { MessageList, MessageShow } from "./pages/messages";

function App() {
  return (
    <BrowserRouter>
      <ConfigProvider theme={RefineThemes.Blue}>
        <AntdApp>
          <Refine
            authProvider={authProvider}
            dataProvider={dataProvider}
            routerProvider={routerBindings}
            resources={[
              {
                name: "dashboard",
                list: "/",
                meta: {
                  label: "Dashboard",
                  icon: <DashboardOutlined />,
                },
              },
              {
                name: "projects",
                list: "/projects",
                create: "/projects/create",
                edit: "/projects/edit/:id",
                show: "/projects/show/:id",
                meta: {
                  label: "Projects",
                  icon: <ProjectOutlined />,
                },
              },
              {
                name: "services",
                list: "/services",
                create: "/services/create",
                edit: "/services/edit/:id",
                show: "/services/show/:id",
                meta: {
                  label: "Services",
                  icon: <ToolOutlined />,
                },
              },
              {
                name: "testimonials",
                list: "/testimonials",
                create: "/testimonials/create",
                edit: "/testimonials/edit/:id",
                show: "/testimonials/show/:id",
                meta: {
                  label: "Testimonials",
                  icon: <StarOutlined />,
                },
              },
              {
                name: "messages",
                list: "/messages",
                show: "/messages/show/:id",
                meta: {
                  label: "Messages",
                  icon: <MailOutlined />,
                },
              },
            ]}
            options={{
              syncWithLocation: true,
              warnWhenUnsavedChanges: true,
            }}
          >
            <Routes>
              <Route
                element={
                  <ThemedLayoutV2>
                    <Outlet />
                  </ThemedLayoutV2>
                }
              >
                <Route index element={<DashboardPage />} />
                
                <Route path="/projects">
                  <Route index element={<ProjectList />} />
                  <Route path="create" element={<ProjectCreate />} />
                  <Route path="edit/:id" element={<ProjectEdit />} />
                  <Route path="show/:id" element={<ProjectShow />} />
                </Route>

                <Route path="/services">
                  <Route index element={<ServiceList />} />
                  <Route path="create" element={<ServiceCreate />} />
                  <Route path="edit/:id" element={<ServiceEdit />} />
                  <Route path="show/:id" element={<ServiceShow />} />
                </Route>

                <Route path="/testimonials">
                  <Route index element={<TestimonialList />} />
                  <Route path="create" element={<TestimonialCreate />} />
                  <Route path="edit/:id" element={<TestimonialEdit />} />
                  <Route path="show/:id" element={<TestimonialShow />} />
                </Route>

                <Route path="/messages">
                  <Route index element={<MessageList />} />
                  <Route path="show/:id" element={<MessageShow />} />
                </Route>
              </Route>
              
              <Route path="/login" element={<LoginPage />} />
            </Routes>
            <UnsavedChangesNotifier />
            <DocumentTitleHandler />
          </Refine>
        </AntdApp>
      </ConfigProvider>
    </BrowserRouter>
  );
}

export default App;
EOF
```

- [ ] **Step 3: Create index exports for each resource**

```bash
# Projects index
cat > apps/admin/src/pages/projects/index.ts << 'EOF'
export { ProjectList } from "./list";
export { ProjectCreate } from "./create";
export { ProjectEdit } from "./edit";
export { ProjectShow } from "./show";
EOF

# Services index
cat > apps/admin/src/pages/services/index.ts << 'EOF'
export { ServiceList } from "./list";
export { ServiceCreate } from "./create";
export { ServiceEdit } from "./edit";
export { ServiceShow } from "./show";
EOF

# Testimonials index
cat > apps/admin/src/pages/testimonials/index.ts << 'EOF'
export { TestimonialList } from "./list";
export { TestimonialCreate } from "./create";
export { TestimonialEdit } from "./edit";
export { TestimonialShow } from "./show";
EOF

# Messages index
cat > apps/admin/src/pages/messages/index.ts << 'EOF'
export { MessageList } from "./list";
export { MessageShow } from "./show";
EOF
```

- [ ] **Step 4: Verify index files created**

```bash
find apps/admin/src/pages -name "index.ts"
```

Expected: Shows index.ts in projects, services, testimonials, messages

---

### Task 25: Final Testing

**Files:**
- None (testing only)

- [ ] **Step 1: Ensure Docker services running**

```bash
docker-compose up -d
```

- [ ] **Step 2: Start all services with script**

```bash
./start-dev.sh &
SCRIPT_PID=$!
sleep 10
```

- [ ] **Step 3: Test API health**

```bash
curl http://localhost:3001/health
```

Expected: {"status":"ok","version":"1.0.0"}

- [ ] **Step 4: Test Web frontend**

```bash
curl -I http://localhost:3000
```

Expected: HTTP 200 OK

- [ ] **Step 5: Test Admin dashboard**

```bash
curl -I http://localhost:3002
```

Expected: HTTP 200 OK

- [ ] **Step 6: Manual browser testing**

Open http://localhost:3002 in browser
Test workflow:
1. Login with Go API credentials
2. Navigate to Dashboard - verify stats
3. Navigate to Projects - test CRUD
4. Navigate to Services - test CRUD
5. Navigate to Testimonials - test CRUD
6. Navigate to Messages - test list and show
7. Logout

- [ ] **Step 7: Stop all services**

```bash
kill $SCRIPT_PID
docker-compose down
```

---

### Task 26: Final Commit and Push

**Files:**
- All Phase 3 files

- [ ] **Step 1: Stage all changes**

```bash
git add apps/admin/
```

- [ ] **Step 2: Commit Phase 3**

```bash
git commit -m "feat: complete admin dashboard with all resources (Phase 3)"
```

- [ ] **Step 3: Push to GitHub**

```bash
git push origin main
```

Expected: Successfully pushes complete implementation

- [ ] **Step 4: Verify GitHub repository**

Open https://github.com/AlexisTak/OpteaTech
Expected: All commits visible, complete monorepo structure

---

### Task 27: Update Documentation

**Files:**
- Create: `README.md` (root)

- [ ] **Step 1: Create comprehensive README.md**

```bash
cat > README.md << 'EOF'
# OpteaTech

Agence web tech — Sites web, logiciels sur mesure & solutions IA.

## Architecture

Monorepo structure with pnpm workspaces:

```
OpteaTech/
├── apps/
│   ├── web/          # Next.js 16 - Public website
│   ├── admin/        # Refine v4 - Admin dashboard
│   └── desktop/      # Tauri - Desktop app
├── api/              # Go 1.21 + Fiber v3 - REST API
└── docker-compose.yml # PostgreSQL + Redis
```

## Tech Stack

**Frontend (apps/web):**
- Next.js 16 with App Router
- TypeScript
- TailwindCSS
- React 19

**Admin Dashboard (apps/admin):**
- Refine v4
- Ant Design
- React Router v6
- TypeScript + Vite

**Backend (api/):**
- Go 1.21
- Fiber v3
- PostgreSQL (pgx)
- Redis
- JWT authentication

**Desktop (apps/desktop):**
- Tauri
- React
- TypeScript

## Development

### Prerequisites

- Node.js >=18.0.0
- pnpm 9.0.0
- Go 1.21+
- Docker & Docker Compose

### Quick Start

1. **Clone and install dependencies:**

```bash
git clone git@github.com:AlexisTak/OpteaTech.git
cd OpteaTech
pnpm install
```

2. **Start infrastructure:**

```bash
docker-compose up -d
```

3. **Start all services:**

```bash
./start-dev.sh
```

This starts:
- Go API: http://localhost:3001
- Web frontend: http://localhost:3000
- Admin dashboard: http://localhost:3002

### Individual Services

```bash
# API only
cd api && PORT=3001 go run ./cmd/server

# Web only
pnpm dev:web

# Admin only
pnpm dev:admin
```

### Build

```bash
# Build all
pnpm build

# Build specific
pnpm build:web
pnpm build:admin
```

## Admin Dashboard

Access: http://localhost:3002

Default credentials (change in production):
- Email: admin@optea.tech
- Password: (set via Go API)

Features:
- Projects management (CRUD)
- Services management (CRUD)
- Testimonials management (CRUD)
- Contact messages (read-only)
- Dashboard with statistics

## API Endpoints

**Public:**
- `GET /api/projects` - List published projects
- `GET /api/services` - List active services
- `GET /api/testimonials` - List active testimonials
- `POST /api/contact` - Submit contact message

**Admin (requires JWT):**
- `POST /api/auth/login` - Login
- `GET /api/admin/dashboard` - Dashboard stats
- `GET /api/admin/projects` - List all projects
- `POST /api/admin/projects` - Create project
- `PUT /api/admin/projects/:id` - Update project
- `DELETE /api/admin/projects/:id` - Delete project
- Similar endpoints for services, testimonials, messages

## Database

PostgreSQL schema includes:
- `projects` - Portfolio projects
- `services` - Services offered
- `testimonials` - Client testimonials
- `contact_messages` - Contact form submissions
- `admin_users` - Admin users

Migrations: `api/migrations/`

## Environment Variables

Create `.env` file:

```env
# Database
POSTGRES_PASSWORD=optea_dev_password
DATABASE_URL=postgresql://optea:optea_dev_password@localhost:5432/optea_tech

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=change-this-in-production
JWT_EXPIRES_IN=15

# API
PORT=3001
```

## Docker

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Reset database
docker-compose down -v
docker-compose up -d
```

## Contributing

1. Create feature branch
2. Make changes
3. Test thoroughly
4. Commit with conventional commits
5. Push and create PR

## License

Proprietary - OpteaTech
EOF
```

- [ ] **Step 2: Commit README**

```bash
git add README.md
git commit -m "docs: add comprehensive README with setup instructions"
```

- [ ] **Step 3: Push final changes**

```bash
git push origin main
```

---

## Validation Checklist

### Phase 1 Validation

- [ ] `pnpm -r list` shows workspaces
- [ ] Docker containers run (`docker ps`)
- [ ] Git repository connected to GitHub
- [ ] start-dev.sh executes without errors

### Phase 2 Validation

- [ ] Admin dashboard loads at localhost:3002
- [ ] Login works with Go API credentials
- [ ] Projects list displays data
- [ ] Projects CRUD operations functional
- [ ] JWT auth prevents unauthorized access

### Phase 3 Validation

- [ ] Dashboard shows statistics
- [ ] Services CRUD operational
- [ ] Testimonials CRUD operational
- [ ] Messages list and detail views work
- [ ] All sidebar navigation items accessible
- [ ] No regressions in Web frontend

### Global Validation

- [ ] Monorepo structure clean
- [ ] All dependencies installed correctly
- [ ] Docker services persistent
- [ ] Git history logical and clean
- [ ] Documentation complete
- [ ] No security issues (JWT, CORS, env vars)

---

**Implementation Complete!** 

Total Tasks: 27  
Estimated Time: ~7-9 hours  
Result: Production-ready monorepo with complete admin dashboard consuming Go API.
