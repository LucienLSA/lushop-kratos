package biz

import (
	"context"
	"errors"
	"goods/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type SpecificationRepo interface {
	CreateSpecification(context.Context, *domain.Specification) (int64, error)
	CreateSpecificationValue(context.Context, int64, []*domain.SpecificationValue) error
	ListByIds(ctx context.Context, id ...*int64) (domain.SpecificationList, error)
}
type SpecificationUsecase struct {
	repo  SpecificationRepo
	gRepo GoodsTypeRepo // 加入商品类型的 repo，用来查询类型是否存在
	tx    Transaction   // 引入 gorm 事务
	log   *log.Helper
}

func NewSpecificationUsecase(repo SpecificationRepo, TypeUc GoodsTypeRepo, tx Transaction,
	logger log.Logger) *SpecificationUsecase {
	return &SpecificationUsecase{
		repo:  repo,
		gRepo: TypeUc,
		tx:    tx,
		log:   log.NewHelper(logger),
	}
}

func (s *SpecificationUsecase) CreateSpecification(ctx context.Context, r *domain.Specification) (int64, error) {
	var (
		id  int64
		err error
	)
	// domain层的判断typeid方法
	if r.IsTypeIDEmpty() {
		return id, errors.New("请选择商品类型进行绑定")
	}

	if r.IsValueEmpty() {
		return id, errors.New("请填写商品规格下的参数")
	}
	// 查询是否存在这个类型
	_, err = s.gRepo.IsExistsByID(ctx, r.TypeID)
	if err != nil {
		return id, err
	}
	// 事务创建Specification和SpecificationValue
	err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
		id, err = s.repo.CreateSpecification(ctx, r)
		if err != nil {
			return err
		}
		err = s.repo.CreateSpecificationValue(ctx, id, r.SpecificationValue) // 插入商品规格
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}
