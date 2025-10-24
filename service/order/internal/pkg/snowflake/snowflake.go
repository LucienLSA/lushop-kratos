package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// Init 初始化雪花算法节点
// nodeID: 节点ID，范围0-1023，用于区分不同的服务实例
func Init(nodeID int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(nodeID)
	})
	return err
}

// GenerateID 生成雪花算法ID
func GenerateID() int64 {
	if node == nil {
		// 如果未初始化，使用默认节点ID 1
		_ = Init(1)
	}
	return node.Generate().Int64()
}

// GenerateIDString 生成字符串格式的雪花算法ID
func GenerateIDString() string {
	if node == nil {
		// 如果未初始化，使用默认节点ID 1
		_ = Init(1)
	}
	return node.Generate().String()
}
