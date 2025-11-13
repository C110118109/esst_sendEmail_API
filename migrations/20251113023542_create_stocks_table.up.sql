-- 現貨報備表 Migration
-- 建立 stocks 表

CREATE TABLE IF NOT EXISTS stocks (
    stock_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- 基本資訊
    stock_name TEXT NOT NULL,                  -- 現貨項目名稱(必填)
    contact_name TEXT NOT NULL,                -- 聯絡人姓名(必填)
    contact_phone TEXT,                        -- 聯絡電話(選填)
    contact_email TEXT,                        -- 聯絡信箱(選填)
    owner TEXT,                                -- 雙欣負責人(選填)
    
    -- 交貨資訊
    expected_delivery_period TEXT,             -- 預計交貨期
    expected_delivery_date DATE,               -- 預計交貨日
    expected_contract_period TEXT,             -- 預計履約期
    contract_start_date DATE,                  -- 履約開始日
    contract_end_date DATE,                    -- 履約結束日
    delivery_address TEXT,                     -- 交貨地址
    special_requirements TEXT,                 -- 特殊需求
    

    
    -- 備註
    remark TEXT,
    
    -- 時間戳記
    created_time TIMESTAMP NOT NULL DEFAULT now(),
    updated_time TIMESTAMP
);

-- 現貨設備表
CREATE TABLE IF NOT EXISTS stock_equipments (
    seq_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_id UUID NOT NULL,                    -- 現貨編號(必填,外鍵)
    part_number TEXT NOT NULL,                 -- 料號(必填)
    quantity INTEGER NOT NULL CHECK (quantity > 0), -- 數量(必填,必須大於0)
    description TEXT,                          -- 說明(選填)
    created_time TIMESTAMP NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_stock
        FOREIGN KEY (stock_id) 
        REFERENCES stocks(stock_id)
        ON DELETE CASCADE
);

-- 建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_stocks_created_time ON stocks(created_time DESC);
CREATE INDEX IF NOT EXISTS idx_stocks_expected_delivery_date ON stocks(expected_delivery_date);
CREATE INDEX IF NOT EXISTS idx_stock_equipments_stock_id ON stock_equipments(stock_id);
CREATE INDEX IF NOT EXISTS idx_stock_equipments_created_time ON stock_equipments(created_time DESC);

-- 新增註解說明
COMMENT ON TABLE stocks IS '現貨報備表';
COMMENT ON COLUMN stocks.stock_id IS '現貨編號(UUID)';
COMMENT ON COLUMN stocks.stock_name IS '現貨項目名稱';
COMMENT ON COLUMN stocks.contact_name IS '聯絡人姓名';
COMMENT ON COLUMN stocks.contact_phone IS '聯絡電話';
COMMENT ON COLUMN stocks.contact_email IS '聯絡信箱';
COMMENT ON COLUMN stocks.owner IS '雙欣負責人';
COMMENT ON COLUMN stocks.expected_delivery_period IS '預計交貨期';
COMMENT ON COLUMN stocks.expected_delivery_date IS '預計交貨日';
COMMENT ON COLUMN stocks.expected_contract_period IS '預計履約期';
COMMENT ON COLUMN stocks.contract_start_date IS '履約開始日';
COMMENT ON COLUMN stocks.contract_end_date IS '履約結束日';
COMMENT ON COLUMN stocks.delivery_address IS '交貨地址';
COMMENT ON COLUMN stocks.special_requirements IS '特殊需求';
COMMENT ON COLUMN stocks.remark IS '備註';
COMMENT ON COLUMN stocks.created_time IS '建立時間';
COMMENT ON COLUMN stocks.updated_time IS '更新時間';

COMMENT ON TABLE stock_equipments IS '現貨設備表';
COMMENT ON COLUMN stock_equipments.seq_id IS '現貨設備編號(UUID)';
COMMENT ON COLUMN stock_equipments.stock_id IS '現貨編號(外鍵)';
COMMENT ON COLUMN stock_equipments.part_number IS '料號';
COMMENT ON COLUMN stock_equipments.quantity IS '數量';
COMMENT ON COLUMN stock_equipments.description IS '設備說明';
COMMENT ON COLUMN stock_equipments.created_time IS '建立時間';
