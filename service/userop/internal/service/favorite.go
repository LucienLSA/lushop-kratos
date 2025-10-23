package service

import (
	"context"
	v1 "userop/api/userop/v1"
	"userop/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserOpService) AddUserFav(ctx context.Context, req *v1.UserFavRequest) (*emptypb.Empty, error) {
	userFav := domain.Favorite{
		UserID:  req.UserId,
		GoodsID: req.GoodsId,
	}
	result := s.fs.AddUserFav(ctx, userFav)
	if result != nil {
		return nil, result
	}
	return &emptypb.Empty{}, nil
}

func (s *UserOpService) DeleteUserFav(ctx context.Context, req *v1.UserFavRequest) (*emptypb.Empty, error) {
	userFav := domain.Favorite{
		UserID:  req.UserId,
		GoodsID: req.GoodsId,
	}
	result := s.fs.DeleteUserFav(ctx, userFav)
	if result != nil {
		return nil, result
	}
	return &emptypb.Empty{}, nil
}

func (s *UserOpService) GetUserFavDetail(ctx context.Context, req *v1.UserFavRequest) (*v1.UserFavResponse, error) {
	userFav := domain.Favorite{
		UserID:  req.UserId,
		GoodsID: req.GoodsId,
	}
	result, err := s.fs.GetUserFavDetail(ctx, userFav)
	if err != nil {
		return nil, err
	}
	return &v1.UserFavResponse{
		UserId:  result.UserID,
		GoodsId: result.GoodsID,
	}, nil
}

func (s *UserOpService) GetFavList(ctx context.Context, req *v1.UserFavRequest) (*v1.UserFavListResponse, error) {
	filter := domain.Favorite{
		UserID:  req.UserId,
		GoodsID: req.GoodsId,
	}
	listResp, err := s.fs.GetFavList(ctx, filter)
	if err != nil {
		return nil, err
	}
	out := &v1.UserFavListResponse{Total: int32(listResp.Total)}
	for _, item := range listResp.List {
		out.Data = append(out.Data, &v1.UserFavResponse{UserId: item.UserID, GoodsId: item.GoodsID})
	}
	return out, nil
}
