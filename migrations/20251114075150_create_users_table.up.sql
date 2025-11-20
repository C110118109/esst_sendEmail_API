-- 建立使用者表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

-- 建立索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- 新增註解
COMMENT ON TABLE users IS '使用者表';
COMMENT ON COLUMN users.id IS '使用者編號(UUID)';
COMMENT ON COLUMN users.username IS '使用者帳號';
COMMENT ON COLUMN users.password IS '密碼(已加密)';
COMMENT ON COLUMN users.role IS '角色 (admin/user)';
COMMENT ON COLUMN users.created_at IS '建立時間';
COMMENT ON COLUMN users.updated_at IS '更新時間';

-- 插入預設管理員帳號 (密碼: admin123)
-- bcrypt hash of "admin123"
INSERT INTO users (id, username, password, role, created_at) 
VALUES (
    uuid_generate_v4(),
    'admin',
    '$2a$10$8K1p/a0dL3LzF7sW5KvN6u5PN5x9lNy0CQvWXJzJ8FvVwqV7aCdoi',
    'admin',
    now()
);