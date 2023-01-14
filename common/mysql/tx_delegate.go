package mysql

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"runtime/debug"
)

//用于事务委托

type TxFunc func(ctx context.Context, txConn *gorm.DB) error

type TxFuncChain []TxFunc

func (tf TxFunc) Do(ctx context.Context, option ...Option) error {
	list := TxFuncChain([]TxFunc{tf})
	return list.Do(ctx, option...)
}

func (tf TxFuncChain) Do(ctx context.Context, option ...Option) error {
	tx, err := getTxConn(ctx, option...)
	if err != nil {
		// log error
		return err
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf(string(debug.Stack()))
			// log error
			tx.Rollback()
		}
		return
	}()

	for _, txFunc := range tf {
		err := txFunc(ctx, tx)
		if err != nil {
			// log error
			tx.Rollback()
			return err
		}
	}
	tx = tx.Commit()
	return tx.Error
}
