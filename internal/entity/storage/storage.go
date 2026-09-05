package goldgym

import "time"

// StorageEntry satu baris foto milik user di menu Storage -- bisa berasal
// dari items.item_photo atau payment_proof, disatukan lewat SourceType.
type StorageEntry struct {
	SourceType  string    `json:"source_type"` // ITEM_PHOTO | PAYMENT_PROOF
	SourceID    int       `json:"source_id"`
	Label       string    `json:"label"`        // keperluan, mis. "Foto Katalog Item"
	ContextText string    `json:"context_text"` // nama item / nomor nota
	SizeKB      int       `json:"size_kb"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// StorageSummary ringkasan pemakaian + daftar foto milik satu user.
type StorageSummary struct {
	UsedKB  int            `json:"used_kb"`
	QuotaKB int            `json:"quota_kb"`
	Entries []StorageEntry `json:"entries"`
}
