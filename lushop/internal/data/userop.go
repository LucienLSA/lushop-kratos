package data

import (
	"context"
	v1 "lushop/api/lushop/v1"
	useropV1 "lushop/api/service/userop/v1"
	"lushop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type useropRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserOpRepo 创建用户操作仓库
func NewUserOpRepo(data *Data, logger log.Logger) biz.UserOpRepo {
	return &useropRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ==================== 地址管理 ====================

// GetAddressList 获取地址列表
func (r *useropRepo) GetAddressList(ctx context.Context) ([]*v1.AddressReply, int32, error) {
	r.log.Info("调用 UserOp gRPC: GetAddressList")

	// 从 JWT 中获取用户ID (这里需要从 context 中获取)
	// userId := GetUserIDFromContext(ctx)
	
	resp, err := r.data.uoc.GetAddressList(ctx, &useropV1.AddressRequest{
		// UserId: userId,
	})
	if err != nil {
		r.log.Errorf("获取地址列表失败: error=%v", err)
		return nil, 0, err
	}

	// 转换为业务对象
	addresses := make([]*v1.AddressReply, 0, len(resp.Data))
	for _, item := range resp.Data {
		addresses = append(addresses, &v1.AddressReply{
			Id:           item.Id,
			Province:     item.Province,
			City:         item.City,
			District:     item.District,
			Address:      item.Address,
			SignerName:   item.SignerName,
			SignerMobile: item.SignerMobile,
		})
	}

	r.log.Infof("获取地址列表成功: total=%d", resp.Total)
	return addresses, resp.Total, nil
}

// CreateAddress 创建地址
func (r *useropRepo) CreateAddress(ctx context.Context, req *v1.CreateAddressReq) (*v1.AddressReply, error) {
	r.log.Infof("调用 UserOp gRPC: CreateAddress province=%s", req.Province)

	resp, err := r.data.uoc.CreateAddress(ctx, &useropV1.AddressRequest{
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	})
	if err != nil {
		r.log.Errorf("创建地址失败: error=%v", err)
		return nil, err
	}

	r.log.Infof("创建地址成功: id=%d", resp.Id)
	return &v1.AddressReply{
		Id:           resp.Id,
		Province:     resp.Province,
		City:         resp.City,
		District:     resp.District,
		Address:      resp.Address,
		SignerName:   resp.SignerName,
		SignerMobile: resp.SignerMobile,
	}, nil
}

// UpdateAddress 更新地址
func (r *useropRepo) UpdateAddress(ctx context.Context, req *v1.UpdateAddressReq) error {
	r.log.Infof("调用 UserOp gRPC: UpdateAddress id=%d", req.Id)

	_, err := r.data.uoc.UpdateAddress(ctx, &useropV1.AddressRequest{
		Id:           req.Id,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	})
	if err != nil {
		r.log.Errorf("更新地址失败: id=%d, error=%v", req.Id, err)
		return err
	}

	r.log.Infof("更新地址成功: id=%d", req.Id)
	return nil
}

// DeleteAddress 删除地址
func (r *useropRepo) DeleteAddress(ctx context.Context, id int32) error {
	r.log.Infof("调用 UserOp gRPC: DeleteAddress id=%d", id)

	_, err := r.data.uoc.DeleteAddress(ctx, &useropV1.AddressRequest{
		Id: id,
	})
	if err != nil {
		r.log.Errorf("删除地址失败: id=%d, error=%v", id, err)
		return err
	}

	r.log.Infof("删除地址成功: id=%d", id)
	return nil
}

// ==================== 留言管理 ====================

// GetMessageList 获取留言列表
func (r *useropRepo) GetMessageList(ctx context.Context, page, pageSize int32) ([]*v1.MessageReply, int32, error) {
	r.log.Infof("调用 UserOp gRPC: MessageList page=%d, pageSize=%d", page, pageSize)

	resp, err := r.data.uoc.MessageList(ctx, &useropV1.MessageRequest{
		// 可以添加分页参数
	})
	if err != nil {
		r.log.Errorf("获取留言列表失败: error=%v", err)
		return nil, 0, err
	}

	// 转换为业务对象
	messages := make([]*v1.MessageReply, 0, len(resp.Data))
	for _, item := range resp.Data {
		messages = append(messages, &v1.MessageReply{
			Id:          item.Id,
			UserId:      item.UserId,
			MessageType: item.MessageType,
			Subject:     item.Subject,
			Message:     item.Message,
			File:        item.File,
		})
	}

	r.log.Infof("获取留言列表成功: total=%d", resp.Total)
	return messages, resp.Total, nil
}

// CreateMessage 创建留言
func (r *useropRepo) CreateMessage(ctx context.Context, req *v1.CreateMessageReq) (*v1.MessageReply, error) {
	r.log.Infof("调用 UserOp gRPC: CreateMessage subject=%s", req.Subject)

	resp, err := r.data.uoc.CreateMessage(ctx, &useropV1.MessageRequest{
		MessageType: req.MessageType,
		Subject:     req.Subject,
		Message:     req.Message,
		File:        req.File,
	})
	if err != nil {
		r.log.Errorf("创建留言失败: error=%v", err)
		return nil, err
	}

	r.log.Infof("创建留言成功: id=%d", resp.Id)
	return &v1.MessageReply{
		Id:          resp.Id,
		UserId:      resp.UserId,
		MessageType: resp.MessageType,
		Subject:     resp.Subject,
		Message:     resp.Message,
		File:        resp.File,
	}, nil
}

// ==================== 收藏管理 ====================

// GetFavoriteList 获取收藏列表
func (r *useropRepo) GetFavoriteList(ctx context.Context) ([]*v1.FavoriteReply, int32, error) {
	r.log.Info("调用 UserOp gRPC: GetFavList")

	resp, err := r.data.uoc.GetFavList(ctx, &useropV1.UserFavRequest{
		// UserId: userId,
	})
	if err != nil {
		r.log.Errorf("获取收藏列表失败: error=%v", err)
		return nil, 0, err
	}

	// 转换为业务对象
	favorites := make([]*v1.FavoriteReply, 0, len(resp.Data))
	for _, item := range resp.Data {
		favorites = append(favorites, &v1.FavoriteReply{
			UserId:  item.UserId,
			GoodsId: item.GoodsId,
		})
	}

	r.log.Infof("获取收藏列表成功: total=%d", resp.Total)
	return favorites, resp.Total, nil
}

// AddFavorite 添加收藏
func (r *useropRepo) AddFavorite(ctx context.Context, goodsId int32) error {
	r.log.Infof("调用 UserOp gRPC: AddUserFav goodsId=%d", goodsId)

	_, err := r.data.uoc.AddUserFav(ctx, &useropV1.UserFavRequest{
		GoodsId: goodsId,
	})
	if err != nil {
		r.log.Errorf("添加收藏失败: goodsId=%d, error=%v", goodsId, err)
		return err
	}

	r.log.Infof("添加收藏成功: goodsId=%d", goodsId)
	return nil
}

// DeleteFavorite 删除收藏
func (r *useropRepo) DeleteFavorite(ctx context.Context, goodsId int32) error {
	r.log.Infof("调用 UserOp gRPC: DeleteUserFav goodsId=%d", goodsId)

	_, err := r.data.uoc.DeleteUserFav(ctx, &useropV1.UserFavRequest{
		GoodsId: goodsId,
	})
	if err != nil {
		r.log.Errorf("删除收藏失败: goodsId=%d, error=%v", goodsId, err)
		return err
	}

	r.log.Infof("删除收藏成功: goodsId=%d", goodsId)
	return nil
}

// CheckFavorite 检查是否收藏
func (r *useropRepo) CheckFavorite(ctx context.Context, goodsId int32) (bool, error) {
	r.log.Infof("调用 UserOp gRPC: GetUserFavDetail goodsId=%d", goodsId)

	resp, err := r.data.uoc.GetUserFavDetail(ctx, &useropV1.UserFavRequest{
		GoodsId: goodsId,
	})
	if err != nil {
		// 如果返回错误，说明没有收藏
		r.log.Infof("商品未收藏: goodsId=%d", goodsId)
		return false, nil
	}

	isFavorite := resp != nil && resp.GoodsId == goodsId
	r.log.Infof("检查收藏成功: goodsId=%d, isFavorite=%v", goodsId, isFavorite)
	return isFavorite, nil
}
