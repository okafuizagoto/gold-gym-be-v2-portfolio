package goldgym

// testings
type Testings struct {
	ID            string  `gorm:"column:id" db:"id" json:"id"`
	Filename      *string `gorm:"column:filename" db:"filename" json:"filename"`
	TestingImages []byte  `gorm:"column:testing" db:"testing" json:"testing"`
}

func (Testings) TableName() string {
	return "image_test"
}
