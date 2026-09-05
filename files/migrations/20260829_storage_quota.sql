-- Kuota storage foto per user (semua role KECUALI ADMIN) -- lihat
-- internal/entity/quota untuk konstanta batasnya (30 MB).
ALTER TABLE data_peserta ADD COLUMN IF NOT EXISTS gold_storage_used_kb INT NOT NULL DEFAULT 0;

-- Ukuran file foto item (bytes) -- item_photo sebelumnya cuma simpan nama
-- file tanpa ukuran, dibutuhkan untuk tampilan MB di menu Storage & buat
-- decrement kuota yang akurat saat foto dihapus.
ALTER TABLE items ADD COLUMN IF NOT EXISTS item_photo_bytes INT NULL;

-- Pemilik numerik bukti pembayaran -- proof_uploaded_by selama ini cuma
-- string bebas (kebetulan selalu berisi gold_id sbg teks krn field "creator"
-- tak pernah diisi kode manapun), butuh kolom int + index utk query per-user
-- yang reliable (menu Storage & kuota).
ALTER TABLE payment_proof ADD COLUMN IF NOT EXISTS proof_gold_id INT NULL;
CREATE INDEX IF NOT EXISTS idx_proof_gold_id ON payment_proof (proof_gold_id);
UPDATE payment_proof SET proof_gold_id = CAST(proof_uploaded_by AS UNSIGNED)
    WHERE proof_gold_id IS NULL AND proof_uploaded_by REGEXP '^[0-9]+$';

-- Metode pembayaran transaksi POS -- belum pernah dilacak sama sekali di
-- th_sale sebelum ini (sale_transpayment = nominal, bukan metode). Dibutuhkan
-- supaya toast Sales History bisa membedakan TUNAI vs TRANSFER tanpa bukti.
ALTER TABLE th_sale ADD COLUMN IF NOT EXISTS sale_pay_type VARCHAR(20) NULL;

-- Backfill sekali: total KB dari bukti pembayaran yang sudah ada (item_photo
-- diabaikan di backfill krn fitur baru, praktis belum ada data lama).
-- Guard `= 0` supaya aman kalau migrasi ini pernah tak sengaja dijalankan ulang.
UPDATE data_peserta dp
LEFT JOIN (
    SELECT proof_gold_id, SUM(proof_bytes) AS total_bytes
    FROM payment_proof WHERE proof_gold_id IS NOT NULL GROUP BY proof_gold_id
) pp ON pp.proof_gold_id = dp.gold_id
SET dp.gold_storage_used_kb = CEIL(COALESCE(pp.total_bytes, 0) / 1024)
WHERE dp.gold_storage_used_kb = 0;

-- History hapus storage -- log saja, tidak ada tampilan sekarang.
CREATE TABLE IF NOT EXISTS storage_delete_history (
    history_id         INT AUTO_INCREMENT PRIMARY KEY,
    gold_id             INT NOT NULL,
    source_type         VARCHAR(20) NOT NULL,   -- 'ITEM_PHOTO' | 'PAYMENT_PROOF'
    source_id           INT NOT NULL,           -- item_id atau proof_id
    original_filename   VARCHAR(255) NOT NULL,
    file_bytes          INT NOT NULL,
    context_label       VARCHAR(255) NULL,      -- nama item / nomor nota
    deleted_by          VARCHAR(100) NULL,
    deleted_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
