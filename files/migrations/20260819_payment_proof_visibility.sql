-- Admin toggle untuk fitur foto bukti pembayaran transfer (upload di POS/belanja
-- pembeli, tombol "lihat" di Sales History). Tiga gerbang independen; fitur
-- tampil HANYA jika ketiganya 'Y' (AND logic) — admin bisa menonaktifkan lewat
-- salah satu saja (global, outlet tertentu, atau satu akun tertentu).
--
-- app_settings: pengaturan global key-value. Baris payment_proof_enabled='Y'
-- adalah default (fitur aktif untuk semua, sampai admin menonaktifkan).
CREATE TABLE IF NOT EXISTS app_settings (
    setting_key   VARCHAR(64) NOT NULL PRIMARY KEY,
    setting_value VARCHAR(16) NOT NULL,
    updated_by    VARCHAR(100) NULL,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO app_settings (setting_key, setting_value)
VALUES ('payment_proof_enabled', 'Y')
ON DUPLICATE KEY UPDATE setting_key = setting_key;

-- outlet_proof_enabled: gerbang per outlet (RETAIL & THERAPY). Default 'Y'
-- supaya perilaku existing (fitur selalu tampil) tidak berubah sampai admin
-- menonaktifkan outlet tertentu.
ALTER TABLE outlet
    ADD COLUMN outlet_proof_enabled CHAR(1) NOT NULL DEFAULT 'Y';

-- gold_proof_enabled: gerbang per akun (penjual RETAIL/THERAPY maupun pembeli,
-- karena keduanya baris data_peserta). Default 'Y'.
ALTER TABLE data_peserta
    ADD COLUMN gold_proof_enabled CHAR(1) NOT NULL DEFAULT 'Y';
