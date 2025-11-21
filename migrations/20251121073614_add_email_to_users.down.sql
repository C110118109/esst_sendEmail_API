-- 回滾 migration 檔案
-- 移除 email 欄位

-- 刪除索引
DROP INDEX IF EXISTS idx_users_email;

-- 移除欄位
ALTER TABLE users 
DROP COLUMN IF EXISTS email;