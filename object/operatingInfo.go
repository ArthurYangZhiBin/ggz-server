package object

import (
	"time"
)

//操作表只能运维修改
type OperatingInfo struct {
	FlowID        string    `json:"flowID"`
	AprovUserId   string    `json:"aprovUserId"`
	Username      string    `json:"username"`
	XconfBefore   string    `json:"xconfBefore"`
	XconfAfter    string    `json:"xconfAfter"`
	ComposeBefore string    `json:"composeBefore"`
	ComposeAfter  string    `json:"composeAfter"`
	CreateTime    time.Time `json:"create_time"` //创建时间
}
type OperatingInfos struct {
	//普通struct类型
	OperatingInfo []OperatingInfo `json:"operatingInfo"`
}
