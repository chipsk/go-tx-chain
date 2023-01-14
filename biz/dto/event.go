package dto

type PolicyInfo struct {
	ChangeType   string `json:"change_type"`   // 变动类型
	PolicyCode   string `json:"policy_code"`   // 保司保单号码
	PolicyStatus string `json:"policy_status"` // 保单状态
	EndCause     string `json:"end_cause"`     // 终止原因
}

type CancelParam struct {
	DriverId   int64       `json:"driver_id"`
	Policy     *PolicyInfo `json:"policy"`
	PolicyCode string      `json:"policy_code"`
}

type ExitParam struct {
	CancelParam
}
