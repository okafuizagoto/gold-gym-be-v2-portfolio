-- Tipe terapi baru: KURSI (30 menit 10.000, 1 jam 20.000).
-- Seed item jasa per outlet THERAPY + baris stock (qty 0 = jasa tanpa batas),
-- pola sama dengan 20260717_booking_therapy_v2.sql.

INSERT INTO items (item_gold_id, item_outcode, item_code, item_name, item_type,
                   item_pack, item_price, item_brand, item_description, item_status)
SELECT o.outlet_gold_id, o.outlet_code, CONCAT('THR-', t.code), t.name, 'STOCK',
       'SESI', t.price, 'THERAPY', t.descr, 'ACTIVE'
FROM outlet o
CROSS JOIN (
    SELECT 'KURSI30' AS code, 'KURSI 30 MENIT' AS name, 10000 AS price, 'Terapi kursi 30 menit' AS descr
    UNION ALL SELECT 'KURSI60', 'KURSI 1 JAM', 20000, 'Terapi kursi 1 jam'
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
