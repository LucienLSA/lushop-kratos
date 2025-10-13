package snowflake

import (
	"sync"
	"testing"
)

func TestInit(t *testing.T) {
	err := Init(1)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestGenerateID(t *testing.T) {
	// 初始化节点
	err := Init(1)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 生成多个ID并验证它们都是唯一的
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := GenerateID()
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
		if ids[id] {
			t.Errorf("Generated duplicate ID: %d", id)
		}
		ids[id] = true
	}
}

func TestGenerateIDString(t *testing.T) {
	// 初始化节点
	err := Init(1)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 生成多个字符串ID并验证它们都是唯一的
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := GenerateIDString()
		if id == "" {
			t.Errorf("Generated ID should not be empty")
		}
		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

func TestGenerateIDWithoutInit(t *testing.T) {
	// 重置全局变量以测试未初始化的情况
	node = nil
	once = sync.Once{}

	// 不调用Init，直接生成ID
	id := GenerateID()
	if id <= 0 {
		t.Errorf("Generated ID should be positive even without explicit init, got: %d", id)
	}
}
