package enum

type EndowmentEvent string

const (
	EndowmentJoinEvent   EndowmentEvent = "EndowmentJoinEvent"   // join
	EndowmentApplyEvent  EndowmentEvent = "EndowmentApplyEvent"  // apply
	EndowmentCancelEvent EndowmentEvent = "EndowmentCancelEvent" // cancel
	EndowmentExitEvent   EndowmentEvent = "EndowmentExitEvent"   // exit
)
