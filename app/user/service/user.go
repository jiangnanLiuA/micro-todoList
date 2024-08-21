package service

import (
	"context"
	"errors"
	"micro-todoList/app/user/repository/dao"
	"micro-todoList/app/user/repository/model"
	"micro-todoList/idl/pb"
	"micro-todoList/pkg/e"

	"sync"

	"gorm.io/gorm"
)

var UserSrvIns *UserSrv   //userService实体
var UserSrvOnce sync.Once // 懒汉式单例模式  只执行一次

type UserSrv struct {
}

// GetUserSrv 懒汉式的单例模式 lazy-loading --> 懒汉式
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() { // 只执行一次
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// GetUserSrv 饿汉式的单例模式  lazy-loading --> 饿汉式  ->  会导致并发问题

func GetUserSrvHugury() *UserSrv {
	if UserSrvIns == nil {
		return new(UserSrv)
	}
	return UserSrvIns
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	user, err := dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		resp.Code = e.ERROR
		return
	}

	if !user.CheckPassword(req.Password) {
		resp.Code = e.InvalidParams
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		return
	}
	resp.Code = e.SUCCESS
	_, err = dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.ERROR
			return
		}
	}
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.ERROR
		return
	}
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.ERROR
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

// 序列化

func BuildUser(item *model.User) *pb.UserModel {
	userModel := pb.UserModel{
		Id:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return &userModel
}
