package dao

import (
	"context"
	"micro_todoList/app/user/repository/model"

	"gorm.io/gorm"
)

// 定义的是对数据库的user model 的curd操作

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUserName(userName string) (user *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name = ?", userName).Find(&user).Error
	// record not found
	// select * from user where user_name = xxx order by id limit 1;  -- First 会报  record not found
	// select * from user where user_name = xxx;   -- Find  不会报  record not found
	if err != nil {
		return
	}
	return
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(in *model.User) (err error) {
	return dao.Model(&model.User{}).Create(&in).Error
}
