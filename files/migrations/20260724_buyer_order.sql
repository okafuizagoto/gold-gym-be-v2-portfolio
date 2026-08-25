-- 2026-07-24: Alur pesanan pembeli (BUYER order flow).
-- Pembeli (role BUYER, atau penjual/admin dalam "Mode Pembeli") memesan barang
-- dari outlet RETAIL milik penjual lain. Pesanan TIDAK langsung jadi nota:
-- pesanan masuk tabel buyer_order dengan status sendiri, dan baru masuk th_sale
-- (jadi nota/omzet penjual) setelah status FINISH.
--
-- Status:
--   WAITING  = menunggu konfirmasi penjual
--   PROCESS  = penjual sudah konfirmasi (in progress)
--   FINISH   = penjual "selesai proses" -> nota dibuat (masuk th_sale)
--   REJECT   = penjual menolak (dengan alasan di order_reject_reason)
--
-- Pembayaran:
--   TUNAI    = order_paid_yn 'N' (belum lunas; penjual tandai lunas di Sales History)
--   TRANSFER = order_paid_yn 'Y' (lunas; pembeli upload bukti transfer)
-- Bukti transfer memakai tabel payment_proof yang sudah ada
-- (proof_sale_id = order_id).
--
-- BELUM DIJALANKAN di DB remote (jalankan manual saat siap deploy).

CREATE TABLE IF NOT EXISTS buyer_order (
    order_id            VARCHAR(64)  NOT NULL,
    order_buyer_id      INT          NOT NULL,               -- gold_id pemesan (data_peserta)
    order_buyer_name    VARCHAR(100) NOT NULL DEFAULT '',
    order_gold_id       INT          NOT NULL,               -- gold_id pemilik outlet/penjual
    order_outcode       VARCHAR(20)  NOT NULL,
    order_outlet_name   VARCHAR(100) NOT NULL DEFAULT '',
    order_total         DECIMAL(15,2) NOT NULL DEFAULT 0,
    order_pay_type      VARCHAR(10)  NOT NULL,               -- TUNAI / TRANSFER
    order_paid_yn       CHAR(1)      NOT NULL DEFAULT 'N',   -- Y jika sudah lunas
    order_status        VARCHAR(10)  NOT NULL DEFAULT 'WAITING',
    order_reject_reason VARCHAR(255) NULL,
    order_sale_id       VARCHAR(64)  NULL,                   -- diisi saat FINISH (relasi th_sale)
    order_created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_updated_at    DATETIME     NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (order_id),
    KEY idx_bo_buyer  (order_buyer_id),
    KEY idx_bo_seller (order_gold_id, order_status),
    KEY idx_bo_outlet (order_gold_id, order_outcode)
);

CREATE TABLE IF NOT EXISTS buyer_order_detail (
    od_id         INT           NOT NULL AUTO_INCREMENT,
    od_order_id   VARCHAR(64)   NOT NULL,
    od_stock_id   VARCHAR(64)   NOT NULL,
    od_stock_name VARCHAR(150)  NOT NULL DEFAULT '',
    od_qty        INT           NOT NULL DEFAULT 0,
    od_price      DECIMAL(15,2) NOT NULL DEFAULT 0,
    od_total      DECIMAL(15,2) NOT NULL DEFAULT 0,
    od_pack       VARCHAR(20)   NULL,
    PRIMARY KEY (od_id),
    KEY idx_bod_order (od_order_id)
);

-- Kurasi outlet oleh ADMIN: outlet penjual mana yang boleh dilihat & dipesan
-- oleh role pembeli. Pembeli HANYA melihat outlet yang tercatat di sini
-- (default: kosong = pembeli tidak melihat outlet apa pun sampai admin
-- menambahkan). outlet_code tidak unik global -> wajib pasangan gold_id+outcode.
CREATE TABLE IF NOT EXISTS buyer_visible_outlet (
    bvo_id       INT         NOT NULL AUTO_INCREMENT,
    bvo_gold_id  INT         NOT NULL,             -- gold_id pemilik outlet
    bvo_outcode  VARCHAR(20) NOT NULL,
    bvo_added_by VARCHAR(50) NULL,                 -- admin yang menambahkan
    bvo_added_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (bvo_id),
    UNIQUE KEY uq_bvo (bvo_gold_id, bvo_outcode)
);
