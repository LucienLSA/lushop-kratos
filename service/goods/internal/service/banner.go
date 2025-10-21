package service

import (
	"context"
	v1 "goods/api/goods/v1"
	"goods/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsService) BannerList(ctx context.Context, req *emptypb.Empty) (*v1.BannerListResponse, error) {
	banner, total, err := g.ban.BannerList(ctx)
	if err != nil {
		return nil, err
	}
	var data []*v1.BannerResponse
	for _, banner := range banner {
		data = append(data, &v1.BannerResponse{
			Id:    banner.ID,
			Url:   banner.Url,
			Index: banner.Index,
		})
	}
	return &v1.BannerListResponse{Data: data, Total: int32(total)}, nil
}

func (g *GoodsService) CreateBanner(ctx context.Context, req *v1.BannerRequest) (*v1.BannerResponse, error) {
	banner := &domain.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}
	result, err := g.ban.CreateBanner(ctx, banner)
	if err != nil {
		return nil, err
	}
	return &v1.BannerResponse{Id: result.ID, Url: result.Url}, nil
}

func (g *GoodsService) DeleteBanner(ctx context.Context, req *v1.BannerRequest) (*emptypb.Empty, error) {
	err := g.ban.DeleteBanner(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateBanner(ctx context.Context, req *v1.BannerRequest) (*emptypb.Empty, error) {
	banner := &domain.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}
	err := g.ban.UpdateBanner(ctx, banner)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
