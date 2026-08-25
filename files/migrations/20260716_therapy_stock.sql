-- Item brand THERAPY (jasa) otomatis punya baris stock agar langsung tampil di
-- menu insert sales tanpa Add Stock. Qty 0 = tanpa batas (tidak divalidasi).
-- Backend menjalankan logika yang sama otomatis saat Add Item / Update Item
-- (EnsureTherapyStock); SQL ini hanya untuk item THERAPY yang sudah ada.

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
