package object

import "time"

type TaskInfo struct {
	TaskId     string    `json:"taskId"`
	CreateTime time.Time `json:"create_time"` //创建时间
}
