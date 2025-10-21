package service

import (
	"context"
	"encoding/json"
	v1 "goods/api/goods/v1"
	"goods/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsService) DeleteCategory(ctx context.Context, r *v1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := g.cac.DeleteCategory(ctx, &domain.CategoryInfo{
		ID: r.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateCategory(ctx context.Context, r *v1.CategoryInfoRequest) (*emptypb.Empty, error) {
	err := g.cac.UpdateCategory(ctx, &domain.CategoryInfo{
		ID:             r.Id,
		Name:           r.Name,
		ParentCategory: r.ParentCategory,
		Level:          r.Level,
		IsTab:          r.IsTab,
	})
	return &emptypb.Empty{}, err
}

func (g *GoodsService) CreateCategory(ctx context.Context, r *v1.CategoryInfoRequest) (*v1.CategoryInfoResponse, error) {
	result, err := g.cac.CreateCategory(ctx, &domain.CategoryInfo{
		Name:           r.Name,
		ParentCategory: r.ParentCategory,
		Level:          r.Level,
		IsTab:          r.IsTab,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CategoryInfoResponse{
		Id:             result.ID,
		Name:           result.Name,
		ParentCategory: result.ParentCategoryID,
		Level:          result.Level,
		IsTab:          result.IsTab,
	}, nil
}

func (g *GoodsService) GetAllCategoryList(ctx context.Context, r *emptypb.Empty) (*v1.CategoryListResponse, error) {
	cate, err := g.cac.CategoryList(ctx)
	if err != nil {
		return nil, err
	}
	// 构造 proto 响应 Data 与 Total
	var data []*v1.CategoryInfoResponse
	for _, c := range cate {
		if c == nil {
			continue
		}
		data = append(data, &v1.CategoryInfoResponse{
			Id:             c.ID,
			Name:           c.Name,
			ParentCategory: c.ParentCategoryID,
			Level:          c.Level,
			IsTab:          c.IsTab,
		})
	}
	jsonData, _ := json.Marshal(cate)
	res := &v1.CategoryListResponse{
		Total:   int32(len(data)),
		Data:    data,
		JsonData: string(jsonData),
	}
	return res, nil
}

// GetSubCategory 获取子分类
func (g *GoodsService) GetSubCategory(ctx context.Context, r *v1.CategoryListRequest) (*v1.SubCategoryListResponse, error) {
	list, err := g.cac.SubCategoryList(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	categoryListRes := v1.SubCategoryListResponse{}
	categoryListRes.Info = &v1.CategoryInfoResponse{
		Id:             list.Category.ID,
		Name:           list.Category.Name,
		ParentCategory: list.Category.ParentCategoryID,
		Level:          list.Category.Level,
		IsTab:          list.Category.IsTab,
	}

	var subCategoryResponse []*v1.CategoryInfoResponse
	for _, subC := range list.SubCategory {
		// list.SubCategory is []*domain.CategoryList, take the Category inside each node
		if subC != nil && subC.Category != nil {
			subCategoryResponse = append(subCategoryResponse, &v1.CategoryInfoResponse{
				Id:             subC.Category.ID,
				Name:           subC.Category.Name,
				ParentCategory: subC.Category.ParentCategoryID,
				Level:          subC.Category.Level,
				IsTab:          subC.Category.IsTab,
			})
		}
	}

	categoryListRes.SubCategorys = subCategoryResponse
	return &categoryListRes, nil
}
