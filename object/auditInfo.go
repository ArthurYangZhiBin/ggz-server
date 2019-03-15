package object

import (
	"ggz-server/store"
	"ggz-server/util"
	"time"
)

type AuditInfo struct {
	NodeId      string    `json:"nodeId"`
	AprovUserId string    `json:"aprovUserId"`
	Username    string    `json:"username"`
	MailAddr    string    `json:"mailAddr"`
	Comment     string    `json:"comment"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"create_time"` //创建时间
	UpdateTime  time.Time `json:"update_time"` //更新时间
}
type AuditInfos struct {
	//普通struct类型
	AuditInfo []AuditInfo `json:"auditInfo"`
}

//根据传入的id号查询节点信息
func (auditInfo *AuditInfo) FindById() ([]byte, error) {
	return store.View(util.AuditInfoKey + "-" + auditInfo.NodeId)
}

//修改信息
func (auditInfo *AuditInfo) Upate(auditInfos []byte, key string) error {
	return store.Store(util.AuditInfoKey+"-"+key, auditInfos)
}
