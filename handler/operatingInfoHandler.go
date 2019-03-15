package handler

import (
	"encoding/json"
	"ggz-server/object"
	"ggz-server/store"
	"ggz-server/util"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//根据流程flowID查询所有这个操作的中间表
func OperatingInfoFind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["flowID"]
	if flowID == "" {
		//查询所有流程id对应的操作，模糊查询key，返回
		data, err := store.Blurry(util.OperatingInfoKey)
		if err != nil {
			if err.Error() != "Key not found" {
				glog.Error(err)
				util.WriteJsonString(w, object.NewServerErrReturnObj())
				return
			}
		}
		/*for _, val := range data {
			log.Println("操作表查询所有" + val)
		}*/
		util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(data))
		return
	}
	//根据flowID查询详细信息
	operatingInfos := object.OperatingInfos{}
	err := object.FindAll(util.OperatingInfoKey+"-"+flowID, &operatingInfos.OperatingInfo)
	if err != nil {
		if err.Error() != "Key not found" {
			glog.Error(err)
			util.WriteJsonString(w, object.NewServerErrReturnObj())
			return
		}
	}
	util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(operatingInfos))
}

//增加
func OperatingInfoAdd(info object.OperatingInfo) error {
	//根据flowID查询详细信息
	operatingInfos := object.OperatingInfos{}
	err := object.FindAll(util.OperatingInfoKey+"-"+info.FlowID, &operatingInfos.OperatingInfo)
	if err != nil {
		if err.Error() != "Key not found" {
			glog.Error(err)
			return err
		}
	}
	operatingInfoJson, _ := json.Marshal(append(operatingInfos.OperatingInfo, info))
	log.Println("插入json" + string(operatingInfoJson))
	err1 := store.Store(util.OperatingInfoKey+"-"+info.FlowID, operatingInfoJson)
	if err1 != nil {
		glog.Error(err)
		return err
	}
	return nil
}
