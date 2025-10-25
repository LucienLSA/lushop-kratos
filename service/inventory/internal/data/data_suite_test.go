package data_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"inventory/internal/data"

	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	testDB       *gorm.DB
	testRedis    *redis.Client
	testData     *data.Data
	pool         *dockertest.Pool
	mysqlRes     *dockertest.Resource
	redisRes     *dockertest.Resource
	mysqlDSN     string
	redisAddr    string
	testCtx      context.Context
)

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testCtx = context.Background()

	// 创建 Docker 连接池
	pool, err = dockertest.NewPool("")
	Ω(err).ShouldNot(HaveOccurred())
	pool.MaxWait = 120 * time.Second

	// 启动 MySQL 容器
	mysqlRes, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_DATABASE=inventory_test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Ω(err).ShouldNot(HaveOccurred())

	mysqlPort := mysqlRes.GetPort("3306/tcp")
	mysqlDSN = fmt.Sprintf("root:secret@tcp(localhost:%s)/inventory_test?charset=utf8mb4&parseTime=True&loc=Local", mysqlPort)
	
	fmt.Printf("MySQL 容器启动在: localhost:%s\n", mysqlPort)

	// 等待 MySQL 就绪
	err = pool.Retry(func() error {
		db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			return err
		}
		testDB = db
		sqlDB, _ := db.DB()
		return sqlDB.Ping()
	})
	Ω(err).ShouldNot(HaveOccurred())

	// 启动 Redis 容器
	redisRes, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Ω(err).ShouldNot(HaveOccurred())

	redisPort := redisRes.GetPort("6379/tcp")
	redisAddr = fmt.Sprintf("localhost:%s", redisPort)
	
	fmt.Printf("Redis 容器启动在: %s\n", redisAddr)

	// 等待 Redis 就绪
	err = pool.Retry(func() error {
		testRedis = redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})
		return testRedis.Ping(testCtx).Err()
	})
	Ω(err).ShouldNot(HaveOccurred())

	// 创建测试数据实例
	testData = &data.Data{}

	// 自动迁移表结构
	err = testDB.AutoMigrate(
		&data.Inventory{},
		&data.StockSellDetail{},
	)
	Ω(err).ShouldNot(HaveOccurred())

	fmt.Println("✅ MySQL 和 Redis 测试环境准备完成")
})

var _ = AfterSuite(func() {
	// 清理 MySQL 容器
	if mysqlRes != nil {
		if err := pool.Purge(mysqlRes); err != nil {
			log.Printf("清理 MySQL 容器失败: %v", err)
		}
	}

	// 清理 Redis 容器
	if redisRes != nil {
		if err := pool.Purge(redisRes); err != nil {
			log.Printf("清理 Redis 容器失败: %v", err)
		}
	}

	fmt.Println("🧹 测试环境已清理")
})
