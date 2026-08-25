-- Perluasan modul Diskon (lanjutan 20260821_discount.sql):
-- 1. discount_scope: baris `discount` sekarang bisa 'ITEM' (per produk, sudah
--    ada) atau 'TOTAL' (persentase dari TOTAL penjualan, auto-aktif per
--    outlet -- untuk TOTAL, discount_item_id/discount_item_name diisi 0/''
--    sebagai sentinel "tidak spesifik ke item", bukan NULL, supaya kolom
--    existing tidak perlu diubah jadi nullable.
-- 2. discount_voucher: kode voucher sekali pakai (diketik manual atau
--    di-generate huruf besar+angka), diskon persen dari TOTAL keranjang.
--    Divalidasi & dikonsumsi (dihapus dari sini) saat dipakai di POS.
-- 3. discount_voucher_history: jejak audit voucher (dibuat/dipakai/dihapus).
-- 4. Kolom snapshot di th_sale: diskon total (auto, per outlet) dan voucher
--    (sekali pakai) yang berlaku untuk satu nota -- keduanya independen,
--    bisa menumpuk dengan diskon per-item di td_sale (lihat 20260821_discount.sql).

ALTER TABLE discount
    ADD COLUMN discount_scope VARCHAR(10) NOT NULL DEFAULT 'ITEM' AFTER discount_outcode;

CREATE TABLE IF NOT EXISTS discount_voucher (
    voucher_id         INT AUTO_INCREMENT PRIMARY KEY,
    voucher_gold_id    INT          NOT NULL,
    voucher_outcode    VARCHAR(20)  NOT NULL,
    voucher_code       VARCHAR(20)  NOT NULL,
    voucher_percent    DECIMAL(5,2) NOT NULL,
    voucher_expired_at DATETIME NULL,
    voucher_created_by VARCHAR(250) NULL,
    voucher_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_voucher_outcode_code (voucher_outcode, voucher_code),
    KEY idx_voucher_gold (voucher_gold_id)
);

CREATE TABLE IF NOT EXISTS discount_voucher_history (
    history_id                 INT AUTO_INCREMENT PRIMARY KEY,
    history_voucher_code       VARCHAR(20)  NOT NULL,
    history_outcode            VARCHAR(20)  NOT NULL,
    history_gold_id            INT          NOT NULL,
    history_percent            DECIMAL(5,2) NOT NULL,
    history_status             VARCHAR(10)  NOT NULL, -- USED / EXPIRED / DELETED
    history_sale_id             VARCHAR(64) NULL,
    history_actor_name         VARCHAR(150) NULL,
    history_actor_role         VARCHAR(20)  NULL,
    history_voucher_created_at DATETIME NULL,
    history_changed_at         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_vhist_outcode (history_outcode),
    KEY idx_vhist_code (history_voucher_code)
);

ALTER TABLE th_sale
    ADD COLUMN sale_total_discount_percent DECIMAL(5,2) NULL,
    ADD COLUMN sale_total_discount_amount  DECIMAL(15,2) NULL,
    ADD COLUMN sale_voucher_code           VARCHAR(20) NULL,
    ADD COLUMN sale_voucher_percent        DECIMAL(5,2) NULL,
    ADD COLUMN sale_voucher_amount         DECIMAL(15,2) NULL;
