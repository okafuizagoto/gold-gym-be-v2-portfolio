-- 2026-08-14: seed data untuk akun ADMIN di DB PRODUCTION (VPS self-hosted).
-- Varian dari 20260718_seed_admin_outlets.sql, khusus production: data_peserta
-- production masih kosong saat migrasi ini dibuat, jadi user PERTAMA yang
-- register (akun admin) otomatis dapat gold_id=1 (AUTO_INCREMENT dari 1),
-- BUKAN gold_id=16 seperti di staging. Semua referensi gold_id diganti ke 1.
--
-- 2 outlet (ADM001 RETAIL "Admin Retail", ADM002 THERAPY "Admin Therapy"),
-- items + stock retail, 6 item jasa terapi + stock (qty 0 = jasa) di ADM002,
-- plus 1 item barang biasa (Aqua 600ml) di outlet therapy.
-- Akun admin harus didaftarkan lebih dulu lewat register API (jadi user
-- pertama di production) sebelum migrasi ini dijalankan, supaya gold_id=1
-- benar-benar terisi.

-- Outlet milik admin (gold_id 1)
INSERT INTO outlet (outlet_id, outlet_gold_id, outlet_code, outlet_name, outlet_type, outlet_address, outlet_status)
VALUES
  (UUID(), 1, 'ADM001', 'Admin Retail', 'RETAIL', 'Outlet retail admin', 'ACTIVE'),
  (UUID(), 1, 'ADM002', 'Admin Therapy', 'THERAPY', 'Outlet terapi admin', 'ACTIVE');

-- Items outlet RETAIL ADM001
INSERT INTO items (item_gold_id, item_outcode, item_code, item_name, item_type, item_pack, item_price, item_brand, item_status)
VALUES
  (1, 'ADM001', 'ITM000001', 'Aqua 600ml', 'STOCK', 'Pcs', 5000, 'AQUA', 'ACTIVE'),
  (1, 'ADM001', 'ITM000002', 'Gas 3kg', 'STOCK', 'Tabung', 25000, 'ELPIJI', 'ACTIVE');

-- Items outlet THERAPY ADM002: 6 jasa terapi + 1 barang biasa
INSERT INTO items (item_gold_id, item_outcode, item_code, item_name, item_type, item_pack, item_price, item_brand, item_status)
VALUES
  (1, 'ADM002', 'THR-SOFA30', 'SOFA 30 MENIT', 'STOCK', 'SESI', 15000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'THR-SOFA60', 'SOFA 1 JAM', 'STOCK', 'SESI', 25000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'THR-DRG30', 'KURSI DRAGON 30 MENIT', 'STOCK', 'SESI', 25000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'THR-DRG60', 'KURSI DRAGON 1 JAM', 'STOCK', 'SESI', 35000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'THR-KURSI30', 'KURSI 30 MENIT', 'STOCK', 'SESI', 10000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'THR-KURSI60', 'KURSI 1 JAM', 'STOCK', 'SESI', 20000, 'THERAPY', 'ACTIVE'),
  (1, 'ADM002', 'ITM000001', 'Aqua 600ml', 'STOCK', 'Pcs', 5000, 'AQUA', 'ACTIVE');

-- Stock semua item di atas (stock_id lanjut sequence STK%06d global).
-- Qty: jasa THERAPY = 0 (tak terbatas), Aqua = 100/50, Gas = 50.
INSERT INTO stock (stock_id, stock_gold_id, stock_outcode, stock_item_id, stock_name, stock_pack, stock_qty, stock_created_at, stock_update_by)
SELECT CONCAT('STK', LPAD(base.max_no + ROW_NUMBER() OVER (ORDER BY i.item_outcode, i.item_id), 6, '0')),
       i.item_gold_id, i.item_outcode, i.item_id, i.item_name, i.item_pack,
       CASE
         WHEN UPPER(i.item_brand) = 'THERAPY' THEN 0
         WHEN i.item_outcode = 'ADM001' AND i.item_code = 'ITM000001' THEN 100
         ELSE 50
       END,
       NOW(), 'SYSTEM'
FROM items i
CROSS JOIN (SELECT COALESCE(MAX(CAST(SUBSTRING(stock_id, 4) AS UNSIGNED)), 0) AS max_no FROM stock) base
WHERE i.item_gold_id = 1 AND i.item_outcode IN ('ADM001', 'ADM002');
