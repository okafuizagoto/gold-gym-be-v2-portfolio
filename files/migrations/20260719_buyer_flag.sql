-- 2026-07-19: flag "sudah mendaftar sebagai pembeli" di akun (data_peserta).
-- Penjual mengkonfirmasi lewat menu Daftar Pembeli -> flag Y -> menu
-- "Mode Pembeli" muncul & menu Daftar Pembeli hilang di aplikasi.
-- Akun role BUYER otomatis Y (mereka memang pembeli).
-- SUDAH DIJALANKAN di DB remote pada 2026-07-19.

ALTER TABLE data_peserta ADD COLUMN gold_buyer_yn CHAR(1) NOT NULL DEFAULT 'N' AFTER gold_role;
UPDATE data_peserta SET gold_buyer_yn = 'Y' WHERE gold_role = 'BUYER';
