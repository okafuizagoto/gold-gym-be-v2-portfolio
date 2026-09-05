-- 2026-09-04: flag ADMIN untuk paksa sembunyikan menu "Daftar Pembeli" dan
-- "Mode Pembeli" dari akun penjual (retail & therapy), terlepas dari status
-- pendaftaran mereka sendiri (gold_buyer_yn). Dipakai layar admin baru
-- "Akses Daftar Pembeli" & "Akses Mode Pembeli" (mirip Outlet untuk Pembeli).
-- Default 'Y' supaya semua akun yang sudah ada tetap tampil menunya seperti
-- sebelum fitur ini ada.

ALTER TABLE data_peserta ADD COLUMN IF NOT EXISTS gold_menu_daftar_pembeli CHAR(1) NOT NULL DEFAULT 'Y' AFTER gold_buyer_yn;
ALTER TABLE data_peserta ADD COLUMN IF NOT EXISTS gold_menu_mode_pembeli CHAR(1) NOT NULL DEFAULT 'Y' AFTER gold_menu_daftar_pembeli;
