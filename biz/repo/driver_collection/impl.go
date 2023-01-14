package driver_collection

import (
	"chipsk/go-tx-chain/biz/enum"
	"chipsk/go-tx-chain/biz/model"
	dbErr "chipsk/go-tx-chain/common/error"
	"chipsk/go-tx-chain/common/mysql"
	"chipsk/go-tx-chain/common/util"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type impl struct {
	conn *gorm.DB
}

func (i *impl) getConn(ctx context.Context) (*gorm.DB, error) {
	if i.conn != nil {
		return i.conn, nil
	}
	conn, err := mysql.GetConn(ctx)
	if err != nil {
		fmt.Println("getConn err: ", err)
		return nil, err
	}
	return conn, err
}

func (i *impl) Get(ctx context.Context, driverId int64, itemId enum.InsuranceType) (*model.DriverCollection, error) {
	conn, err := i.getConn(ctx)
	if err != nil {
		return nil, err
	}
	model := &model.DriverCollection{
		Uid: driverId,
	}

	conn = conn.Table(model.TableName())
	db := conn.Where("uid = ?", driverId).Where("item_id = ?", itemId).First(model)

	if db.Error != nil && !dbErr.RecordNotFound(db) {
		fmt.Println("db error: ", err)
		return nil, db.Error
	}

	if dbErr.RecordNotFound(db) {
		fmt.Println("record_not_found: uid=", driverId)
		return nil, mysql.RecordNotFoundErr
	}

	return model, nil
}

func (i *impl) Create(ctx context.Context, dc *model.DriverCollection) (int64, error) {
	conn, err := i.getConn(ctx)
	if err != nil {
		return 0, err
	}

	conn = conn.Table(dc.TableName())
	db := conn.Create(dc)

	if db.Error != nil {
		fmt.Println("db error: ", err)
		return 0, db.Error
	}

	if db.RowsAffected == 0 {
		fmt.Println("err=record insert error")
		return 0, errors.New(fmt.Sprintf("rows affect = 0||%s", util.JsonString(dc)))
	}

	return dc.ID, nil
}

func (i *impl) Update(ctx context.Context, dc *model.DriverCollection) error {
	conn, err := i.getConn(ctx)
	if err != nil {
		return err
	}

	conn = conn.Table(dc.TableName()).Where("id = ?", dc.ID).Save(dc)
	db := conn.Create(dc)

	if db.Error != nil {
		fmt.Println("db error: ", err)
		return db.Error
	}

	if db.RowsAffected == 0 {
		fmt.Println("err=record update error")
		return errors.New(fmt.Sprintf("rows affect = 0||%s", util.JsonString(dc)))
	}

	return nil
}
