package repo

import (
	"context"
	"ucenter-rpc/internal/model"
)

type UserAddressRepo interface {
	Save(ctx context.Context, address *model.UserAddress) error
	FindByCoinId(ctx context.Context, userId int64, coinId int64) (list []*model.UserAddress, err error)
	FindByUserId(ctx context.Context, userId int64) (list []*model.UserAddress, err error)
}
