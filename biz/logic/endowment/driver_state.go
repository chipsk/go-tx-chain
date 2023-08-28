package endowment

import (
	"chipsk/go-tx-chain/biz/enum"
	"chipsk/go-tx-chain/biz/model"
	"chipsk/go-tx-chain/biz/repo/driver_collection"
	"chipsk/go-tx-chain/common/mysql"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type transferMap map[enum.EndowmentEvent]enum.DriverState

var driverStateMap = map[enum.DriverState]transferMap{

	enum.DriverNotJoinedState: { //1 未加入
		enum.EndowmentJoinEvent:  enum.DriverNotJoinedState,
		enum.EndowmentApplyEvent: enum.DriverWaitForEffectiveState, // 投保成功事件: 待生效
	},

	enum.DriverWaitForEffectiveState: { //2 待生效

		enum.EndowmentCancelEvent: enum.DriverJoinedState,
	},
}

type DriverStateService struct {
	DriverCollection *model.DriverCollection `json:"driver_collection"`
}

func (d *DriverStateService) GetDst(event enum.EndowmentEvent) (enum.DriverState, error) {
	transfer, _ := driverStateMap[enum.DriverState(d.DriverCollection.Status)]
	for e, dstState := range transfer {
		if event == e {
			return dstState, nil
		}
	}
	return enum.DriverState(-1), errors.New(fmt.Sprintf("driver_state: %s can not do %s action, driver_id: %d", enum.DriverState(d.DriverCollection.Status).Desc(), d.DriverCollection.Uid))
}

func (d *DriverStateService) Can(event enum.EndowmentEvent) (bool, error) {
	_, err := d.GetDst(event)
	return err == nil, err
}

func initEndowmentDriverCollection(ctx context.Context, driverId int64) (*model.DriverCollection, error) {

	return &model.DriverCollection{
		Uid:    driverId,
		Status: int(enum.DriverNotJoinedState),
	}, nil
}

func NewDriverState(ctx context.Context, driverId int64, driverCollection *model.DriverCollection) (svr *DriverStateService, err error) {
	if driverCollection == nil {
		driverCollection, err = initEndowmentDriverCollection(ctx, driverId)
		if err != nil {
			return nil, err
		}
	}
	return &DriverStateService{
		DriverCollection: driverCollection,
	}, nil
}

func (d *DriverStateService) changeStatus(ctx context.Context, event enum.EndowmentEvent) (mysql.TxFunc, error) {
	dst, err := d.GetDst(event)
	if err != nil {
		return nil, err
	}

	f := func(ctx context.Context, tx *gorm.DB) error {
		d.DriverCollection.Status = dst.Int()
		return driver_collection.NewWithTx(tx).Update(ctx, d.DriverCollection)
	}

	return f, nil
}

// Apply
// 如果司机状态信息不存在则新增
// 如果司机状态信息存在，且为未加入计划状态，则更新为加入计划
// 其他不做处理
func (d *DriverStateService) Apply(ctx context.Context) (mysql.TxFunc, error) {
	dst, err := d.GetDst(enum.EndowmentApplyEvent)
	if err != nil {
		return nil, err
	}

	f := func(ctx context.Context, tx *gorm.DB) error {
		d.DriverCollection.Status = dst.Int()
		return driver_collection.NewWithTx(tx).Update(ctx, d.DriverCollection)
	}
	return f, nil
}

func (d *DriverStateService) Join(ctx context.Context) (mysql.TxFunc, error) {
	dst, err := d.GetDst(enum.EndowmentJoinEvent)
	if err != nil {
		return nil, err
	}

	// todo
	fmt.Println(dst)
	return nil, err
}

func (d *DriverStateService) ClaimStop(ctx context.Context) (mysql.TxFunc, error) {
	return d.changeStatus(ctx, enum.EndowmentClaimStopEvent)
}

func (d *DriverStateService) Cancel(ctx context.Context) (mysql.TxFunc, error) {
	return d.changeStatus(ctx, enum.EndowmentCancelEvent)
}

func (d *DriverStateService) Exit(ctx context.Context) (mysql.TxFunc, error) {
	return d.changeStatus(ctx, enum.EndowmentExitEvent)
}
