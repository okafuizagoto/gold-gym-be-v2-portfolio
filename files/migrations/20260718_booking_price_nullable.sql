-- 2026-07-18: fix error "booking_price cannot be null" saat insert booking UNPAID.
-- Backend (gorm) mengirim NULL eksplisit untuk booking belum bayar; kolom NOT NULL
-- menolaknya. Dibuat nullable + default 0 supaya backend versi lama pun aman.
-- (Backend baru selalu menulis 0 untuk booking belum bayar.)
-- SUDAH DIJALANKAN di DB remote pada 2026-07-18.

ALTER TABLE booking MODIFY booking_price DECIMAL(20,2) NULL DEFAULT 0.00;
