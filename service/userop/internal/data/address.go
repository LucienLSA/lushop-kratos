package data

import (
	"userop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type addressRepo struct {
	data *Data
	log  *log.Helper
}

// NewAddressRepo .
func NewAddressRepo(data *Data, logger log.Logger) biz.AddressRepo {
	return &addressRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
