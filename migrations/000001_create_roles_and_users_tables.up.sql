CREATE TABLE roles (
  id UUID PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_by UUID,
  updated_by UUID,
  deleted_by UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_roles_slug_unique ON roles (slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_deleted_at ON roles (deleted_at);

CREATE TABLE users (
  id UUID PRIMARY KEY,
  role_id UUID NOT NULL REFERENCES roles(id),
  username VARCHAR(100) NOT NULL,
  password TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_by UUID,
  updated_by UUID,
  deleted_by UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_username_unique ON users (username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role_id ON users (role_id);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
