package tran

import (
	"common/zerodb"
	"common/zerodb/gorms"
	"gorm.io/gorm"
)

type TransactionImpl struct {
	conn zerodb.DbConn
}

func NewTransaction(db *gorm.DB) *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.New(db),
	}
}

func (t *TransactionImpl) Action(f func(conn zerodb.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}
