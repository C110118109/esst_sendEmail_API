-- 專案報備系統 - 第二階段欄位 Migration
-- 執行此 SQL 以新增第二階段所需的欄位

-- 新增第二階段欄位到 projects 表
ALTER TABLE projects 
ADD COLUMN IF NOT EXISTS expected_delivery_period TEXT,
ADD COLUMN IF NOT EXISTS expected_delivery_date DATE,
ADD COLUMN IF NOT EXISTS expected_contract_period TEXT,
ADD COLUMN IF NOT EXISTS contract_start_date DATE,
ADD COLUMN IF NOT EXISTS contract_end_date DATE,
ADD COLUMN IF NOT EXISTS delivery_address TEXT,
ADD COLUMN IF NOT EXISTS special_requirements TEXT,
ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'step1',
ADD COLUMN IF NOT EXISTS updated_time TIMESTAMP;

-- 為現有資料設定預設狀態
UPDATE projects 
SET status = 'step1' 
WHERE status IS NULL OR status = '';

-- 建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_updated_time ON projects(updated_time);
CREATE INDEX IF NOT EXISTS idx_projects_expected_delivery_date ON projects(expected_delivery_date);

-- 新增註解說明
COMMENT ON COLUMN projects.expected_delivery_period IS '預計交貨期';
COMMENT ON COLUMN projects.expected_delivery_date IS '預計交貨日';
COMMENT ON COLUMN projects.expected_contract_period IS '預計履約期';
COMMENT ON COLUMN projects.contract_start_date IS '履約開始日';
COMMENT ON COLUMN projects.contract_end_date IS '履約結束日';
COMMENT ON COLUMN projects.delivery_address IS '交貨地址';
COMMENT ON COLUMN projects.special_requirements IS '特殊需求';
COMMENT ON COLUMN projects.status IS '專案狀態 (step1: 第一階段, step2: 第二階段, completed: 已完成)';
COMMENT ON COLUMN projects.updated_time IS '更新時間';