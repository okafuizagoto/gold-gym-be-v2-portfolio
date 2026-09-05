-- Nomor meja (gabungan, dipisah koma, mis. "A1, A3") ditempel langsung di
-- th_sale saat insert -- supaya Sales History (list & detail) bisa
-- menampilkannya tanpa join ke sale_meja per baris. NULL = tidak ada meja
-- dipilih untuk transaksi ini (mis. bukan outlet retail, atau kasir tidak
-- memilih meja).
ALTER TABLE th_sale ADD COLUMN IF NOT EXISTS sale_meja_names VARCHAR(255) NULL;
