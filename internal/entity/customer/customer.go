package customer

import "time"

type Customer struct {
	CustID         int        `json:"cust_id" gorm:"column:cust_id;primaryKey;autoIncrement" db:"cust_id"`
	CustGoldID     int        `json:"cust_gold_id" gorm:"column:cust_gold_id" db:"cust_gold_id"`
	CustOutcode    string     `json:"cust_outcode" gorm:"column:cust_outcode" db:"cust_outcode"`
	CustCode       string     `json:"cust_code" gorm:"column:cust_code" db:"cust_code"`
	// CustTypeID     int        `json:"cust_type_id" gorm:"column:cust_type_id" db:"cust_type_id"`
	CustName       string     `json:"cust_name" gorm:"column:cust_name" db:"cust_name"`
	CustOutletName string     `json:"cust_outlet_name" gorm:"column:cust_outlet_name" db:"cust_outlet_name"`
	CustPhone      *string    `json:"cust_phone" gorm:"column:cust_phone" db:"cust_phone"`
	CustAddress    *string    `json:"cust_address" gorm:"column:cust_address" db:"cust_address"`
	CustEmail      *string    `json:"cust_email" gorm:"column:cust_email" db:"cust_email"`
	CustStatus     string     `json:"cust_status" gorm:"column:cust_status" db:"cust_status"`
	CustCreatedAt  time.Time  `json:"cust_created_at" gorm:"column:cust_created_at;autoCreateTime" db:"cust_created_at"`
	CustUpdatedAt  *time.Time `json:"cust_updated_at" gorm:"column:cust_updated_at;autoUpdateTime" db:"cust_updated_at"`
	CustCreatedBy  string     `json:"cust_created_by" gorm:"column:cust_created_by" db:"cust_created_by"`
}

type InsertCustomer struct {
	CustID         int        `json:"cust_id" gorm:"column:cust_id;primaryKey;autoIncrement" db:"cust_id"`
	CustGoldID     int        `json:"cust_gold_id" gorm:"column:cust_gold_id" db:"cust_gold_id"`
	CustOutcode    string     `json:"cust_outcode" gorm:"column:cust_outcode" db:"cust_outcode"`
	CustCode       string     `json:"cust_code" gorm:"column:cust_code" db:"cust_code"`
	// CustTypeID     int        `json:"cust_type_id" gorm:"column:cust_type_id" db:"cust_type_id"`
	CustName       string     `json:"cust_name" gorm:"column:cust_name" db:"cust_name"`
	CustOutletName string     `json:"cust_outlet_name" gorm:"column:cust_outlet_name" db:"cust_outlet_name"`
	CustPhone      string     `json:"cust_phone" gorm:"column:cust_phone" db:"cust_phone"`
	CustAddress    string     `json:"cust_address" gorm:"column:cust_address" db:"cust_address"`
	CustEmail      string     `json:"cust_email" gorm:"column:cust_email" db:"cust_email"`
	CustStatus     string     `json:"cust_status" gorm:"column:cust_status" db:"cust_status"`
	CustCreatedAt  time.Time  `json:"cust_created_at" gorm:"column:cust_created_at;autoCreateTime" db:"cust_created_at"`
	CustUpdatedAt  *time.Time `json:"cust_updated_at" gorm:"column:cust_updated_at;autoUpdateTime" db:"cust_updated_at"`
	CustCreatedBy  string     `json:"cust_created_by" gorm:"column:cust_created_by" db:"cust_created_by"`
}

type CustomerType struct {
	CtID          int        `json:"ct_id" gorm:"column:ct_id;primaryKey;autoIncrement" db:"ct_id"`
	CtGoldID      int        `json:"ct_gold_id" gorm:"column:ct_gold_id" db:"ct_gold_id"`
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

type InsertCustomerData struct {
	CustomerData []Customer `json:"data"`
}

// BulkInsertCustomer payload yang dipublish ke Kafka untuk insert massal customer
// (action "insert_customers_bulk"). GoldID = pemilik (penjual) dari token handler.
type BulkInsertCustomer struct {
	GoldID int        `json:"gold_id"`
	Items  []Customer `json:"items"`
}

type UpdateCustomerData struct {
	UpdateData Customer `json:"data"`
}

func (Customer) TableName() string {
	return "customer"
}

func (CustomerType) TableName() string {
	return "customer_type"
}
