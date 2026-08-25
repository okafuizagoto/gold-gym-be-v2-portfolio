package customer

import "time"

type CustomerType struct {
	CtID          int        `json:"ct_id" gorm:"column:ct_id;primaryKey;autoIncrement" db:"ct_id"`
	CtGoldID      int        `json:"ct_gold_id" gorm:"column:ct_gold_id" db:"ct_gold_id"`
	CtOutcode string     `json:"ct_outcode" gorm:"column:ct_outcode" db:"ct_outcode"`
	CtCode        string     `json:"ct_code" gorm:"column:ct_code" db:"ct_code"`
	CtName        string     `json:"ct_name" gorm:"column:ct_name" db:"ct_name"`
	CtCategory    string     `json:"ct_category" gorm:"column:ct_category" db:"ct_category"`
	CtDescription string     `json:"ct_description" gorm:"column:ct_description" db:"ct_description"`
	CtStatus      string     `json:"ct_status" gorm:"column:ct_status" db:"ct_status"`
	CtCreatedAt   time.Time  `json:"ct_created_at" gorm:"column:ct_created_at;autoCreateTime" db:"ct_created_at"`
	CtUpdatedAt   *time.Time `json:"ct_updated_at" gorm:"column:ct_updated_at;autoUpdateTime" db:"ct_updated_at"`
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type InsertCustomerTypeData struct {
	CustomerData []CustomerType `json:"data"`
}

type UpdateCustomerTypeData struct {
	UpdateData CustomerType `json:"data"`
}

func (CustomerType) TableName() string {
	return "customer_type"
}
