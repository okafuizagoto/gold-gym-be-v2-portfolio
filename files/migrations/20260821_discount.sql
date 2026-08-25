-- Modul Diskon:
-- 1. Tabel discount: diskon per item per outlet (persen ATAU nominal tetap),
--    dibuat/diedit oleh penjual (SELLER/ADMIN) untuk outlet miliknya sendiri.
-- 2. Tabel discount_history: jejak audit setiap INSERT/UPDATE/DELETE diskon
--    (siapa, kapan, aksi apa, snapshot data diskon saat itu) — mengikuti pola
--    booking_remove_log tapi dengan kolom history_action karena di sini ada
--    3 aksi, bukan cuma delete.
-- 3. Kolom snapshot diskon di td_sale: saat item dengan diskon aktif masuk
--    nota, detail diskon (id/tipe/nilai/waktu diskon dibuat) DISALIN ke baris
--    td_sale supaya riwayat transaksi tidak berubah walau diskon aslinya
--    nanti diedit/dihapus.

CREATE TABLE IF NOT EXISTS discount (
    discount_id         INT AUTO_INCREMENT PRIMARY KEY,
    discount_gold_id    INT          NOT NULL,
    discount_outcode    VARCHAR(20)  NOT NULL,
    discount_item_id    INT          NOT NULL,
    discount_item_name  VARCHAR(150) NOT NULL,
    discount_type       VARCHAR(10)  NOT NULL, -- 'PERCENT' atau 'NOMINAL'
    discount_value      DECIMAL(15,2) NOT NULL, -- 0-100 kalau PERCENT, rupiah kalau NOMINAL
    discount_status     VARCHAR(10)  NOT NULL DEFAULT 'ACTIVE', -- ACTIVE/NONACTIVE
    discount_created_by VARCHAR(250) NULL,
    discount_created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    discount_updated_at DATETIME     NULL ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_discount_outcode (discount_outcode),
    KEY idx_discount_item (discount_item_id, discount_outcode),
    KEY idx_discount_gold (discount_gold_id)
);

CREATE TABLE IF NOT EXISTS discount_history (
    history_id                  INT AUTO_INCREMENT PRIMARY KEY,
    history_discount_id         INT          NOT NULL,
    history_action               VARCHAR(10)  NOT NULL, -- INSERT / UPDATE / DELETE
    history_gold_id              INT          NOT NULL,
    history_actor_name           VARCHAR(150) NULL,
    history_actor_role           VARCHAR(20)  NULL,
    history_outcode              VARCHAR(20)  NOT NULL,
    history_item_id              INT          NOT NULL,
    history_item_name            VARCHAR(150) NULL,
    history_discount_type        VARCHAR(10)  NULL,
    history_discount_value       DECIMAL(15,2) NULL,
    history_discount_status      VARCHAR(10)  NULL,
    history_discount_created_at  DATETIME NULL,
    history_changed_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_history_discount (history_discount_id),
    KEY idx_history_outcode (history_outcode)
);

-- Snapshot diskon per baris nota (NULL = baris ini tidak memakai diskon).
ALTER TABLE td_sale
    ADD COLUMN sale_discount_id         INT NULL,
    ADD COLUMN sale_discount_type       VARCHAR(10) NULL,
    ADD COLUMN sale_discount_value      DECIMAL(15,2) NULL,
    ADD COLUMN sale_discount_amount     DECIMAL(15,2) NULL,
    ADD COLUMN sale_discount_created_at DATETIME NULL;
