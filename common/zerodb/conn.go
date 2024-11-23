package zerodb

import "gorm.io/gorm"

type DbConn interface {
	Begin()
	Rollback()
	Commit()
}

type ZeroDB struct {
	Conn *gorm.DB
}
