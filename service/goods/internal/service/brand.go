package service

import (
	"context"
	v1 "goods/api/goods/v1"
	"goods/internal/biz"
	"goods/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsService) CreateBrand(ctx context.Context, r *v1.BrandRequest) (*v1.BrandInfoResponse, error) {
	brand := domain.Brand{
		ID:   r.Id,
		Name: r.Name,
		Logo: r.Logo,
	}
	result, err := g.bc.CreateBrand(ctx, &brand)
	if err != nil {
		return nil, err
	}
	return &v1.BrandInfoResponse{Id: int32(result.ID)}, nil
}

func (g *GoodsService) UpdateBrand(ctx context.Context, r *v1.BrandRequest) (*emptypb.Empty, error) {
	brand := domain.Brand{
		ID:   r.Id,
		Name: r.Name,
		Logo: r.Logo,
	}
	err := g.bc.UpdateBrand(ctx, &brand)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsService) BrandList(ctx context.Context, r *v1.BrandFilterRequest) (*v1.BrandListResponse, error) {
	pg := &biz.Pagination{PageNum: int(r.Pages), PageSize: int(r.PagePerNums)}
	list, total, err := g.bc.BrandList(ctx, pg)
	if err != nil {
		return nil, err
	}
	resp := &v1.BrandListResponse{Total: int32(total)}
	for _, b := range list {
		resp.Data = append(resp.Data, &v1.BrandInfoResponse{
			Id:   b.ID,
			Name: b.Name,
			Logo: b.Logo,
		})
	}
	return resp, nil
}

func (g *GoodsService) DeleteBrand(ctx context.Context, r *v1.BrandRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, g.bc.DeleteBrand(ctx, r.Id)
}
