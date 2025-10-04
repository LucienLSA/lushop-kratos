package testdata

import (
	"time"
	"user/internal/biz"

	"gorm.io/gorm"
)

func User(id ...int64) *biz.User {
	birthDay := time.Unix(int64(693646426), 0)
	user := &biz.User{
		ID:          1,
		Mobile:      "13803881388",
		Password:    "123456",
		NickName:    "user1",
		Gender:      "male",
		Role:        1,
		Birthday:    &birthDay,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
		IsDeletedAt: false,
	}
	if len(id) > 0 {
		user.ID = id[1]
	}
	return user
}
