package data_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"userauth/internal/biz"
	"userauth/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"io"
)

var (
	redisClient *redis.Client
	dataObj     *data.Data
	repo        biz.AuthRepo
	pool        *dockertest.Pool
	resource    *dockertest.Resource
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
	
	// 启动 Redis 容器
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Ω(err).ShouldNot(HaveOccurred())
	
	// 设置容器过期时间
	_ = resource.Expire(120)
	
	// 等待 Redis 启动
	hostAndPort := resource.GetHostPort("6379/tcp")
	fmt.Printf("Redis 容器启动在: %s\n", hostAndPort)
	
	err = pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: hostAndPort,
			DB:   0,
		})
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		return redisClient.Ping(ctx).Err()
	})
	Ω(err).ShouldNot(HaveOccurred())
	
	// 创建 Data 对象
	logger := log.NewStdLogger(io.Discard)
	dataObj, cleanup, err := data.NewData(redisClient, logger)
	Ω(err).ShouldNot(HaveOccurred())
	
	// 注册清理函数
	DeferCleanup(cleanup)
	
	// 创建 AuthRepo
	repo = data.NewAuthRepo(dataObj, logger)
	
	fmt.Println("✅ Redis 测试环境准备完成")
})

var _ = AfterSuite(func() {
	if redisClient != nil {
		_ = redisClient.Close()
	}
	
	if pool != nil && resource != nil {
		err := pool.Purge(resource)
		if err != nil {
			fmt.Printf("清理 Redis 容器失败: %v\n", err)
		}
	}
	
	fmt.Println("🧹 测试环境已清理")
})

// 每个测试后清理数据
var _ = AfterEach(func() {
	if redisClient != nil {
		ctx := context.Background()
		_ = redisClient.FlushDB(ctx).Err()
	}
})
