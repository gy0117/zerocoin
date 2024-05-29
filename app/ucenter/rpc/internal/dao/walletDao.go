package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.WalletRepo = (*WalletDao)(nil)

type WalletDao struct {
	conn *gorms.GormConn
}

func (w *WalletDao) GetByAddress(ctx context.Context, address string) (mw *model.UserWallet, err error) {
	session := w.conn.Session(ctx)
	err = session.Model(&model.UserWallet{}).Where("address=?", address).Take(&mw).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (w *WalletDao) GetAllAddress(ctx context.Context, coinName string) (list []string, err error) {
	session := w.conn.Session(ctx)
	err = session.Model(&model.UserWallet{}).Where("coin_name=?", coinName).Select("address").Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (w *WalletDao) UpdateAddress(ctx context.Context, wallet *model.UserWallet) error {
	session := w.conn.Session(ctx)
	sql := "UPDATE user_wallet SET address=?, address_private_key=? WHERE id=?"
	return session.Model(&model.UserWallet{}).Exec(sql, wallet.Address, wallet.AddressPrivateKey, wallet.Id).Error
}

func (w *WalletDao) FindByUserId(ctx context.Context, userId int64) (list []*model.UserWallet, err error) {
	session := w.conn.Session(ctx)
	err = session.Model(&model.UserWallet{}).Where("user_id=?", userId).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (w *WalletDao) UpdateWallet(ctx context.Context, conn zerodb.DbConn, id int64, balance float64, frozenBalance float64) error {
	txConn := conn.(*gorms.GormConn)
	tx := txConn.Tx(ctx)
	sql := "UPDATE user_wallet SET balance=?, frozen_balance=? WHERE id=?"
	err := tx.Model(&model.UserWallet{}).Where("id=?", id).Exec(sql, balance, frozenBalance, id).Error
	return err
}

func (w *WalletDao) Freeze(ctx context.Context, userId int64, money float64, symbol string) error {
	session := w.conn.Session(ctx)
	sql := "update user_wallet set balance=balance-?, frozen_balance=frozen_balance+? where user_id=? and coin_name=?"
	exec := session.Model(&model.UserWallet{}).Exec(sql, money, money, userId, symbol)
	err := exec.Error
	if err != nil {
		return err
	}
	if exec.RowsAffected <= 0 {
		return errors.New("there is no data to update")
	}
	return nil
}

func (w *WalletDao) FreezeWithConn(ctx context.Context, conn zerodb.DbConn, userId int64, money float64, symbol string) error {
	txConn := conn.(*gorms.GormConn)
	tx := txConn.Tx(ctx)

	sql := "update user_wallet set balance=balance-?, frozen_balance=frozen_balance+? where user_id=? and coin_name=?"
	return tx.Model(&model.UserWallet{}).Exec(sql, money, money, userId, symbol).Error
}

func (w *WalletDao) FindByCoinName(ctx context.Context, userId int64, coinName string) (data *model.UserWallet, err error) {
	session := w.conn.Session(ctx)
	err = session.Model(&model.UserWallet{}).Where("user_id=? and coin_name=?", userId, coinName).Take(&data).Error
	if err == gorm.ErrRecordNotFound { // 这个不算错误
		return nil, nil
	}
	return
}

func (w *WalletDao) FindByCoinId(ctx context.Context, userId int64, coinId int64) (data *model.UserWallet, err error) {
	session := w.conn.Session(ctx)
	err = session.Model(&model.UserWallet{}).Where("user_id=? and coin_id=?", userId, coinId).Take(&data).Error
	if err == gorm.ErrRecordNotFound { // 这个不算错误
		return nil, nil
	}
	return
}

func (w *WalletDao) Save(ctx context.Context, walletData *model.UserWallet) error {
	session := w.conn.Session(ctx)
	return session.Save(&walletData).Error
}

func NewWalletDao(db *zerodb.ZeroDB) repo.WalletRepo {
	return &WalletDao{
		conn: gorms.New(db.Conn),
	}
}
