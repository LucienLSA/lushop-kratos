package data

import (
	"inventory/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

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
