-- 回滾 migration 檔案
-- 刪除現貨報備相關資料表

-- 刪除索引
DROP INDEX IF EXISTS idx_stock_equipments_created_time;
DROP INDEX IF EXISTS idx_stock_equipments_stock_id;
DROP INDEX IF EXISTS idx_stocks_expected_delivery_date;
DROP INDEX IF EXISTS idx_stocks_created_time;

-- 刪除資料表(按依賴順序,先刪除有外鍵的表)
DROP TABLE IF EXISTS stock_equipments;
DROP TABLE IF EXISTS stocks;
