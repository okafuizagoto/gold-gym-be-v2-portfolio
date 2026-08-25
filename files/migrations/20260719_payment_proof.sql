-- 2026-07-19: bukti pembayaran transfer bank (upload foto dari menu POS).
-- File fisik disimpan backend di direktori PHOTO_STORAGE_DIR
-- (default /root/storages/photos); tabel ini menyimpan metadatanya.
-- Kuota global: SUM(proof_bytes) maksimal 10.000.000.000 bytes (10 GB) —
-- upload ditolak dengan pesan "tolong hubungi admin" jika kuota penuh.
-- SUDAH DIJALANKAN di DB remote pada 2026-07-19.

CREATE TABLE IF NOT EXISTS payment_proof (
  proof_id INT AUTO_INCREMENT PRIMARY KEY,
  proof_sale_id VARCHAR(50) NOT NULL COMMENT 'link ke th_sale.sale_id',
  proof_filename VARCHAR(255) NOT NULL COMMENT 'nama file tersimpan (uuid.ext)',
  proof_original_name VARCHAR(255) NULL COMMENT 'nama file asli dari user',
  proof_mime VARCHAR(100) NULL,
  proof_bytes BIGINT NOT NULL COMMENT 'ukuran file (bytes)',
  proof_path VARCHAR(255) NOT NULL COMMENT 'direktori penyimpanan di server',
  proof_uploaded_by VARCHAR(100) NULL,
  proof_uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_proof_sale (proof_sale_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
