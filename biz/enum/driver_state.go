package enum

import (
	"chipsk/go-tx-chain/common/consts"
	"fmt"
)

type DriverState int

const (
	DriverNotJoinedState        DriverState = consts.DriverNotJoinedState        //1 未加入
	DriverWaitForEffectiveState DriverState = consts.DriverWaitForEffectiveState //2 待生效
	DriverJoinedState           DriverState = consts.DriverJoinedState           //3 已加入

	DriverCancelState DriverState = consts.DriverCancelState //6 已退保

)

func (d DriverState) Desc() string {
	m := map[DriverState]string{

		DriverNotJoinedState:        "未加入",
		DriverWaitForEffectiveState: "待生效",
		DriverJoinedState:           "已加入",

		DriverCancelState: "已退保",
	}

	desc, ok := m[d]
	if !ok {
		return fmt.Sprintf("unknown: %d", d)
	}
	return desc
}

func (d DriverState) Int() int {
	return int(d)
}
