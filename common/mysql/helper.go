package mysql

import (
	"chipsk/go-tx-chain/conf"
	"context"
	"gorm.io/gorm"
)

func GetConn(ctx context.Context) (*gorm.DB, error) {
	db, err := Client.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	return db, db.Error
}

func getTxConn(ctx context.Context, option ...Option) (*gorm.DB, error) {
	db, err := Client.GetDB(ctx, option...)
	if err != nil {
		return nil, err
	}
	// begin transaction
	tx := db.Begin()
	return tx, tx.Error
}

func TableIndexByUid(uid int64) int64 {
	shards := conf.Viper.GetInt64("mysql.table_default_shards_num")
	index := uid % shards
	return index
}
