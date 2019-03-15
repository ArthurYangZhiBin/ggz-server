package main

import (
	"ggz-server/handler"
	"ggz-server/object"
	"ggz-server/store"
	"testing"
	"time"
)

func TestCombination(t *testing.T) {
	operatingInfo := object.OperatingInfo{
		FlowID:        "1",
		AprovUserId:   "1",
		Username:      "1",
		XconfBefore:   "1",
		XconfAfter:    "1",
		ComposeBefore: "1",
		ComposeAfter:  "1",
		CreateTime:    time.Now().Local(),
	}
	handler.OperatingInfoAdd(operatingInfo)
	// 关闭数据连接
	store.Close()
}
func TestUserTask(t *testing.T) {
	taskInfo := object.TaskInfo{
		TaskId:     "2",
		CreateTime: time.Now().Local(),
	}
	handler.UserTaskAddOrUpdate("add", taskInfo, "1", "unfin")
	// 关闭数据连接
	store.Close()
}
