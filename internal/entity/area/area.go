package area

import "time"

type Area struct {
	AreaID       int        `gorm:"column:area_id;primaryKey" db:"area_id" json:"area_id"`
	AreaGoldID   int        `gorm:"column:area_gold_id" db:"area_gold_id" json:"area_gold_id"`
	AreaOutcode  string     `gorm:"column:area_outcode" db:"area_outcode" json:"area_outcode"`
	AreaName     string     `gorm:"column:area_name" db:"area_name" json:"area_name"`
	AreaType     string     `gorm:"column:area_type" db:"area_type" json:"area_type"`
	AreaCreateAt time.Time  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	AreaUpdateAt *time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (Area) TableName() string {
	return "area"
}

const (
	AreaTypeIndoor  = "INDOOR"
	AreaTypeOutdoor = "OUTDOOR"
)

type InsertArea struct {
	AreaName string `json:"area_name"`
	AreaType string `json:"area_type"`
}

type InsertAreaData struct {
	Outcode    string       `json:"outcode"`
	InsertData []InsertArea `json:"data"`
}
