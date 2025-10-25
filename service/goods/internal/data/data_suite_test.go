package data_test

import (
	"context"
	"goods/internal/conf"
	"goods/internal/data"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	cleaner func()
	Db      *data.Data
	ctx     context.Context
)

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goods Data Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.Background()
	
	// 使用 Docker 启动测试数据库
	con, f := data.DockerMysql("mysql", "8.0")
	cleaner = f

	// 配置数据库连接
	config := &conf.Data{
		Database: &conf.Data_Database{
			Driver: "mysql",
			Source: con,
		},
	}

	// 初始化数据层
	db := data.NewDB(config)
	mySQLdb, _, err := data.NewData(config, nil, db, nil, nil)
	Ω(err).ShouldNot(HaveOccurred())
	
	Db = mySQLdb

	// 执行数据库迁移
	err = initialize(db)
	Ω(err).ShouldNot(HaveOccurred())
})

// initialize 自动建表
func initialize(db *gorm.DB) error {
	// 创建测试数据库
	if err := db.Exec("CREATE DATABASE IF NOT EXISTS test_goods_db").Error; err != nil {
		return errors.WithStack(err)
	}

	// 切换到测试数据库
	if err := db.Exec("USE test_goods_db").Error; err != nil {
		return errors.WithStack(err)
	}

	// 先创建依赖表
	if err := db.AutoMigrate(&data.Category{}, &data.Brand{}, &data.Banner{}); err != nil {
		return errors.WithStack(err)
	}

	// 手动创建 Goods 表，跳过外键关联
	err := db.Migrator().AutoMigrate(&GoodsForTest{})
	return errors.WithStack(err)
}

// GoodsForTest 用于测试的商品表结构（移除关联字段）
type GoodsForTest struct {
	data.BaseFields
	CategoryID      int32           `gorm:"type:int;not null;comment:'商品分类ID'"`
	BrandsID        int32           `gorm:"type:int;not null"`
	OnSale          bool            `gorm:"default:false;not null;comment:'是否特价'"`
	GoodsSn         string          `gorm:"type:varchar(50);not null;comment:'商品编号'"`
	Name            string          `gorm:"type:varchar(100);not null;comment:'商品名称'"`
	ClickNum        int32           `gorm:"type:int;default:0;not null;comment:'商品点击数'"`
	SoldNum         int32           `gorm:"type:int;default:0;not null;comment:'商品销量'"`
	FavNum          int32           `gorm:"type:int;default:0;not null;comment:'商品收藏数'"`
	MarketPrice     float32         `gorm:"not null;comment:'商品市场价'"`
	ShopPrice       float32         `gorm:"not null;comment:'商品实际价'"`
	GoodsBrief      string          `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	ShipFree        bool            `gorm:"default:false;not null;comment:'是否免运费'"`
	Images          data.GormList   `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages      data.GormList   `gorm:"type:varchar(5000);not null;comment:'商品详情图片'"`
	GoodsFrontImage string          `gorm:"type:varchar(200);not null;comment:'商品封面图'"`
	IsNew           bool            `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot           bool            `gorm:"default:false;not null;comment:'是否热卖'"`
}

// TableName 指定表名
func (GoodsForTest) TableName() string {
	return "goods"
}

var _ = AfterSuite(func() {
	if cleaner != nil {
		cleaner()
	}
})

// CleanTestData 清理测试数据的辅助函数
func CleanTestData() {
	if Db != nil {
		Db.DB(ctx).Exec("DELETE FROM goods WHERE goods_sn LIKE 'TEST%'")
	}
}
