package Initialize

// ActionDescribeTriggerLog 后端数据结构
type ActionDescribeTriggerLog struct {
	DescribeMessage string `json:"DescribeMessage"`
	ActionType      int8   `json:"ActionType"`
	SrcAddr         string `json:"SrcAddr"`
	DstAddr         string `json:"DstAddr"`
}
