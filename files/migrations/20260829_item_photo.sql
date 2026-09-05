-- Foto item (opsional, maks 2MB divalidasi di backend) untuk kartu produk di
-- POS. Kolom cuma menyimpan nama file -- file fisik ada di VPS
-- /root/storages/photos (folder yang sama dipakai foto bukti pembayaran,
-- dibedakan lewat prefix nama file "item_").
ALTER TABLE items ADD COLUMN IF NOT EXISTS item_photo VARCHAR(255) NULL;
