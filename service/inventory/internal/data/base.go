package data

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type GormDetailList []GoodsDetail

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type BaseFields struct {
	ID        int64          `gorm:"primarykey;type:int" json:"id"` // bigint
	CreatedAt time.Time      `gorm:"column:add_time" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// Paginate 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
