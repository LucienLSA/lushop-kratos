package service

import (
	"context"
	v1 "goods/api/goods/v1"
	"goods/internal/biz"
	"goods/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsService) CategoryBrandList(ctx context.Context, r *v1.CategoryBrandFilterRequest) (*v1.CategoryBrandListResponse, error) {
	pg := &biz.Pagination{
		PageNum:  int(r.Pages),
		PageSize: int(r.PagePerNums),
	}
	list, total, err := g.cacb.CategoryBrandList(ctx, pg)
	if err != nil {
		return nil, err
	}
	resp := &v1.CategoryBrandListResponse{
		Total: int32(total),
	}
	for _, item := range *list {
		resp.Data = append(resp.Data, &v1.CategoryBrandResponse{
			Id: int32(item.ID),
			Brand: &v1.BrandInfoResponse{
				Id: int32(item.BrandsID),
			},
			Category: &v1.CategoryInfoResponse{
				Id: int32(item.CategoryID),
			},
		})
	}
	return resp, nil
}

func (g *GoodsService) GetCategoryBrandList(ctx context.Context, r *v1.CategoryInfoRequest) (*v1.BrandListResponse, error) {
	categoryBrands, err := g.cacb.GetCategoryBrandList(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	// 收集品牌ID
	var ids []int32
	if categoryBrands != nil {
		for _, cb := range *categoryBrands {
			ids = append(ids, cb.BrandsID)
		}
	}
	// 批量查询品牌信息
	brands, err := g.bc.ListByIds(ctx, ids...)
	if err != nil {
		return nil, err
	}
	resp := &v1.BrandListResponse{Total: int32(len(brands))}
	for _, b := range brands {
		resp.Data = append(resp.Data, &v1.BrandInfoResponse{
			Id:   int32(b.ID),
			Name: b.Name,
			Logo: b.Logo,
		})
	}
	return resp, nil
}

func (g *GoodsService) CreateCategoryBrand(ctx context.Context, r *v1.CategoryBrandRequest) (*v1.CategoryBrandResponse, error) {
	categoryBrands := &domain.CategoryBrand{
		ID:         r.Id,
		CategoryID: r.CategoryId,
		BrandsID:   r.BrandId,
	}
	result, err := g.cacb.CreateCategoryBrand(ctx, categoryBrands)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryBrandResponse{
		Id: int32(result.ID),
		Category: &v1.CategoryInfoResponse{
			Id: int32(result.CategoryID),
		},
		Brand: &v1.BrandInfoResponse{
			Id: int32(result.BrandsID),
		},
	}, nil
}

func (g *GoodsService) DeleteCategoryBrand(ctx context.Context, r *v1.CategoryBrandRequest) (*emptypb.Empty, error) {
	err := g.cacb.DeleteCategoryBrand(ctx, &domain.CategoryBrand{
		ID: r.Id,
	})
	return &emptypb.Empty{}, err
}

func (g *GoodsService) UpdateCategoryBrand(ctx context.Context, r *v1.CategoryBrandRequest) (*emptypb.Empty, error) {
	categoryBrands := &domain.CategoryBrand{
		ID:         r.Id,
		CategoryID: r.CategoryId,
		BrandsID:   r.BrandId,
	}
	err := g.cacb.UpdateCategoryBrand(ctx, categoryBrands)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
