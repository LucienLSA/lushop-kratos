package data_test

import (
	"fmt"
	"testing"

	"userop/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
)

var (
	db       *gorm.DB
	dataObj  *data.Data // 用于测试中创建 repo
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = BeforeSuite(func() {
	var err error

	// 创建 Docker 连接池
	pool, err = dockertest.NewPool("")
	Ω(err).ShouldNot(HaveOccurred())

	err = pool.Client.Ping()
	Ω(err).ShouldNot(HaveOccurred())

	// 启动 MySQL 容器
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_DATABASE=userop_test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Ω(err).ShouldNot(HaveOccurred())

	// 设置容器过期时间
	_ = resource.Expire(120)

	// 等待 MySQL 启动
	hostAndPort := resource.GetHostPort("3306/tcp")
	fmt.Printf("MySQL 容器启动在: %s\n", hostAndPort)

	dsn := fmt.Sprintf("root:secret@tcp(%s)/userop_test?charset=utf8mb4&parseTime=True&loc=Local", hostAndPort)

	err = pool.Retry(func() error {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		return sqlDB.Ping()
	})
	Ω(err).ShouldNot(HaveOccurred())

	// 自动迁移表结构
	err = db.AutoMigrate(
		&data.Address{},
		&data.LeavingMessages{},
		&data.UserFav{},
	)
	if err != nil {
		fmt.Printf("AutoMigrate 错误: %v\n", err)
	}
	Ω(err).ShouldNot(HaveOccurred())

	// 创建 Data 对象（使用 NewData 函数）
	logger := log.NewStdLogger(io.Discard)
	var cleanup func()
	dataObj, cleanup, err = data.NewData(nil, logger, db, nil)
	Ω(err).ShouldNot(HaveOccurred())
	
	// 注册清理函数
	DeferCleanup(cleanup)

	fmt.Println("✅ MySQL 测试环境准备完成")
})

var _ = AfterSuite(func() {
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}

	if pool != nil && resource != nil {
		err := pool.Purge(resource)
		if err != nil {
			fmt.Printf("清理 MySQL 容器失败: %v\n", err)
		}
	}

	fmt.Println("🧹 测试环境已清理")
})

// 每个测试后清理数据
var _ = AfterEach(func() {
	if db != nil {
		// 清理所有测试数据
		db.Exec("DELETE FROM addresses WHERE signer_name LIKE 'TEST%'")
		db.Exec("DELETE FROM leaving_messages WHERE subject LIKE 'TEST%'")
		db.Exec("DELETE FROM user_favs WHERE user_id >= 9000")
	}
})
