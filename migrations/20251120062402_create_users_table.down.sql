-- 回滾 migration 檔案
-- 刪除使用者表

-- 刪除索引
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_username;

-- 刪除資料表
DROP TABLE IF EXISTS users;