package biz

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserOpRepo 用户操作仓库接口
type UserOpRepo interface {
	// 地址管理
	GetAddressList(ctx context.Context) ([]*v1.AddressReply, int32, error)
	CreateAddress(ctx context.Context, req *v1.CreateAddressReq) (*v1.AddressReply, error)
	UpdateAddress(ctx context.Context, req *v1.UpdateAddressReq) error
	DeleteAddress(ctx context.Context, id int32) error

	// 留言管理
	GetMessageList(ctx context.Context, page, pageSize int32) ([]*v1.MessageReply, int32, error)
	CreateMessage(ctx context.Context, req *v1.CreateMessageReq) (*v1.MessageReply, error)

	// 收藏管理
	GetFavoriteList(ctx context.Context) ([]*v1.FavoriteReply, int32, error)
	AddFavorite(ctx context.Context, goodsId int32) error
	DeleteFavorite(ctx context.Context, goodsId int32) error
	CheckFavorite(ctx context.Context, goodsId int32) (bool, error)
}

// UserOpUsecase 用户操作用例
type UserOpUsecase struct {
	repo UserOpRepo
	log  *log.Helper
}

// NewUserOpUsecase 创建用户操作用例
func NewUserOpUsecase(repo UserOpRepo, logger log.Logger) *UserOpUsecase {
	return &UserOpUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// ==================== 地址管理 ====================

// GetAddressList 获取地址列表
func (uc *UserOpUsecase) GetAddressList(ctx context.Context) (*v1.AddressListReply, error) {
	uc.log.Info("获取地址列表")
	
	addresses, total, err := uc.repo.GetAddressList(ctx)
	if err != nil {
		return nil, err
	}
	
	return &v1.AddressListReply{
		Total: total,
		Data:  addresses,
	}, nil
}

// CreateAddress 创建地址
func (uc *UserOpUsecase) CreateAddress(ctx context.Context, req *v1.CreateAddressReq) (*v1.AddressReply, error) {
	uc.log.Infof("创建地址: province=%s, city=%s", req.Province, req.City)
	return uc.repo.CreateAddress(ctx, req)
}

// UpdateAddress 更新地址
func (uc *UserOpUsecase) UpdateAddress(ctx context.Context, req *v1.UpdateAddressReq) error {
	uc.log.Infof("更新地址: id=%d", req.Id)
	return uc.repo.UpdateAddress(ctx, req)
}

// DeleteAddress 删除地址
func (uc *UserOpUsecase) DeleteAddress(ctx context.Context, req *v1.DeleteAddressReq) error {
	uc.log.Infof("删除地址: id=%d", req.Id)
	return uc.repo.DeleteAddress(ctx, req.Id)
}

// ==================== 留言管理 ====================

// GetMessageList 获取留言列表
func (uc *UserOpUsecase) GetMessageList(ctx context.Context, req *v1.GetMessageListReq) (*v1.MessageListReply, error) {
	uc.log.Infof("获取留言列表: page=%d, pageSize=%d", req.Page, req.PageSize)
	
	messages, total, err := uc.repo.GetMessageList(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	
	return &v1.MessageListReply{
		Total: total,
		Data:  messages,
	}, nil
}

// CreateMessage 创建留言
func (uc *UserOpUsecase) CreateMessage(ctx context.Context, req *v1.CreateMessageReq) (*v1.MessageReply, error) {
	uc.log.Infof("创建留言: subject=%s", req.Subject)
	return uc.repo.CreateMessage(ctx, req)
}

// ==================== 收藏管理 ====================

// GetFavoriteList 获取收藏列表
func (uc *UserOpUsecase) GetFavoriteList(ctx context.Context) (*v1.FavoriteListReply, error) {
	uc.log.Info("获取收藏列表")
	
	favorites, total, err := uc.repo.GetFavoriteList(ctx)
	if err != nil {
		return nil, err
	}
	
	return &v1.FavoriteListReply{
		Total: total,
		Data:  favorites,
	}, nil
}

// AddFavorite 添加收藏
func (uc *UserOpUsecase) AddFavorite(ctx context.Context, req *v1.AddFavoriteReq) error {
	uc.log.Infof("添加收藏: goodsId=%d", req.GoodsId)
	return uc.repo.AddFavorite(ctx, req.GoodsId)
}

// DeleteFavorite 删除收藏
func (uc *UserOpUsecase) DeleteFavorite(ctx context.Context, req *v1.DeleteFavoriteReq) error {
	uc.log.Infof("删除收藏: goodsId=%d", req.GoodsId)
	return uc.repo.DeleteFavorite(ctx, req.GoodsId)
}

// CheckFavorite 检查是否收藏
func (uc *UserOpUsecase) CheckFavorite(ctx context.Context, req *v1.CheckFavoriteReq) (*v1.CheckFavoriteReply, error) {
	uc.log.Infof("检查收藏: goodsId=%d", req.GoodsId)
	
	isFavorite, err := uc.repo.CheckFavorite(ctx, req.GoodsId)
	if err != nil {
		return nil, err
	}
	
	return &v1.CheckFavoriteReply{
		IsFavorite: isFavorite,
	}, nil
}
