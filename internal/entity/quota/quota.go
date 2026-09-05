package quota

// MaxUserStorageKB kuota total foto per user (item + bukti pembayaran
// digabung, disimpan di data_peserta.gold_storage_used_kb), berlaku untuk
// semua role KECUALI ADMIN.
const MaxUserStorageKB = 30 * 1024

// StorageDeleteHistory satu baris log foto yang dihapus dari menu Storage --
// murni audit, belum ada tampilan untuk tabel ini.
type StorageDeleteHistory struct {
	HistoryID        int    `gorm:"column:history_id;primaryKey;autoIncrement" json:"history_id"`
	GoldID           int    `gorm:"column:gold_id" json:"gold_id"`
	SourceType       string `gorm:"column:source_type" json:"source_type"`
	SourceID         int    `gorm:"column:source_id" json:"source_id"`
	OriginalFilename string `gorm:"column:original_filename" json:"original_filename"`
	FileBytes        int    `gorm:"column:file_bytes" json:"file_bytes"`
	ContextLabel     string `gorm:"column:context_label" json:"context_label"`
	DeletedBy        string `gorm:"column:deleted_by" json:"deleted_by"`
}

func (StorageDeleteHistory) TableName() string {
	return "storage_delete_history"
}

const (
	SourceTypeItemPhoto    = "ITEM_PHOTO"
	SourceTypePaymentProof = "PAYMENT_PROOF"
)
