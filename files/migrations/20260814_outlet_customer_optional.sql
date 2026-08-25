-- POS "boleh tanpa customer" kini jadi atribut per-outlet dengan DEFAULT boleh
-- kosong. Sebelumnya tersimpan di tabel pos_customer_optional (baris = boleh
-- kosong; default WAJIB). Keputusan baru (2026-08-14): default OPSIONAL (boleh
-- dikosongi); admin menandai outlet retail yang HARUS mengisi customer ('N').
--
-- outlet_customer_optional: 'Y' = boleh POS tanpa customer (default),
--                           'N' = wajib mengisi nama customer.
ALTER TABLE outlet
    ADD COLUMN outlet_customer_optional CHAR(1) NOT NULL DEFAULT 'Y';

-- Semua outlet lama otomatis 'Y' (opsional) sesuai default baru. Tidak ada data
-- dari pos_customer_optional yang perlu dipindah: outlet yang dulunya "boleh
-- kosong" tetap 'Y', dan outlet yang dulunya "wajib" (tanpa baris) kini ikut
-- jadi opsional sesuai keputusan default baru — admin dapat menandainya 'N' lagi.
-- Tabel pos_customer_optional dibiarkan (tidak dipakai lagi) untuk keamanan;
-- boleh di-DROP manual setelah dipastikan tidak diperlukan.
