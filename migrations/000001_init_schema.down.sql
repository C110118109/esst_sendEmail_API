-- 回滾 migration 檔案
-- 刪除索引
DROP INDEX IF EXISTS idx_equipments_created_time;
DROP INDEX IF EXISTS idx_projects_created_time;
DROP INDEX IF EXISTS idx_equipments_p_id;

-- 刪除資料表(按依賴順序,先刪除有外鍵的表)
DROP TABLE IF EXISTS equipments;
DROP TABLE IF EXISTS projects;

-- 移除 UUID 擴充套件(可選,如果其他地方不使用的話)
-- DROP EXTENSION IF EXISTS "uuid-ossp";