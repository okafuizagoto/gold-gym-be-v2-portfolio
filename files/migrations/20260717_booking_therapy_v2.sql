-- Booking terapi v2:
-- 1. Kolom tipe terapi (SOFA / DRAGON) di booking
-- 2. Tabel log penghapusan booking (remove oleh penjual/admin)
-- 3. Kolom nama toko pembeli (pembeli yang didaftarkan penjual)
-- 4. Seed item jasa terapi per tipe & durasi di setiap outlet THERAPY
--    (SOFA 30 MENIT 15.000, SOFA 1 JAM 25.000,
--     KURSI DRAGON 30 MENIT 25.000, KURSI DRAGON 1 JAM 35.000)
-- 5. Baris stock untuk item THERAPY baru (qty 0 = jasa tanpa batas)

ALTER TABLE booking
    ADD COLUMN booking_therapy_type VARCHAR(20) NULL AFTER booking_duration;

CREATE TABLE IF NOT EXISTS booking_remove_log (
    log_id           INT AUTO_INCREMENT PRIMARY KEY,
    log_booking_id   VARCHAR(50)  NOT NULL,
    log_gold_id      INT          NOT NULL,
    log_outcode      VARCHAR(20)  NOT NULL,
    log_booking_date DATE         NULL,
    log_booking_start DATETIME    NULL,
    log_duration     INT          NULL,
    log_therapy_type VARCHAR(20)  NULL,
    log_cust_id      INT          NULL,
    log_cust_name    VARCHAR(100) NULL,
    log_status       VARCHAR(20)  NULL,
    log_removed_by   VARCHAR(250) NULL,
    log_removed_role VARCHAR(20)  NULL,
    log_removed_at   DATETIME     DEFAULT CURRENT_TIMESTAMP,
    KEY idx_remove_log_booking (log_booking_id),
    KEY idx_remove_log_outcode (log_outcode)
);

ALTER TABLE data_peserta
    ADD COLUMN gold_toko VARCHAR(100) NULL AFTER gold_nama;

INSERT INTO items (item_gold_id, item_outcode, item_code, item_name, item_type,
                   item_pack, item_price, item_brand, item_description, item_status)
SELECT o.outlet_gold_id, o.outlet_code, CONCAT('THR-', t.code), t.name, 'STOCK',
       'SESI', t.price, 'THERAPY', t.descr, 'ACTIVE'
FROM outlet o
CROSS JOIN (
    SELECT 'SOFA30' AS code, 'SOFA 30 MENIT' AS name, 15000 AS price, 'Terapi sofa 30 menit' AS descr
    UNION ALL SELECT 'SOFA60', 'SOFA 1 JAM', 25000, 'Terapi sofa 1 jam'
    UNION ALL SELECT 'DRG30', 'KURSI DRAGON 30 MENIT', 25000, 'Terapi kursi dragon 30 menit'
    UNION ALL SELECT 'DRG60', 'KURSI DRAGON 1 JAM', 35000, 'Terapi kursi dragon 1 jam'
) t
LEFT JOIN items i ON i.item_outcode = o.outlet_code
    AND UPPER(i.item_name) = t.name AND UPPER(i.item_brand) = 'THERAPY'
WHERE o.outlet_type = 'THERAPY' AND i.item_id IS NULL;

INSERT INTO stock (stock_id, stock_gold_id, stock_outcode, stock_item_id, stock_name, stock_pack, stock_qty, stock_created_at, stock_update_by)
SELECT CONCAT('STK', LPAD(base.maxnum + ROW_NUMBER() OVER (ORDER BY i.item_id), 6, '0')),
       i.item_gold_id, i.item_outcode, i.item_id, i.item_name, i.item_pack, 0, NOW(), 'SYSTEM'
FROM items i
CROSS JOIN (
    SELECT COALESCE(MAX(CAST(SUBSTRING(stock_id, 4) AS UNSIGNED)), 0) AS maxnum
    FROM stock WHERE stock_id LIKE 'STK%'
) base
LEFT JOIN stock s ON s.stock_item_id = i.item_id AND s.stock_outcode = i.item_outcode AND s.stock_gold_id = i.item_gold_id
WHERE UPPER(i.item_brand) = 'THERAPY' AND s.stock_id IS NULL;
