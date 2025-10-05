package data_test

import (
	"context"
	"testing"
	"user/internal/conf"
	"user/internal/data"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 测试data方法
func TestData(t *testing.T) {
	//  Ginkgo 测试通过调用 Fail(description string) 功能来表示失败
	// 使用 RegisterFailHandler 将此函数传递给 Gomega 。这是 Ginkgo 和 Gomega 之间的唯一连接点
	RegisterFailHandler(Fail)
	// 通知 Ginkgo 启动测试套件。如果任何 specs 失败，Ginkgo 将自动使 testing.T 失败
	RunSpecs(t, "test biz data")
}

var cleaner func() // 定义删除mysql容器回调函数
var Db *data.Data
var ctx context.Context

// ginkgo 使用 BeforeEach 为您的 Specs 设置状态
var _ = BeforeSuite(func() {
	// 执行测试数据库操作之前，链接之前 docker 容器创建的 mysql
	con, f := data.DockerMysql("mysql", "8.0")
	// 测试完成，关闭容器的回调方法
	cleaner = f
	// 2. 配置数据库连接
	config := &conf.Data{Database: &conf.Data_Database{Driver: "mysql", Source: con}}
	// 3. 初始化数据层
	db := data.NewDB(config)
	mySQLdb, _, err := data.NewData(config, nil, db, nil)
	if err != nil {
		return
	}
	Db = mySQLdb
	// 4. 执行数据库迁移
	err = initialize(db)
	if err != nil {
		return
	}
	Expect(err).NotTo(HaveOccurred())
})

// initialize AutoMigrate gorm自动建表
func initialize(db *gorm.DB) error {
	// 首先创建测试数据库
	if err := db.Exec("CREATE DATABASE IF NOT EXISTS test_db").Error; err != nil {
		return errors.WithStack(err)
	}

	// 切换到测试数据库
	if err := db.Exec("USE test_db").Error; err != nil {
		return errors.WithStack(err)
	}
	// 自动创建表结构
	err := db.AutoMigrate(&data.User{})
	return errors.WithStack(err)
}

// 测试结束后 通过回调函数，关闭并删除 docker 创建的容器
var _ = AfterSuite(func() {
	cleaner()
})
