-- 2026-08-15: data_peserta.gold_updated_by & gold_last_login_host NOT NULL
-- tanpa DEFAULT -> registrasi user baru (InsertGoldUser, struct GetGoldUsers)
-- tidak mengisi kolom ini sama sekali, gagal dengan error 1364 di sql_mode
-- strict (MariaDB self-hosted VPS). Di Hostinger lama sql_mode lebih longgar
-- jadi tidak pernah muncul. Kolom ini murni audit (siapa update terakhir /
-- host login terakhir), aman default string kosong untuk user baru.
ALTER TABLE data_peserta MODIFY COLUMN gold_updated_by VARCHAR(7) NOT NULL DEFAULT '';
ALTER TABLE data_peserta MODIFY COLUMN gold_last_login_host VARCHAR(255) NOT NULL DEFAULT '';
