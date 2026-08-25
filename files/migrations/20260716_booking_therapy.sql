-- Migrasi fitur booking THERAPY + role pembeli
-- Dijalankan pada: u868654674_gold_gym_bez (MariaDB)

-- 1. Tipe outlet: RETAIL (default) atau THERAPY
ALTER TABLE outlet
  ADD COLUMN IF NOT EXISTS outlet_type VARCHAR(20) NOT NULL DEFAULT 'RETAIL' AFTER outlet_name;

-- 2. gold_id dibuat auto_increment agar registrasi pembeli tidak bentrok PK
ALTER TABLE data_peserta
  MODIFY gold_id INT(11) NOT NULL AUTO_INCREMENT;

-- 3. Tabel booking slot terapi
CREATE TABLE IF NOT EXISTS booking (
  booking_id VARCHAR(50) NOT NULL,
  booking_gold_id INT(11) NOT NULL,
  booking_outcode VARCHAR(20) NOT NULL,
  booking_date DATE NOT NULL,
  booking_start DATETIME NOT NULL,
  booking_duration INT(11) NOT NULL DEFAULT 60 COMMENT 'menit: 30 atau 60',
  booking_cust_id INT(11) DEFAULT NULL COMMENT 'NULL = pembeli belum terdaftar',
  booking_cust_name VARCHAR(100) NOT NULL,
  booking_registered_yn CHAR(1) NOT NULL DEFAULT 'N',
  booking_status VARCHAR(20) NOT NULL DEFAULT 'UNPAID' COMMENT 'PAID / UNPAID / EXPIRED / CANCELLED',
  booking_item_id INT(11) DEFAULT NULL,
  booking_item_name VARCHAR(100) DEFAULT NULL,
  booking_price DECIMAL(20,2) NOT NULL DEFAULT 0,
  booking_sale_id VARCHAR(50) DEFAULT NULL COMMENT 'link ke th_sale saat PAID',
  booking_created_by VARCHAR(250) DEFAULT NULL,
  booking_created_role VARCHAR(20) DEFAULT NULL COMMENT 'BUYER / SELLER',
  booking_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  booking_updated_at DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (booking_id),
  KEY idx_booking_outlet_date (booking_gold_id, booking_outcode, booking_date),
  KEY idx_booking_start (booking_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
