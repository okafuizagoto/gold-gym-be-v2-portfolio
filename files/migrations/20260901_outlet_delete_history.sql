-- History hapus outlet -- selama ini DeleteOutlet melakukan hard delete
-- (DELETE FROM outlet) tanpa jejak sama sekali. Simpan snapshot baris
-- outlet + siapa yang menghapus sebelum baris aslinya dihapus.
CREATE TABLE IF NOT EXISTS outlet_delete_history (
    history_id          INT AUTO_INCREMENT PRIMARY KEY,
    outlet_gold_id       INT NOT NULL,
    outlet_id            VARCHAR(36) NOT NULL,
    outlet_code          VARCHAR(20) NOT NULL,
    outlet_name          VARCHAR(255) NOT NULL,
    outlet_type          VARCHAR(50) NULL,
    outlet_address       VARCHAR(255) NULL,
    outlet_status        VARCHAR(20) NULL,
    outlet_created_at    DATETIME NULL,
    deleted_by           INT NOT NULL,
    deleted_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
