-- 2026-08-18: data_peserta.gold_otp tidak pernah kedaluwarsa -- OTP 6 digit
-- (1 juta kemungkinan) reset password bisa di-brute-force tanpa batas waktu
-- (temuan security review). Tambah kolom timestamp supaya OTP bisa
-- ditolak setelah 10 menit (dicek di service layer, UpdateDataPeserta).
ALTER TABLE data_peserta ADD COLUMN gold_otp_created_at DATETIME NULL;
