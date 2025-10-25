package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3" // 注意这个包的引入
)

// 启动一个临时的 MySQL Docker 容器用于测试
// 隔离测试环境，避免污染开发或生产数据库
func DockerMysql(img, version string) (string, func()) {
	return innerDockerMysql(img, version)
}

// 初始化 Docker mysql 容器
func innerDockerMysql(img, version string) (string, func()) {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run(img, version, []string{
		"MYSQL_ROOT_PASSWORD=123456",
		"MYSQL_ROOT_HOST=%",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	conStr := fmt.Sprintf("root:123456@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp"))
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", conStr)
		if err != nil {
			return err
		}
		db.SetConnMaxLifetime(60 * time.Second)
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// 回调函数关闭容器
	return conStr, func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
