package driver_collection

import (
	"chipsk/go-tx-chain/biz/enum"
	"chipsk/go-tx-chain/biz/model"
	"context"
	"gorm.io/gorm"
)

type Interface interface {
	Get(ctx context.Context, driverId int64, itemId enum.InsuranceType) (*model.DriverCollection, error)
	Create(ctx context.Context, dc *model.DriverCollection) (int64, error)
	Update(ctx context.Context, dc *model.DriverCollection) error
}

func New() Interface {
	return &impl{}
}

func NewWithTx(db *gorm.DB) Interface {
	return &impl{conn: db}
}
