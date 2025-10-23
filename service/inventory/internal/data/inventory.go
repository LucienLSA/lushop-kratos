package data

import (
	"context"
	"fmt"
	"inventory/internal/biz"
	"inventory/internal/domain"
	"sort"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"gorm.io/gorm"
)

type Inventory struct {
	BaseFields
	Goods   int32 `gorm:"type:int;index;comment:商品id"`
	Stocks  int32 `gorm:"type:int;comment:仓库"`
	Version int32 `gorm:"type:int;comment:分布式锁-乐观锁"` // 分布式锁的乐观锁
}
type InventoryNew struct {
	BaseFields
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
	Freeze  int32 `gorm:"type:int"` //冻结库存
}

type Delivery struct {
	Goods   int32  `gorm:"type:int;index"`
	Nums    int32  `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status  string `gorm:"type:varchar(200)"` // 1.代表等待支付，2.代表支付成功，3.支付失败
}

type StockSellDetail struct {
	BaseFields
	OrderSn string         `gorm:"type:varchar(200);index:idx_order_sn,unique;comment:订单编号"`
	Status  int32          `gorm:"type:varchar(200);comment:1.已扣减,2.已归还"` // 1.代表已扣减，2.代表已归还，3.失败
	Detail  GormDetailList `gorm:"type:varchar(200);comment:详细商品"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

type inventoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewInventoryRepo .
func NewInventoryRepo(data *Data, logger log.Logger) biz.InventoryRepo {
	return &inventoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (p *Inventory) ToDomain() *domain.Inventory {
	return &domain.Inventory{
		Goods:  int32(p.Goods),
		Stocks: int32(p.Stocks),
	}
}
func (r *inventoryRepo) GetInvById(ctx context.Context, goodsId int32) (*domain.Inventory, error) {
	var inv Inventory
	err := r.data.DB(ctx).Where(&Inventory{Goods: goodsId}).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return inv.ToDomain(), nil
}

func (r *inventoryRepo) AddInv(ctx context.Context, inv *domain.Inventory) error {
	db := r.data.DB(ctx)
	var m Inventory
	// try to find existing inventory by goods id
	res := db.Where("goods = ?", inv.Goods).First(&m)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			// create new record
			m = Inventory{
				Goods:  inv.Goods,
				Stocks: inv.Stocks,
			}
			if err := db.Create(&m).Error; err != nil {
				return errors.InternalServer("ADD_INV_ERROR", err.Error())
			}
			return nil
		}
		return errors.InternalServer("GET_INV_ERROR", res.Error.Error())
	}
	// update existing stocks
	m.Stocks = inv.Stocks
	if err := db.Save(&m).Error; err != nil {
		return errors.InternalServer("UPDATE_INV_ERROR", err.Error())
	}
	return nil
}

func (r *inventoryRepo) Sell(ctx context.Context, sell *domain.SellInfo) error {
	// initialize redsync with existing redis client
	pool := goredis.NewPool(r.data.rdb)
	rs := redsync.New(pool)

	return r.data.ExecTx(ctx, func(ctx context.Context) error {
		db := r.data.DB(ctx)
		// sort by GoodsID to avoid potential lock order deadlocks
		items := make([]domain.GoodsInvInfo, len(sell.GoodsInvInfo))
		copy(items, sell.GoodsInvInfo)
		sort.Slice(items, func(i, j int) bool { return items[i].GoodsID < items[j].GoodsID })

		// deduct stock with guarded updates to prevent oversell, lock per item sequentially
		details := make(GormDetailList, 0, len(items))
		for _, gi := range items {
			m := rs.NewMutex(
				fmt.Sprintf("goods_%d", gi.GoodsID),
				redsync.WithExpiry(5*time.Second),
				redsync.WithTries(3),
				redsync.WithRetryDelay(100*time.Millisecond),
			)
			if err := m.Lock(); err != nil {
				return errors.InternalServer("LOCK_ERROR", "failed to acquire distributed lock")
			}
			res := db.Model(&Inventory{}).
				Where("goods = ? AND stocks >= ?", gi.GoodsID, gi.Nums).
				Update("stocks", gorm.Expr("stocks - ?", gi.Nums))
			if res.Error != nil {
				_, _ = m.Unlock()
				return errors.InternalServer("DEDUCT_INV_ERROR", res.Error.Error())
			}
			if res.RowsAffected == 0 {
				_, _ = m.Unlock()
				return errors.BadRequest("STOCK_NOT_ENOUGH", "not enough stock or goods not found")
			}
			details = append(details, GoodsDetail{Goods: gi.GoodsID, Num: gi.Nums})
			_, _ = m.Unlock()
		}
		rec := StockSellDetail{OrderSn: sell.OrderSn, Status: 1, Detail: details}
		if err := db.Create(&rec).Error; err != nil {
			return errors.InternalServer("SAVE_SELL_DETAIL_ERROR", err.Error())
		}
		return nil
	})
}
