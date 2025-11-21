-- 為使用者表新增 email 欄位
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS email TEXT;

-- 新增索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 新增註解
COMMENT ON COLUMN users.email IS '使用者信箱';