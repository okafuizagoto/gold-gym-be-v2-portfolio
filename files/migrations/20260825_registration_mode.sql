-- Admin toggle untuk mode pendaftaran akun mandiri (layar "Daftar Akun",
-- register_screen.dart, tanpa OTP). Nilai: BOTH (default, boleh pilih
-- pembeli/penjual seperti sekarang), BUYER_ONLY, atau SELLER_ONLY.
--
-- Pakai tabel app_settings yang sudah ada (lihat
-- 20260819_payment_proof_visibility.sql) — tidak perlu ALTER TABLE, kolom
-- setting_value VARCHAR(16) sudah cukup untuk ketiga nilai di atas.
INSERT INTO app_settings (setting_key, setting_value)
VALUES ('registration_mode', 'BOTH')
ON DUPLICATE KEY UPDATE setting_key = setting_key;
