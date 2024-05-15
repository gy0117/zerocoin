package repo

import (
	"context"
	"ucenter-rpc/internal/model"
	"zero-common/zerodb"
)

type WalletRepo interface {
	FindByCoinName(ctx context.Context, userId int64, coinName string) (*model.UserWallet, error)
	FindByCoinId(ctx context.Context, userId int64, coinId int64) (*model.UserWallet, error)
	FindByUserId(ctx context.Context, userId int64) ([]*model.UserWallet, error)
	Save(ctx context.Context, walletData *model.UserWallet) error
	Freeze(ctx context.Context, userId int64, money float64, symbol string) error
	FreezeWithConn(ctx context.Context, conn zerodb.DbConn, userId int64, money float64, symbol string) error
	UpdateWallet(ctx context.Context, conn zerodb.DbConn, id int64, balance float64, frozenBalance float64) error
	UpdateAddress(ctx context.Context, wallet *model.UserWallet) error
	GetAllAddress(ctx context.Context, coinName string) ([]string, error)
	GetByAddress(ctx context.Context, address string) (*model.UserWallet, error)
}
