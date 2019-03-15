package handler

import (
	"encoding/json"
	"ggz-server/object"
	"ggz-server/store"
	"ggz-server/util"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

//查看完成详情
func UserTaskFindFin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aprovUserId := vars["aprovUserId"]
	reqType := vars["type"]
	result, err, _ := UserTaskFind(aprovUserId, reqType)
	if err != nil {
		if err.Error() != "Key not found" {
			glog.Error(err)
			util.WriteJsonString(w, object.NewServerErrReturnObj())
			return
		}
	}
	util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(result))
}
func UserTaskAddOrUpdate(method string, rep object.TaskInfo, aprovUserId string, reqType string) error {
	result, err, key := UserTaskFind(aprovUserId, reqType)
	if err != nil {
		if err.Error() != "Key not found" {
			return err
		}
	}
	log.Println("用户任务修改,key:" + key)
	if method == "add" {
		resultJson, _ := json.Marshal(append(result, rep))
		log.Println("插入json" + string(resultJson))
		err1 := store.Store(key, resultJson)
		if err1 != nil {
			glog.Error(err)
			return err
		}
		return nil
	} else {
		//删除，先遍历，然后找对应的流程id删除
		isUpdate := false
		taskInfos := make([]object.TaskInfo, 0)
		for _, val := range result {
			if val.TaskId == rep.TaskId {
				isUpdate = true
				continue
			}
			taskInfos = append(taskInfos, val)
		}
		if !isUpdate {
			return errors.New("删除失败，根据员工号查询无任务信息")
		}
		taskInfosJson, _ := json.Marshal(taskInfos)
		log.Println("重新加入json" + string(taskInfosJson))
		addErr := store.Store(key, taskInfosJson)
		if addErr != nil {
			glog.Error(addErr)
			return addErr
		}
		return nil
	}
}

//查询
func UserTaskFind(aprovUserId string, method string) (result []object.TaskInfo, error error, key string) {
	//var key string
	if method == "fin" {
		//已经处理
		key = util.UserTaskFinKey + aprovUserId
	}
	if method == "unfin" {
		//已经处理
		key = util.UserTaskUnFinKey + aprovUserId
	}
	//根据不同的key查询
	taskInfos := []object.TaskInfo{}
	err := object.FindAll(key, &taskInfos)
	return taskInfos, err, key
}
