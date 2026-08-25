-- 2026-07-25: Customer di POS.
-- 1) th_sale membawa sumber nama customer + flag tampil-di-nota (utk THERAPY).
--    sale_customer_source: PESERTA (nama toko gold_toko), CUSTOMER (tabel customer),
--                          MANUAL (ketikan kasir), BOOKING (dari nama booking terapi).
--    sale_customer_show:   Y/N — untuk outlet THERAPY, apakah nama customer dicetak
--                          di nota (toggle di POS/booking). RETAIL selalu Y.
-- 2) pos_customer_optional: daftar outlet RETAIL yang boleh transaksi POS TANPA
--    mengisi nama customer (diatur admin). Default: kosong = customer wajib diisi.
--
-- BELUM DIJALANKAN di DB remote.

ALTER TABLE th_sale
    ADD COLUMN sale_customer_source VARCHAR(10) NULL AFTER sale_salescustomer,
    ADD COLUMN sale_customer_show   CHAR(1) NOT NULL DEFAULT 'Y' AFTER sale_customer_source;

CREATE TABLE IF NOT EXISTS pos_customer_optional (
    pco_id       INT         NOT NULL AUTO_INCREMENT,
    pco_gold_id  INT         NOT NULL,             -- gold_id pemilik outlet (penjual)
    pco_outcode  VARCHAR(20) NOT NULL,
    pco_added_by VARCHAR(50) NULL,                 -- admin yang menambahkan
    pco_added_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (pco_id),
    UNIQUE KEY uq_pco (pco_gold_id, pco_outcode)
);
