package endowment

import (
	"chipsk/go-tx-chain/biz/dto"
	"chipsk/go-tx-chain/biz/enum"
	"chipsk/go-tx-chain/biz/repo/driver_collection"
	"chipsk/go-tx-chain/common/mysql"
	"context"
	"fmt"
)

type stateEnum int

const (
	driverStateEnum stateEnum = 1 << iota // driver_state_enum
	orderStateEnum                        // order_state_enum
	policyStateEnum                       // policy_state_enum
	billStateEnum                         // bill_state_enum
)

type EventService struct {
}

func NewEventService() (*EventService, error) {
	return &EventService{}, nil
}

/**
event的执行分为
1. 当前状态是否可以执行Event(是否有对应handler)
2. 执行前状态机获取(获取driver/bill/order/policy等)
3. 执行逻辑检查, 每个状态机执行事务调用
4. 具体执行(事务, 三方调用)
	a. 事务的产生应该由各个状态机负责
	b. 状态切换导致的三方调用可以在事务执行结束后完成(当前产生的三方调用比较少, 所以尝试如此收敛, 具体执行灵活判断)
5. 错误处理
*/

func (s *EventService) Cancel(ctx context.Context, param *dto.CancelParam) error {
	dState, oState, pState, bState, err := s.initState(ctx, param.DriverId, enum.EndowmentCancelEvent, billStateEnum, policyStateEnum, orderStateEnum, driverStateEnum)

}

func (s *EventService) initState(ctx context.Context, driverId int64, event enum.EndowmentEvent, states ...stateEnum) (driverState *DriverStateService, orderState *OrderStateService, policyState *PolicyStateService, billState *BillStateService, err error) {
	st := stateEnum(0)
	for _, state := range states {
		st = st | state
	}

	// 状态机是按照 driver -> order -> policy -> bill 自上而下调用
	// order -> 保单 -> 账单
	v := make([]func(), 0, 4)
	deferFunc := func(f func()) {
		v = append(v, f)
	}

	defer func() {
		// 出栈过程
		for len(v) > 0 {
			v[len(v)-1]()
			if err != nil {
				fmt.Println("index:initState.init||err: ", err)
				return
			}
			v = v[:len(v)-1]
		}
	}()

	switch {
	case st&billStateEnum > 0:
		deferFunc(func() {
			billState, err = s.getBillState(ctx, driverId, policyState.Policy.PolicyID)
		})
		fallthrough
	case st&policyStateEnum > 0:
		deferFunc(func() {
			policyState, err = s.getPolicyState(ctx, driverId, orderState.Order.OrderId, event)
		})
		fallthrough
	case st&orderStateEnum > 0:
		deferFunc(func() {
			orderState, err = s.getOrderState(ctx, driverId, event)
		})
	}
	if st&driverStateEnum > 0 {
		deferFunc(func() {
			driverState, err = s.getDriverState(ctx, driverId, event)
		})
	}
	return
}

func (s *EventService) getDriverState(ctx context.Context, driverId int64, event enum.EndowmentEvent) (*DriverStateService, error) {
	dcd, err := driver_collection.New().Get(ctx, driverId, enum.DriverEndowmentInsurance)
	if err != nil && err != mysql.RecordNotFoundErr {
		return nil, err
	}

	// get driver state
	driverState, err := NewDriverState(ctx, driverId, dcd)
	if err != nil {
		return nil, err
	}

	// load driver state
	ok, err := driverState.Can(event)
	if !ok {
		return nil, err
	}

	return driverState, nil

}
