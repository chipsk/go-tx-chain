package consts

// 司机状态	driver_collection table
// 1: 未加入, 2: 待生效, 3: 已加入
const (
	DriverNotJoinedState        = 1
	DriverWaitForEffectiveState = 2
	DriverJoinedState           = 3

	DriverCancelState = 6
)
