package service

import (
	"testing"
	"time"

	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

// Service 层测试
// 注意：Service 层主要负责数据转换和 gRPC 接口适配
// 业务逻辑测试已在 Biz 层完成，这里主要测试转换函数

// TestUserResponse 测试 UserResponse 转换函数（有生日）
func TestUserResponse(t *testing.T) {
	birthday := time.Now()
	user := &biz.User{
		ID:       1,
		Mobile:   "13800138000",
		Password: "encrypted_password",
		NickName: "测试用户",
		Gender:   "male",
		Role:     1,
		Birthday: &birthday,
	}

	resp := UserResponse(user)

	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.Id)
	assert.Equal(t, user.Mobile, resp.Mobile)
	assert.Equal(t, user.Password, resp.Password)
	assert.Equal(t, user.NickName, resp.NickName)
	assert.Equal(t, user.Gender, resp.Gender)
	assert.Equal(t, int32(user.Role), resp.Role)
	assert.Equal(t, uint64(birthday.Unix()), resp.Birthday)
}

// TestUserResponse_NoBirthday 测试没有生日的情况
func TestUserResponse_NoBirthday(t *testing.T) {
	user := &biz.User{
		ID:       1,
		Mobile:   "13800138000",
		Password: "encrypted_password",
		NickName: "测试用户",
		Gender:   "male",
		Role:     1,
		Birthday: nil,
	}

	resp := UserResponse(user)

	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.Id)
	assert.Equal(t, uint64(0), resp.Birthday)
}

// TestNewUserService 测试 Service 创建
func TestNewUserService(t *testing.T) {
	logger := log.NewStdLogger(nil)
	uc := biz.NewUserUsecase(nil, logger)
	service := NewUserService(uc, logger)

	assert.NotNil(t, service)
	assert.NotNil(t, service.uc)
	assert.NotNil(t, service.log)
}

// 表驱动测试示例：测试 UserResponse 的不同场景
func TestUserResponse_TableDriven(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name            string
		user            *biz.User
		expectedBirthday uint64
	}{
		{
			name: "完整用户信息",
			user: &biz.User{
				ID:       1,
				Mobile:   "13800138000",
				NickName: "用户1",
				Gender:   "male",
				Role:     1,
				Birthday: &now,
			},
			expectedBirthday: uint64(now.Unix()),
		},
		{
			name: "无生日信息",
			user: &biz.User{
				ID:       2,
				Mobile:   "13800138001",
				NickName: "用户2",
				Gender:   "female",
				Role:     1,
				Birthday: nil,
			},
			expectedBirthday: 0,
		},
		{
			name: "管理员用户",
			user: &biz.User{
				ID:       3,
				Mobile:   "13800138002",
				NickName: "管理员",
				Gender:   "male",
				Role:     2,
				Birthday: &now,
			},
			expectedBirthday: uint64(now.Unix()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := UserResponse(tt.user)

			assert.NotNil(t, resp)
			assert.Equal(t, tt.user.ID, resp.Id)
			assert.Equal(t, tt.user.Mobile, resp.Mobile)
			assert.Equal(t, tt.user.NickName, resp.NickName)
			assert.Equal(t, tt.user.Gender, resp.Gender)
			assert.Equal(t, int32(tt.user.Role), resp.Role)
			assert.Equal(t, tt.expectedBirthday, resp.Birthday)
		})
	}
}
