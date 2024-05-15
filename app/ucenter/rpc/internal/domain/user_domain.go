package domain

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter-rpc/internal/dao"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/tools"
	"zero-common/zerodb"
)

type UserDomain struct {
	userRepo repo.UserRepo
}

func NewUserDomain(db *zerodb.ZeroDB) *UserDomain {
	return &UserDomain{
		userRepo: dao.NewUserMemberDao(db),
	}
}

func (d *UserDomain) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	um, err := d.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("database Exception")
	}
	return um, nil
}

func (d *UserDomain) FindByUserName(ctx context.Context, username string) (*model.User, error) {
	um, err := d.userRepo.FindByUserName(ctx, username)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("database Exception")
	}
	return um, nil
}

func (d *UserDomain) FindUserById(ctx context.Context, userId int64) (*model.User, error) {
	um, err := d.userRepo.FindUserById(ctx, userId)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("database Exception")
	}
	if um == nil {
		return nil, errors.New("not found record")
	}
	return um, nil
}

func (d *UserDomain) Register(ctx context.Context, username, phone, password, country, promotion, partner string) error {
	// password 进行md5加密，同时加盐，md5加密不安全（通过彩虹表进行破解）
	newUser := model.NewUser()
	// 设置默认值
	if err := tools.Default(newUser); err != nil {
		return err
	}
	// 赋值
	newUser.Username = username
	newUser.MobilePhone = phone
	newUser.Country = country
	newUser.PromotionCode = promotion
	newUser.FillSuperPartner(partner)
	newUser.UserLevel = model.GENERAL
	newUser.Avatar = "https://p2.itc.cn/q_70/images03/20230902/721191166cd2465c9db74d5b52a3e7bc.png"

	salt, newPwd := tools.Encode(password, nil)
	newUser.Password = newPwd
	newUser.Salt = salt

	// 保存到数据库
	err := d.userRepo.Save(ctx, newUser)
	if err != nil {
		logx.Error(err)
		return errors.New("failed to register, err is " + err.Error())
	}
	return nil
}
