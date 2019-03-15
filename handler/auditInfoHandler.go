package handler

import (
	"encoding/json"
	"ggz-server/object"
	"ggz-server/util"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

//根据nodeid查询所有，可根据名字的条件查询
func Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeId := vars["nodeId"]
	username := vars["username"]
	if nodeId == "" {
		glog.Error("nodeId 为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}
	data, err := FindAll(nodeId)
	if err != nil {
		if err.Error() != "Key not found" {
			glog.Error(err)
			util.WriteJsonString(w, object.NewServerErrReturnObj())
			return
		}
	}
	if username == "" || data == nil {
		util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(data))
		return
	}
	auditInfos := make([]object.AuditInfo, 0)
	for _, v1 := range data.AuditInfo {
		if username == v1.Username {
			auditInfos = append(auditInfos, v1)
		}
	}
	util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(auditInfos))
}

//修改
func Update(w http.ResponseWriter, r *http.Request) {
	AddOrUpdate(w, r, "update")
}

//增加
func Add(w http.ResponseWriter, r *http.Request) {
	AddOrUpdate(w, r, "add")
}

//删除
func Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aprovUserId := vars["aprovUserId"]
	nodeId := vars["nodeId"]
	if aprovUserId == "" || nodeId == "" {
		glog.Error("员工工号或nodeId 为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}
	data, err := FindAll(nodeId)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	if data == nil {
		util.WriteJsonString(w, object.NewReturnObj(500, "无审批信息，请检查nodeId是否正确", nodeId))
		return
	}
	auditInfos := make([]object.AuditInfo, 0)
	for _, v1 := range data.AuditInfo {
		if aprovUserId == v1.AprovUserId {
			continue
		}
		auditInfos = append(auditInfos, v1)
	}
	auditInfosJson, _ := json.Marshal(auditInfos)
	auditInfo := object.AuditInfo{AprovUserId: aprovUserId, NodeId: nodeId}
	addErr := auditInfo.Upate(auditInfosJson, nodeId)
	if addErr != nil {
		glog.Error(addErr)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	util.WriteJsonString(w, object.NewSuccessReturnObj())
}

//查看详情
func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aprovUserId := vars["aprovUserId"]
	nodeId := vars["nodeId"]
	if aprovUserId == "" || nodeId == "" {
		glog.Error("员工工号或nodeId 为空")
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}
	data, err := FindAll(nodeId)
	if err != nil {
		glog.Error(err)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	if data == nil {
		util.WriteJsonString(w, object.NewReturnObj(500, "无审批信息，请检查nodeId是否正确", nodeId))
		return
	}
	for _, v1 := range data.AuditInfo {
		if aprovUserId == v1.AprovUserId {
			util.WriteJsonString(w, object.NewSuccessWithDataReturnObj(v1))
			return
		}
	}
	auditInfo := object.AuditInfo{
		NodeId:      nodeId,
		AprovUserId: aprovUserId,
	}
	util.WriteJsonString(w, object.NewReturnObj(500, "请核对员工编号是否正确", auditInfo))
}

func FindAll(nodeId string) (*object.AuditInfos, error) {
	auditInfo := object.AuditInfo{NodeId: nodeId}
	data, err := auditInfo.FindById()
	if err != nil {
		return nil, err
	}
	if string(data) == "" {
		return nil, nil
	}
	log.Println("查询结果" + string(data))
	auditInfos := object.AuditInfos{}
	jsonErr := json.Unmarshal(data, &auditInfos.AuditInfo)
	if jsonErr != nil {
		glog.Error(jsonErr)
		return nil, err
	}
	return &auditInfos, nil
}
func AddOrUpdate(w http.ResponseWriter, r *http.Request, method string) {
	status := 0
	nodeId := r.FormValue("nodeId")
	aprovUserId := r.FormValue("aprovUserId")
	username := r.FormValue("username")
	mailAddr := r.FormValue("mailAddr")
	comment := r.FormValue("comment")
	statusString := r.FormValue("status")
	if statusString != "" {
		statusInt, e := strconv.Atoi(statusString)
		if e == strconv.ErrRange || e == strconv.ErrSyntax {
			util.WriteJsonString(w, object.NewReturnObj(500, "状态格式错误", statusString))
			return
		}
		status = statusInt
	}
	auditInfo := object.AuditInfo{
		NodeId:      nodeId,
		AprovUserId: aprovUserId,
		Username:    username,
		MailAddr:    mailAddr,
		Comment:     comment,
		Status:      status,
		CreateTime:  time.Now().Local(),
	}
	//查询出所有，遍历
	data, err := FindAll(nodeId)
	if err != nil {
		if err.Error() != "Key not found" {
			glog.Error(err)
			util.WriteJsonString(w, object.NewServerErrReturnObj())
			return
		}
	}
	auditInfos := make([]object.AuditInfo, 0)
	isUpdate := false
	if data == nil {
		//初始表，增加
		if method == "add" {
			auditInfos = append(auditInfos, auditInfo)
			auditInfoJson, _ := json.Marshal(auditInfos)
			addErr := auditInfo.Upate(auditInfoJson, nodeId)
			if addErr != nil {
				glog.Error(addErr)
				util.WriteJsonString(w, object.NewServerErrReturnObj())
				return
			}
			util.WriteJsonString(w, object.NewSuccessReturnObj())
			return
		}
		util.WriteJsonString(w, object.NewParamErrReturnObj())
		return
	}
	if method == "add" {
		for _, v1 := range data.AuditInfo {
			if aprovUserId == v1.AprovUserId {
				util.WriteJsonString(w, object.NewReturnObj(500, "工号已存在，添加失败", auditInfo))
				return
			}
			auditInfos = append(auditInfos, v1)
		}
		//代表增加
		auditInfos = append(auditInfos, auditInfo)
	} else if method == "update" {
		for _, v1 := range data.AuditInfo {
			if aprovUserId == v1.AprovUserId {
				//代表已经包含修改
				if comment != "" {
					v1.Comment = comment
				}
				if mailAddr != "" {
					v1.MailAddr = mailAddr
				}
				if statusString != "" {
					v1.Status = status
				}
				if username != "" {
					v1.Username = username
				}
				v1.UpdateTime = time.Now().Local()
				isUpdate = true
			}
			auditInfos = append(auditInfos, v1)
		}
		if !isUpdate {
			//代表没有修改的人
			util.WriteJsonString(w, object.NewReturnObj(500, "工号不存在，修改失败", auditInfo))
			return
		}
	} else {
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	auditInfosJson, _ := json.Marshal(auditInfos)
	log.Println("插入json" + string(auditInfosJson))
	addErr := auditInfo.Upate(auditInfosJson, nodeId)
	if addErr != nil {
		glog.Error(addErr)
		util.WriteJsonString(w, object.NewServerErrReturnObj())
		return
	}
	util.WriteJsonString(w, object.NewSuccessReturnObj())
}
