package object

import (
	"encoding/json"
	"ggz-server/store"
	"github.com/golang/glog"
	"log"
)

func FindAll(key string, result interface{}) error {
	data, err := store.View(key)
	if err != nil {
		return err
	}
	if string(data) == "" {
		return nil
	}
	log.Println("查询结果" + string(data))
	jsonErr := json.Unmarshal(data, &result)
	if jsonErr != nil {
		glog.Error(jsonErr)
		return err
	}
	result = &result
	return nil
}
