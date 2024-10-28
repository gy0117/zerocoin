package dao

import (
	"context"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/repo"
	"gorm.io/gorm"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.OrderRepo = (*OrderDao)(nil)

type OrderDao struct {
	conn *gorms.GormConn
}

func (dao *OrderDao) UpdateOrderComplete(ctx context.Context, orderId string, tradedAmount float64, turnover float64) error {
	session := dao.conn.Session(ctx)
	sql := "UPDATE exchange_order SET traded_amount=?, turnover=?, status=? WHERE order_id=?"
	return session.Model(&model.ExchangeOrder{}).Exec(sql, tradedAmount, turnover, model.OrderStatus_Completed, orderId).Error
}

// FindTradingOrderBySymbol 查询正在交易的订单，且满足symbol条件
func (dao *OrderDao) FindTradingOrderBySymbol(ctx context.Context, symbol string) (list []*model.ExchangeOrder, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("symbol=? and status=?", symbol, model.OrderStatus_Trading).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (dao *OrderDao) UpdateOrderStatus(ctx context.Context, orderId string, status int) error {
	session := dao.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Update("status", status).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (dao *OrderDao) CancelOrder(ctx context.Context, orderId string) error {
	session := dao.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Update("status", model.OrderStatus_Canceled).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (dao *OrderDao) FindByOrderId(ctx context.Context, orderId string) (eo *model.ExchangeOrder, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Take(&eo).Error
	// 将无记录和err区分开来，在dao层处理好，然后在domain层进一步判断
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (dao *OrderDao) Save(ctx context.Context, conn zerodb.DbConn, order *model.ExchangeOrder) error {
	txConn := conn.(*gorms.GormConn)
	tx := txConn.Tx(ctx)
	return tx.Save(&order).Error
}

func (dao *OrderDao) FindCurrentTradingCount(ctx context.Context, memId int64, symbol string, direction int) (total int64, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("user_id=? and symbol=? and direction=? and status=?", memId, symbol, direction, model.OrderStatus_Trading).Count(&total).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return
}

func (dao *OrderDao) GetHistoryOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := dao.conn.Session(ctx)
	// 找到第pageNo页的数据，一页的数据量为pageSize
	err = session.Model(&model.ExchangeOrder{}).Where("symbol=? and user_id=?", symbol, memId).Limit(int(pageSize)).Offset(int((pageNo - 1) * pageSize)).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	err = session.Model(&model.ExchangeOrder{}).Where("symbol=? and user_id=?", symbol, memId).Count(&total).Error
	return
}

func (dao *OrderDao) GetCurrentOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := dao.conn.Session(ctx)
	// 找到第pageNo页的数据，一页的数据量为pageSize
	err = session.Model(&model.ExchangeOrder{}).Where("symbol=? and user_id=? and status=?", symbol, memId, model.OrderStatus_Trading).Limit(int(pageSize)).Offset(int((pageNo - 1) * pageSize)).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	err = session.Model(&model.ExchangeOrder{}).Where("symbol=? and user_id=?  and status=?", symbol, memId, model.OrderStatus_Trading).Count(&total).Error
	return
}

func NewOrderDao(db *zerodb.ZeroDB) *OrderDao {
	return &OrderDao{
		conn: gorms.New(db.Conn),
	}
}
