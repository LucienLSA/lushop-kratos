package data

import (
	"context"
	"time"

	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 定义数据表结构体
type User struct {
	ID          int64      `gorm:"primarykey"`
	Mobile      string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号码，用户唯一标识';not null"`
	Password    string     `gorm:"type:varchar(100);not null "` // 用户密码的保存需要注意是否加密
	NickName    string     `gorm:"type:varchar(25) comment '用户昵称'"`
	Birthday    *time.Time `gorm:"type:datetime comment '出生日日期'"`
	Gender      string     `gorm:"column:gender;default:male;type:varchar(16) comment 'female:女,male:男'"`
	Role        int        `gorm:"column:role;default:1;type:int comment '1:普通用户，2:管理员'"`
	CreatedAt   time.Time  `gorm:"column:add_time"`
	UpdatedAt   time.Time  `gorm:"column:update_time"`
	DeletedAt   gorm.DeletedAt
	IsDeletedAt bool
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

// CreateUser .
func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	var user User
	// 验证是否已经创建用户
	result := r.data.db.Where(&biz.User{Mobile: u.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = u.Mobile
	user.NickName = u.NickName
	user.Password = setPassword(u.Password)
	res := r.data.db.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
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
