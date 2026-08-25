package goldgym

type LoginToken struct {
	GoldToken string `gorm:"column:gold_token" db:"gold_token" json:"gold_token"`
}

type LoginTokenDataPeserta struct {
	GoldToken string `gorm:"column:gold_token" db:"gold_token" json:"gold_token"`
	GoldEmail string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
}

func (LoginToken) TableName() string {
	return "data_token"
}

func (LoginTokenDataPeserta) TableName() string {
	return "data_token"
}
