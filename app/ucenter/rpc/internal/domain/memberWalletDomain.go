package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"grpc-common/market/types/market"
	"ucenter-rpc/internal/dao"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/tran"
)

type WalletDomain struct {
	walletRepo  repo.WalletRepo
	transaction tran.Transaction
}

func NewWalletDomain(db *zerodb.ZeroDB) *WalletDomain {
	return &WalletDomain{
		walletRepo:  dao.NewWalletDao(db),
		transaction: tran.NewTransaction(db.Conn),
	}
}

// FindWallet 返回的是数据库model结构体
func (d *WalletDomain) FindWallet(ctx context.Context, userId int64, coinName string, coin *market.Coin) (*model.UserWalletCoin, error) {
	// 1. 找到userId用户拥有coinName的数据
	walletData, err := d.walletRepo.FindByCoinName(ctx, userId, coinName)
	if err != nil {
		return nil, err
	}
	if walletData == nil {
		// 2. 如果不存在匹配的数据，则创建默认数据保存
		memberWallet, memberWalletCoin := model.NewWalletData(userId, coin)
		err = d.walletRepo.Save(ctx, memberWallet)
		if err != nil {
			return nil, err
		}
		return memberWalletCoin, nil
	}

	mwc := &model.UserWalletCoin{}
	if err = copier.Copy(mwc, walletData); err != nil {
		return nil, err
	}
	mwc.Coin = coin
	return mwc, nil
}

func (d *WalletDomain) Freeze(ctx context.Context, userId int64, money float64, symbol string) error {
	// 冻结之前，查看钱包的钱是否足够
	wallet, err := d.walletRepo.FindByCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	if wallet.Balance < money {
		return errors.New("insufficient balance")
	}
	return d.walletRepo.Freeze(ctx, userId, money, symbol)
}

func (d *WalletDomain) FreezeWithConn(ctx context.Context, conn zerodb.DbConn, userId int64, money float64, symbol string) error {
	// 冻结之前，查看钱包的钱是否足够
	wallet, err := d.walletRepo.FindByCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	fmt.Printf("ucenter | Freeze | wallet: %+v\n", wallet)
	if wallet.Balance < money {
		return errors.New("余额不足-3")
	}
	return d.walletRepo.FreezeWithConn(ctx, conn, userId, money, symbol)
}

func (d *WalletDomain) FindWalletByMemIdAndCoinName(ctx context.Context, userId int64, symbol string) (*model.UserWallet, error) {
	wallet, err := d.walletRepo.FindByCoinName(ctx, userId, symbol)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("no wallet")
	}
	return wallet, nil
}

func (d *WalletDomain) FindWalletByCoinId(ctx context.Context, userId, coinId int64) (*model.UserWallet, error) {
	wallet, err := d.walletRepo.FindByCoinId(ctx, userId, coinId)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("no wallet")
	}
	return wallet, nil
}

func (d *WalletDomain) UpdateWalletCoinAndBase(ctx context.Context, baseWallet *model.UserWallet, coinWallet *model.UserWallet) error {
	// 同时更新成功
	return d.transaction.Action(func(conn zerodb.DbConn) error {
		err := d.walletRepo.UpdateWallet(ctx, conn, baseWallet.Id, baseWallet.Balance, baseWallet.FrozenBalance)
		if err != nil {
			return err
		}
		err = d.walletRepo.UpdateWallet(ctx, conn, coinWallet.Id, coinWallet.Balance, coinWallet.FrozenBalance)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *WalletDomain) FindWalletsByUserId(ctx context.Context, userId int64) ([]*model.UserWallet, error) {
	memberWallets, err := d.walletRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if memberWallets == nil {
		return nil, errors.New("record is not exist")
	}
	return memberWallets, err
}

func (d *WalletDomain) UpdateWalletAddress(ctx context.Context, wallet *model.UserWallet) error {
	return d.walletRepo.UpdateAddress(ctx, wallet)
}

func (d *WalletDomain) GetAddress(ctx context.Context, coinName string) ([]string, error) {
	allAddress, err := d.walletRepo.GetAllAddress(ctx, coinName)
	if err != nil {
		return nil, err
	}
	if allAddress == nil {
		return nil, errors.New("没有记录")
	}
	return allAddress, nil
}

func (d *WalletDomain) GetByAddress(ctx context.Context, address string) (*model.UserWallet, error) {
	memberWallet, err := d.walletRepo.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if memberWallet == nil {
		return nil, errors.New("no record")
	}
	return memberWallet, nil
}
