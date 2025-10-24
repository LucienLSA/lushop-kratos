package service

import (
	"context"
	v1 "lushop/api/lushop/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ==================== 地址管理 ====================

// GetAddressList 获取地址列表 HTTP API
func (s *LushopService) GetAddressList(ctx context.Context, req *emptypb.Empty) (*v1.AddressListReply, error) {
	s.log.Info("HTTP API: 获取地址列表")
	return s.useropUc.GetAddressList(ctx)
}

// CreateAddress 创建地址 HTTP API
func (s *LushopService) CreateAddress(ctx context.Context, req *v1.CreateAddressReq) (*v1.AddressReply, error) {
	s.log.Infof("HTTP API: 创建地址 province=%s, city=%s", req.Province, req.City)
	return s.useropUc.CreateAddress(ctx, req)
}

// UpdateAddress 更新地址 HTTP API
func (s *LushopService) UpdateAddress(ctx context.Context, req *v1.UpdateAddressReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 更新地址 id=%d", req.Id)
	err := s.useropUc.UpdateAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteAddress 删除地址 HTTP API
func (s *LushopService) DeleteAddress(ctx context.Context, req *v1.DeleteAddressReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 删除地址 id=%d", req.Id)
	err := s.useropUc.DeleteAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ==================== 留言管理 ====================

// GetMessageList 获取留言列表 HTTP API
func (s *LushopService) GetMessageList(ctx context.Context, req *v1.GetMessageListReq) (*v1.MessageListReply, error) {
	s.log.Infof("HTTP API: 获取留言列表 page=%d, pageSize=%d", req.Page, req.PageSize)
	return s.useropUc.GetMessageList(ctx, req)
}

// CreateMessage 创建留言 HTTP API
func (s *LushopService) CreateMessage(ctx context.Context, req *v1.CreateMessageReq) (*v1.MessageReply, error) {
	s.log.Infof("HTTP API: 创建留言 subject=%s", req.Subject)
	return s.useropUc.CreateMessage(ctx, req)
}

// ==================== 收藏管理 ====================

// GetFavoriteList 获取收藏列表 HTTP API
func (s *LushopService) GetFavoriteList(ctx context.Context, req *emptypb.Empty) (*v1.FavoriteListReply, error) {
	s.log.Info("HTTP API: 获取收藏列表")
	return s.useropUc.GetFavoriteList(ctx)
}

// AddFavorite 添加收藏 HTTP API
func (s *LushopService) AddFavorite(ctx context.Context, req *v1.AddFavoriteReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 添加收藏 goodsId=%d", req.GoodsId)
	err := s.useropUc.AddFavorite(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteFavorite 删除收藏 HTTP API
func (s *LushopService) DeleteFavorite(ctx context.Context, req *v1.DeleteFavoriteReq) (*emptypb.Empty, error) {
	s.log.Infof("HTTP API: 删除收藏 goodsId=%d", req.GoodsId)
	err := s.useropUc.DeleteFavorite(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CheckFavorite 检查是否收藏 HTTP API
func (s *LushopService) CheckFavorite(ctx context.Context, req *v1.CheckFavoriteReq) (*v1.CheckFavoriteReply, error) {
	s.log.Infof("HTTP API: 检查收藏 goodsId=%d", req.GoodsId)
	return s.useropUc.CheckFavorite(ctx, req)
}
