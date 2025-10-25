package data

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// NewTestData 创建用于测试的 Data 实例
// 这个函数仅用于测试，允许直接设置 db 和 rdb
func NewTestData(db *gorm.DB, rdb *redis.Client) *Data {
	return &Data{
		db:  db,
		rdb: rdb,
	}
}
