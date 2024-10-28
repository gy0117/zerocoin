package tran

import (
	"gorm.io/gorm"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
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
