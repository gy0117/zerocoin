package tran

import "common/zerodb"

type Transaction interface {
	Action(func(conn zerodb.DbConn) error) error
}
