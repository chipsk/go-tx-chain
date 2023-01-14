package endowment

import (
	"chipsk/go-tx-chain/biz/enum"
	"chipsk/go-tx-chain/biz/model"
	"context"
	"errors"
	"fmt"
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
