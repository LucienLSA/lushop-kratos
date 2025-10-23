package data

import (
	"context"
	"userop/internal/biz"
	"userop/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	LEAVING_MESSAGES = iota + 1
	COMPLAINT
	INQUIRY
	POST_SALE
	WANT_TO_BUY
)

type LeavingMessages struct {
	BaseFields
	User        int32  `gorm:"type:int;index;comment:用户ID"`
	MessageType int32  `gorm:"type:int;comment:留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)"`
	Subject     string `gorm:"type:varchar(100);comment:主题"`

	Message string `gorm:"comment:详细信息"`
	File    string `gorm:"type:varchar(200);comment:附件url"`
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

type Address struct {
	BaseFields
	User         int32  `gorm:"type:int;index;comment:用户id"`
	Province     string `gorm:"type:varchar(10);comment:省"`
	City         string `gorm:"type:varchar(10);comment:市"`
	District     string `gorm:"type:varchar(20);comment:区域"`
	Address      string `gorm:"type:varchar(100);comment:详细地址"`
	SignerName   string `gorm:"type:varchar(20);comment:收货人名称"`
	SignerMobile string `gorm:"type:varchar(11);comment:收货人手机号"`
}

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

func (f *addressRepo) GetAddressList(ctx context.Context, filter domain.Address) (*domain.AddressListResponse, error) {
	var (
		list  []domain.Address
		total int64
	)
	db := f.data.DB(ctx).Model(&domain.Address{})
	if filter.ID != 0 {
		db = db.Where("id = ?", filter.ID)
	}
	if filter.UserID != 0 {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Province != "" {
		db = db.Where("province = ?", filter.Province)
	}
	if filter.City != "" {
		db = db.Where("city = ?", filter.City)
	}
	if filter.District != "" {
		db = db.Where("district = ?", filter.District)
	}
	if filter.Address != "" {
		db = db.Where("address = ?", filter.Address)
	}
	if filter.SignerName != "" {
		db = db.Where("signer_name = ?", filter.SignerName)
	}
	if filter.SignerMobile != "" {
		db = db.Where("signer_mobile = ?", filter.SignerMobile)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	resp := &domain.AddressListResponse{Total: total}
	for i := range list {
		// take address of each element to append pointers
		addr := list[i]
		resp.List = append(resp.List, &addr)
	}
	return resp, nil
}

func (f *addressRepo) CreateAddress(ctx context.Context, address domain.Address) error {
	return f.data.DB(ctx).Create(&address).Error
}

func (f *addressRepo) DeleteAddress(ctx context.Context, address domain.Address) error {
	// Prefer delete by ID if provided; otherwise use other non-zero fields
	db := f.data.DB(ctx)
	if address.ID != 0 {
		return db.Where("id = ?", address.ID).Delete(&domain.Address{}).Error
	}
	return db.Where(&address).Delete(&domain.Address{}).Error
}

func (f *addressRepo) UpdateAddress(ctx context.Context, address domain.Address) error {
	if address.ID == 0 {
		return f.data.DB(ctx).Save(&address).Error
	}
	return f.data.DB(ctx).Model(&domain.Address{}).Where("id = ?", address.ID).Updates(address).Error
}

func (f *addressRepo) GetMessageList(ctx context.Context, filter domain.Message) (*domain.MessageListResponse, error) {
	var (
		list  []LeavingMessages
		total int64
	)
	db := f.data.DB(ctx).Model(&LeavingMessages{})
	if filter.ID != 0 {
		db = db.Where("id = ?", filter.ID)
	}
	if filter.UserID != 0 {
		db = db.Where("user = ?", filter.UserID)
	}
	if filter.MessageType != 0 {
		db = db.Where("message_type = ?", filter.MessageType)
	}
	if filter.Subject != "" {
		db = db.Where("subject = ?", filter.Subject)
	}
	if filter.Message != "" {
		db = db.Where("message = ?", filter.Message)
	}
	if filter.File != "" {
		db = db.Where("file = ?", filter.File)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	resp := &domain.MessageListResponse{Total: total}
	for i := range list {
		m := list[i]
		resp.List = append(resp.List, &domain.Message{
			ID:          int32(m.ID),
			UserID:      m.User,
			MessageType: m.MessageType,
			Subject:     m.Subject,
			Message:     m.Message,
			File:        m.File,
		})
	}
	return resp, nil
}

func (f *addressRepo) CreateMessage(ctx context.Context, msg domain.Message) error {
	m := LeavingMessages{
		User:        msg.UserID,
		MessageType: msg.MessageType,
		Subject:     msg.Subject,
		Message:     msg.Message,
		File:        msg.File,
	}
	return f.data.DB(ctx).Create(&m).Error
}
