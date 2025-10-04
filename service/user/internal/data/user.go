package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 定义数据表结构体
type User struct {
	ID          int64      `gorm:"primarykey"`
	Mobile      string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号码，用户唯一标识';not null"`
	Password    string     `gorm:"type:varchar(100);not null "` // 用户密码的保存需要注意是否加密
	NickName    string     `gorm:"type:varchar(25) comment '用户昵称'"`
	Birthday    *time.Time `gorm:"type:datetime comment '出生日期'"`
	Gender      string     `gorm:"column:gender;default:male;type:varchar(16) comment 'female:女,male:男'"`
	Role        int        `gorm:"column:role;default:1;type:int comment '1:普通用户, 2:管理员'"`
	CreatedAt   time.Time  `gorm:"column:add_time"`
	UpdatedAt   time.Time  `gorm:"column:update_time"`
	DeletedAt   gorm.DeletedAt
	IsDeletedAt bool
}

func (User) TableName() string {
	return "user"
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ModelToResponse 转换 user 表中所有字段的值
func modelToResponse(user User) biz.User {
	userInfoRsp := biz.User{
		ID:        user.ID,
		Mobile:    user.Mobile,
		Password:  user.Password,
		NickName:  user.NickName,
		Gender:    user.Gender,
		Role:      user.Role,
		Birthday:  user.Birthday,
		CreatedAt: user.CreatedAt,
	}
	return userInfoRsp
}

// CreateUser .
func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	var user User
	// 验证是否已经创建用户
	result := r.data.db.Where(&biz.User{Mobile: u.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, errors.New(500, "USER_EXIST", "用户已存在"+u.Mobile)
	}
	user.Mobile = u.Mobile
	user.NickName = u.NickName
	user.Password = setPassword(u.Password)
	res := r.data.db.Create(&user)
	if res.Error != nil {
		return nil, errors.New(500, "CREAT_USER_ERROR", "用户创建失败")
	}
	return &biz.User{
		ID:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     user.Role,
	}, nil
}

const (
	PasswordCost = 12 //密码加密难度
)

// SetPassword 设置密码
func setPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return " "
	}
	// 将原密码进行加密
	encodePwd := string(bytes)
	return encodePwd
}

// paginate 分页
func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// ListUser .
func (r *userRepo) ListUser(ctx context.Context, pageNum, pageSize int) ([]*biz.User, int, error) {
	var users []User
	result := r.data.db.Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	if result.Error != nil {
		return nil, 0, errors.New(500, "FIND_USER_ERROR", "find user error")
	}
	total := int(result.RowsAffected)
	r.data.db.Scopes(paginate(pageNum, pageSize)).Find(&users)
	rv := make([]*biz.User, 0)
	for _, u := range users {
		rv = append(rv, &biz.User{
			ID:       u.ID,
			Mobile:   u.Mobile,
			Password: u.Password,
			NickName: u.NickName,
			Gender:   u.Gender,
			Role:     u.Role,
			Birthday: u.Birthday,
		})
	}
	return rv, total, nil
}

// GetUser .
func (r *userRepo) UserByMobile(ctx context.Context, mobile string) (*biz.User, error) {
	var user User
	result := r.data.db.Where(&User{Mobile: mobile}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	if result.Error != nil {
		return nil, errors.New(500, "FIND_USER_ERROR", "find user error")
	}
	if result.RowsAffected == 0 {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	res := modelToResponse(user)
	return &res, nil
}

// UpdateUser .
func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (bool, error) {
	var userInfo User
	result := r.data.db.Where(&User{ID: user.ID}).First(&userInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	if result.Error != nil {
		return false, errors.New(500, "FIND_USER_ERROR", "find user error")
	}
	if result.RowsAffected == 0 {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	userInfo.NickName = user.NickName
	userInfo.Birthday = user.Birthday
	userInfo.Gender = user.Gender
	res := r.data.db.Save(&userInfo)
	if res.Error != nil {
		return false, errors.InternalServer("USER_UPDATE_ERROR", "user save error")
	}
	return true, nil
}

// // CheckPassword 校验密码
// func checkPassword(password string) string {
// 	// 第一个参数是hashed已经加密过的密码，第二个参数请求中输入的参数密码
// 	fmt.Println(password)
// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// CheckPassword .
func (r *userRepo) CheckPassword(ctx context.Context, psd, encryptedPsd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPsd), []byte(psd))
	if err != nil {
		return false, errors.InternalServer("USER_PWD_ERROR", "user check pwd error")
	}
	return true, nil
}

// GetUserById .
func (r *userRepo) GetUserById(ctx context.Context, Id int64) (*biz.User, error) {
	var user User
	result := r.data.db.Where(&User{ID: Id}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	if result.Error != nil {
		return nil, errors.New(500, "FIND_USER_ERROR", "find user error")
	}
	if result.RowsAffected == 0 {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	res := modelToResponse(user)
	return &res, nil
}
