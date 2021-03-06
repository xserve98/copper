package csql

import (
	"context"

	"github.com/jinzhu/gorm"
)

const connCtxKey = "csql/tx"

// GetConn returns a db connection from the context or the given default connection if context is empty.
func GetConn(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(connCtxKey).(*gorm.DB)
	if !ok {
		return db
	}

	return tx
}
