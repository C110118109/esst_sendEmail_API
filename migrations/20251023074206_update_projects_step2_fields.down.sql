-- 移除索引
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_updated_time;
DROP INDEX IF EXISTS idx_projects_expected_delivery_date;

-- 移除欄位
ALTER TABLE projects 
DROP COLUMN IF EXISTS expected_delivery_period,
DROP COLUMN IF EXISTS expected_delivery_date,
DROP COLUMN IF EXISTS expected_contract_period,
DROP COLUMN IF EXISTS contract_start_date,
DROP COLUMN IF EXISTS contract_end_date,
DROP COLUMN IF EXISTS delivery_address,
DROP COLUMN IF EXISTS special_requirements,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS updated_time;