CREATE TABLE IF NOT EXISTS client_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_name VARCHAR(100) NOT NULL,
  client_email VARCHAR(200) NOT NULL,
  client_company VARCHAR(100),
  client_phone VARCHAR(30),
  service_type VARCHAR(50) NOT NULL CHECK (service_type IN ('site_web', 'logiciel', 'ia', 'conseil', 'autre')),
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  budget_range VARCHAR(50),
  deadline DATE,
  attachments JSONB NOT NULL DEFAULT '[]'::jsonb,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  status VARCHAR(30) NOT NULL DEFAULT 'nouveau' CHECK (status IN ('nouveau', 'en_etude', 'devis_envoye', 'accepte', 'en_cours', 'en_revision', 'livre', 'termine', 'annule')),
  progress INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  quote_amount NUMERIC(10,2),
  quote_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
  quote_valid_until DATE,
  quote_accepted_at TIMESTAMPTZ,
  quote_pdf_url TEXT,
  access_token_hash VARCHAR(64) NOT NULL UNIQUE,
  token_created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  token_expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '90 days'),
  token_last_used TIMESTAMPTZ,
  token_use_count INTEGER NOT NULL DEFAULT 0,
  email_verified BOOLEAN NOT NULL DEFAULT false,
  email_sent_at TIMESTAMPTZ,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_milestones (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES client_requests(id) ON DELETE CASCADE,
  title VARCHAR(200) NOT NULL,
  description TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'done', 'blocked')),
  order_index INTEGER NOT NULL DEFAULT 0,
  due_date DATE,
  completed_at TIMESTAMPTZ,
  is_visible BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES client_requests(id) ON DELETE CASCADE,
  sender_type VARCHAR(10) NOT NULL CHECK (sender_type IN ('admin', 'client')),
  sender_name VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  attachments JSONB NOT NULL DEFAULT '[]'::jsonb,
  is_read BOOLEAN NOT NULL DEFAULT false,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_deliverables (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES client_requests(id) ON DELETE CASCADE,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  file_url TEXT NOT NULL,
  file_type VARCHAR(50) NOT NULL DEFAULT 'file',
  file_size BIGINT,
  version VARCHAR(20),
  is_visible BOOLEAN NOT NULL DEFAULT true,
  download_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client_access_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID REFERENCES client_requests(id) ON DELETE CASCADE,
  action VARCHAR(50) NOT NULL,
  ip_address VARCHAR(45),
  user_agent TEXT,
  country VARCHAR(2),
  success BOOLEAN NOT NULL DEFAULT true,
  failure_reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_requests_token_hash ON client_requests(access_token_hash);
CREATE INDEX IF NOT EXISTS idx_requests_email ON client_requests(client_email);
CREATE INDEX IF NOT EXISTS idx_requests_status ON client_requests(status);
CREATE INDEX IF NOT EXISTS idx_milestones_request ON project_milestones(request_id, order_index);
CREATE INDEX IF NOT EXISTS idx_messages_request ON project_messages(request_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_deliverables_request ON project_deliverables(request_id);
CREATE INDEX IF NOT EXISTS idx_access_logs_request ON client_access_logs(request_id, created_at DESC);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_client_requests_updated_at ON client_requests;
CREATE TRIGGER trg_client_requests_updated_at
  BEFORE UPDATE ON client_requests
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_project_milestones_updated_at ON project_milestones;
CREATE TRIGGER trg_project_milestones_updated_at
  BEFORE UPDATE ON project_milestones
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();