package enum

type EndowmentEvent string

const (
	// EndowmentApplyEvent 1.投保成功通知
	EndowmentApplyEvent EndowmentEvent = "EndowmentApplyEvent" // apply 投保成功

	// EndowmentReconciliationOrderEvent 2.日交易对账
	EndowmentReconciliationOrderEvent EndowmentEvent = "EndowmentReconciliationOrderEvent"

	// 3.司机状态变动通知
	EndowmentValidEvent     EndowmentEvent = "EndowmentValidEvent"     // valid 保单生效
	EndowmentCancelEvent    EndowmentEvent = "EndowmentCancelEvent"    // cancel 犹豫期扯单变动
	EndowmentExitEvent      EndowmentEvent = "EndowmentExitEvent"      // exit 退保变动
	EndowmentClaimStopEvent EndowmentEvent = "EndowmentClaimStopEvent" // claim stop 理赔终止

	// EndowmentJoinEvent 用户手动事件
	EndowmentJoinEvent EndowmentEvent = "EndowmentJoinEvent" // join
)
