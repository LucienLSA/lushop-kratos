package data

import (
	"context"
	"time"

	v1 "lushop/api/lushop/v1"
	userService "lushop/api/service/user/v1"
	userauthV1 "lushop/api/service/userauth/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// userAuthRepoGRPC 通过 gRPC 调用用户认证服务实现 UserRepo 接口
// 认证域方法（验证码/SMS/Token/黑名单）走 UserAuth 服务
// 用户资料方法（CRUD）继续走 User 服务
type userAuthRepoGRPC struct {
	data *Data
	log  *log.Helper
}

// NewUserAuthRepoGRPC 创建 gRPC 版本的用户仓库（用于统一治理方案）
func NewUserAuthRepoGRPC(data *Data, logger log.Logger) biz.UserRepo {
	return &userAuthRepoGRPC{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/user-auth-grpc")),
	}
}

// ========== 用户资料相关（继续通过 User 服务） ==========

func (r *userAuthRepoGRPC) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	createUser, err := r.data.uc.CreateUser(ctx, &userService.CreateUserInfo{
		NickName: user.NickName,
		Password: user.Password,
		Mobile:   user.Mobile,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:       createUser.Id,
		Mobile:   createUser.Mobile,
		NickName: createUser.NickName,
	}, nil
}

func (r *userAuthRepoGRPC) UserByMobile(ctx context.Context, mobile string) (*biz.User, error) {
	byMobile, err := r.data.uc.GetUserByMobile(ctx, &userService.MobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Mobile:   byMobile.Mobile,
		ID:       byMobile.Id,
		NickName: byMobile.NickName,
		Password: byMobile.Password,
		Role:     int(byMobile.Role),
		Birthday: int64(byMobile.Birthday),
		Gender:   byMobile.Gender,
	}, nil
}

func (r *userAuthRepoGRPC) UserById(ctx context.Context, id int64) (*biz.User, error) {
	user, err := r.data.uc.GetUserById(ctx, &userService.IdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:       user.Id,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int(user.Role),
		Password: user.Password,
		Birthday: int64(user.Birthday),
	}, nil
}

func (r *userAuthRepoGRPC) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	_, err := r.data.uc.UpdateUser(ctx, &userService.UpdateUserInfo{
		Id:       user.ID,
		NickName: user.NickName,
		Password: user.Password,
		Gender:   user.Gender,
		Birthday: uint64(user.Birthday),
	})
	if err != nil {
		return nil, err
	}

	updatedUser, err := r.data.uc.GetUserById(ctx, &userService.IdRequest{
		Id: user.ID,
	})
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:       updatedUser.Id,
		Mobile:   updatedUser.Mobile,
		NickName: updatedUser.NickName,
		Password: updatedUser.Password,
		Birthday: int64(updatedUser.Birthday),
		Gender:   updatedUser.Gender,
		Role:     int(updatedUser.Role),
	}, nil
}

func (r *userAuthRepoGRPC) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	result, err := r.data.uc.CheckPassword(ctx, &userService.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPassword,
	})
	if err != nil {
		return false, err
	}
	return result.Success, nil
}

func (r *userAuthRepoGRPC) ListUsers(ctx context.Context, req *v1.ListUsersReq) ([]*biz.User, int, error) {
	pageInfo := &userService.PageInfo{
		Pn:    uint32(req.Page),
		PSize: uint32(req.PageSize),
	}

	userListResp, err := r.data.uc.GetUserList(ctx, pageInfo)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*biz.User, 0, len(userListResp.Data))
	for _, user := range userListResp.Data {
		result = append(result, &biz.User{
			ID:       user.Id,
			Mobile:   user.Mobile,
			NickName: user.NickName,
			Role:     int(user.Role),
			Password: user.Password,
		})
	}

	return result, int(userListResp.Total), nil
}

// ========== 黑名单相关（通过 UserAuth 服务） ==========

func (r *userAuthRepoGRPC) StoreLogoutBlacklist(ctx context.Context, userId int64) error {
	_, err := r.data.uac.AddToBlacklist(ctx, &userauthV1.AddToBlacklistReq{
		UserId:     userId,
		TtlSeconds: 0, // 使用服务端默认 TTL
	})
	if err != nil {
		r.log.Errorf("添加黑名单失败: %v", err)
		return err
	}
	return nil
}

func (r *userAuthRepoGRPC) CheckLogoutBlacklist(ctx context.Context, userId int64) (bool, error) {
	reply, err := r.data.uac.CheckBlacklist(ctx, &userauthV1.CheckBlacklistReq{
		UserId: userId,
	})
	if err != nil {
		r.log.Errorf("检查黑名单失败: %v", err)
		return false, err
	}
	return reply.IsBlacklisted, nil
}

func (r *userAuthRepoGRPC) StoreLogoutBlacklistWithTTL(ctx context.Context, userId int64, ttl time.Duration) error {
	_, err := r.data.uac.AddToBlacklist(ctx, &userauthV1.AddToBlacklistReq{
		UserId:     userId,
		TtlSeconds: int64(ttl.Seconds()),
	})
	if err != nil {
		r.log.Errorf("添加黑名单（自定义TTL）失败: %v", err)
		return err
	}
	return nil
}
