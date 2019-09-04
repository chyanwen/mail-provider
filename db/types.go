package db

type Strategy struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       string            `json:"tags"`
	MaxStep    int               `json:"max_step"`
	Priority   int               `json:"priority"`
	Func       string            `json:"func"`       
	Operator   string            `json:"op"`   
	RightValue string           `json:"right_value"`
	RunBegin   string            `json:"run_begin"`
	RunEnd     string            `json:"run_end"`
	Note       string            `json:"note"`
	TplId      int               `json:"tpl_id"`
}

type AlarmVerification struct {
	Id              int         `json:"id"`
	StrategyId      int         `json:"strategy_id"`
	StrategyCopyId  int         `json:"strategy_copyid"`
	CreatedTime     int         `json:"createdtime"`
	VerificationStatus  int     `json:"verification_status"`
}
