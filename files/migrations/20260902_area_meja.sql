-- Fitur "Atur Meja" (retail): area (indoor/outdoor) dan meja per outlet,
-- + snapshot meja yang dipakai per nota (sale_meja) supaya nota lama tetap
-- benar walau meja/area kelak berubah. Tanpa FOREIGN KEY, konsisten dengan
-- migrasi lain di project ini -- integritas relasi di level aplikasi saja.

CREATE TABLE IF NOT EXISTS area (
    area_id      INT AUTO_INCREMENT PRIMARY KEY,
    area_gold_id INT NOT NULL,
    area_outcode VARCHAR(20) NOT NULL,
    area_name    VARCHAR(100) NOT NULL,
    area_type    VARCHAR(10) NOT NULL,   -- 'INDOOR' | 'OUTDOOR'
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NULL ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_area_outcode (area_outcode)
);

CREATE TABLE IF NOT EXISTS meja (
    meja_id       INT AUTO_INCREMENT PRIMARY KEY,
    meja_gold_id  INT NOT NULL,
    meja_outcode  VARCHAR(20) NOT NULL,
    meja_area_id  INT NOT NULL,
    meja_name     VARCHAR(50) NOT NULL,
    meja_capacity INT NOT NULL DEFAULT 1,
    meja_status   VARCHAR(10) NOT NULL DEFAULT 'KOSONG',  -- 'KOSONG' | 'ISI'
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NULL ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_meja_outcode (meja_outcode),
    KEY idx_meja_area (meja_area_id),
    UNIQUE KEY uq_meja_outcode_name (meja_outcode, meja_name)
);

-- Snapshot meja yang dipakai per nota (mendukung gabungan >1 meja).
-- meja_name disalin apa adanya, bukan JOIN ke tabel meja, supaya nota lama
-- tetap benar walau baris meja kelak berubah.
CREATE TABLE IF NOT EXISTS sale_meja (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    sale_id    VARCHAR(64) NOT NULL,
    meja_id    INT NOT NULL,
    meja_name  VARCHAR(50) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_sale_meja_sale_id (sale_id)
);
