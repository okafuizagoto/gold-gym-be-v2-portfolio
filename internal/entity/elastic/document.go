package elastic

// UserDocument — dokumen yang di-index ke ES
type UserDocument struct {
	GoldId      int    `json:"gold_id"`
	GoldEmail   string `json:"gold_email"`
	GoldNama    string `json:"gold_nama"`
	GoldNomorHp string `json:"gold_nomorhp"`
	IndexedAt   string `json:"indexed_at"`
}

// SearchResult — response dari ES search
type SearchResult struct {
	Total int            `json:"total"`
	Hits  []UserDocument `json:"hits"`
}
