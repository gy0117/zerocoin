package tran

import "zero-common/zerodb"

type Transaction interface {
	Action(func(conn zerodb.DbConn) error) error
}
