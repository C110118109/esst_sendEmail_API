-- 啟用 UUID 擴充套件
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 專案表
CREATE TABLE IF NOT EXISTS projects (
    p_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    p_name TEXT NOT NULL,                      -- 專案名稱(必填)
    contact_name TEXT NOT NULL,                -- 聯絡人姓名(必填)
    contact_phone TEXT,                        -- 聯絡電話(選填) ← 修正
    contact_email TEXT,                        -- 聯絡信箱(選填) ← 修正
    owner TEXT,                                -- 負責人(選填) ← 修正
    remark TEXT,                               -- 備註(選填)
    created_time TIMESTAMP NOT NULL DEFAULT now()
);

-- 設備表
CREATE TABLE IF NOT EXISTS equipments (
    eq_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    p_id UUID NOT NULL,                        -- 專案編號(必填,外鍵)
    part_number TEXT NOT NULL,                 -- 料號(必填)
    quantity INTEGER NOT NULL CHECK (quantity > 0), -- 數量(必填,必須大於0) ← 新增檢查
    description TEXT,                          -- 說明(選填)
    created_time TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_project                      -- 外鍵約束
        FOREIGN KEY (p_id) 
        REFERENCES projects(p_id)
        ON DELETE CASCADE                      -- 刪除專案時同時刪除設備 ← 新增
);

-- 建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_equipments_p_id ON equipments(p_id);
CREATE INDEX IF NOT EXISTS idx_projects_created_time ON projects(created_time DESC);
CREATE INDEX IF NOT EXISTS idx_equipments_created_time ON equipments(created_time DESC);

-- 新增註解說明
COMMENT ON TABLE projects IS '專案表';
COMMENT ON COLUMN projects.p_id IS '專案編號(UUID)';
COMMENT ON COLUMN projects.p_name IS '專案名稱';
COMMENT ON COLUMN projects.contact_name IS '聯絡人姓名';
COMMENT ON COLUMN projects.contact_phone IS '聯絡電話';
COMMENT ON COLUMN projects.contact_email IS '聯絡信箱';
COMMENT ON COLUMN projects.owner IS '負責人';
COMMENT ON COLUMN projects.remark IS '備註';
COMMENT ON COLUMN projects.created_time IS '建立時間';

COMMENT ON TABLE equipments IS '設備表';
COMMENT ON COLUMN equipments.eq_id IS '設備編號(UUID)';
COMMENT ON COLUMN equipments.p_id IS '專案編號(外鍵)';
COMMENT ON COLUMN equipments.part_number IS '料號';
COMMENT ON COLUMN equipments.quantity IS '數量';
COMMENT ON COLUMN equipments.description IS '設備說明';
COMMENT ON COLUMN equipments.created_time IS '建立時間';