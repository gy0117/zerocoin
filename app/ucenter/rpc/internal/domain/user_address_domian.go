package domain

import (
	"context"
	"errors"
	"ucenter-rpc/internal/dao"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
)

type UserAddressDomain struct {
	mAddressRepo repo.UserAddressRepo
}

func NewUserAddressDomain(db *zerodb.ZeroDB) *UserAddressDomain {
	return &UserAddressDomain{
		mAddressRepo: dao.NewUserAddressDao(db),
	}
}

func (d *UserAddressDomain) FindAddressesByCoinId(ctx context.Context, userId int64, coinId int64) ([]*model.UserAddress, error) {
	list, err := d.mAddressRepo.FindByCoinId(ctx, userId, coinId)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, errors.New("no record")
	}
	return list, nil
}
